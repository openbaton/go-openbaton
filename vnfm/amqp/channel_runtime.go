package amqp

import (
	"errors"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/streadway/amqp"
)

var (
	ErrTimedOut = errors.New("timed out")
)

func (acnl *amqpChannel) closeQueues() {
	acnl.setStatus(channel.Quitting)

	close(acnl.statusChan)
	close(acnl.sendQueue)
	close(acnl.subChan)
	close(acnl.quitChan)
	close(acnl.receiverDeliveryChan)
}

func (acnl *amqpChannel) setStatus(newStatus channel.Status) {
	for i := 0; i < acnl.numOfWorkers; i++ {
		acnl.statusChan <- newStatus
	}

	acnl.status = newStatus
}

// spawn spawns the main handler for AMQP communications.
func (acnl *amqpChannel) spawn() error {
	errChan, err := acnl.setup()
	if err != nil {
		return err
	}

	acnl.register()

	acnl.spawnWorkers()
	acnl.setStatus(channel.Running)

	go func() {
		for {
			select {
			case <-acnl.quitChan:
				if err := acnl.conn.Close(); err != nil {
					acnl.l.Errorf("while closing AMQP Connection: %v\n", err)

					acnl.closeQueues()

					if err := acnl.unregister(); err != nil {
						acnl.l.Errorf("unregister failed: %v\n", err)
					}

					return
				}
				// Close will cause the reception of nil on errChan.

			case err = <-errChan:
				// The connection closed cleanly after invoking Close().
				if err == nil {
					// notify the receiving end and listeners
					acnl.closeQueues()

					return
				}

				acnl.setStatus(channel.Reconnecting)

				// The connection crashed for some reason. Try to bring it up again.
				for {
					if errChan, err = acnl.setup(); err != nil {
						acnl.l.Errorln("can't re-establish connection with AMQP; queues stalled. Retrying in 30 seconds.")
						time.Sleep(30 * time.Second)
					} else {
						acnl.setStatus(channel.Running)
						break
					}
				}

			}
		}
	}()

	return nil
}

// spawnReceiver spawns a goroutine which handles the reception of
// incoming messages from the NFVO on a dedicated queue.
// The receiver main channel is updated by setup() with a new
// consumer each time the connection is reestablished.
func (acnl *amqpChannel) spawnReceiver() {
	go func() {
		acnl.l.Infoln("receiver: spawned")

		// list of channels to which incoming messages will be broadcasted.
		notifyChans := []chan<- messages.NFVMessage{}

		var deliveryChan <-chan amqp.Delivery
	RecvLoop:
		for {
			select {
			// setup delivers a new channel to this receiver, to
			// be listened for Deliveries.
			case deliveryChan = <-acnl.receiverDeliveryChan:
				if deliveryChan == nil {
					break RecvLoop
				}
				acnl.l.Debugln("receiver: new delivery channel received")
				// chan updated

			// receives and adds a chan to the list of notifyChans
			case notifyChan := <-acnl.subChan:
				acnl.l.Debugln("receiver: new notify channel received")
				notifyChans = append(notifyChans, notifyChan)

			case delivery, ok := <-deliveryChan:
				if ok {
					msg, err := messages.Unmarshal(delivery.Body)
					if err != nil {
						acnl.l.Errorf("while receiving message: %v\n", err)
						continue RecvLoop
					}

					last := 0
					for _, c := range notifyChans {
						select {
						// message sent successfully.
						case c <- msg:
							// keep the channel around for the next time
							notifyChans[last] = c
							last++

						// nobody is listening at the other end of the channel.
						case <-time.After(1 * time.Second):
							close(c)
						}
					}

					// notifyChans trimmed of dead chans
					notifyChans = notifyChans[last:]
				} else {
					// make deliveryChan nil if someone closes it:
					// a closed channel always immediately returns a zero value, thus never
					// allowing the select to block.
					// A nil channel always blocks.
					deliveryChan = nil
				}
			}
		}

		// closing all the notification channels
		for _, cnl := range notifyChans {
			close(cnl)
		}

		acnl.l.Infoln("receiver: exiting")
	}()
}
func (acnl *amqpChannel) spawnWorkers() {
	for i := 0; i < acnl.numOfWorkers; i++ {
		go acnl.worker(i)
	}
}

func (acnl *amqpChannel) worker(id int) {
	acnl.l.Infof("starting worker %d\n", id)

	status := channel.Stopped

	// explanation: a read on a nil channel will
	// block forever. This lambda ensures that we will accept jobs only
	// when the status is valid.
	work := func() chan *exchange {
		if status == channel.Running {
			return acnl.sendQueue
		}

		return nil
	}

WorkerLoop:
	for {
		select {
		// Updates the status. If it becomes Running, the next loop will accept incoming jobs again
		case status = <-acnl.statusChan:
			if status == channel.Quitting {
				break WorkerLoop
			}

		case exc := <-work():
			if exc.replyChan != nil { // RPC request
				resp, err := acnl.rpc(exc.queue, exc.msg)

				exc.replyChan <- response{resp, err}
			} else { //send only
				if err := acnl.publish(exc.queue, exc.msg); err != nil {
					acnl.l.Errorf("while publishing from worker #%d: %v\n", id, err)
				}
			}
		}
	}

	acnl.l.Infof("quitting worker %d\n", id)
}
