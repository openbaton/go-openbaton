package amqp

import (
	"errors"
	"time"

	"github.com/mcilloni/go-openbaton/vnfm/channel"
)

var (
	ErrTimedOut = errors.New("timed out")
)

func (acnl *amqpChannel) closeQueues() {
	close(acnl.sendQueue)
	close(acnl.quitChan)

	for _, cnl := range acnl.notifyChans {
		close(cnl)
	}
}

func (acnl *amqpChannel) receiver() {

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
			case notifyChan := <-acnl.subChan:
				acnl.notifyChans = append(acnl.notifyChans, notifyChan)

			case <-acnl.quitChan:
				if err := acnl.conn.Close(); err != nil {
					acnl.l.Errorf("while closing AMQP Connection: %v\n", err)

					acnl.closeQueues()

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

	for {
		select {
		case status = <-acnl.statusChan:
			// Updates the status. If it becomes Running, the next loop will accept incoming jobs again

		case exc := <-work():
			// the sender expects a reply
			if exc.replyChan != nil {

			}
		}
	}

	acnl.l.Infof("quitting worker %d\n", id)
}
