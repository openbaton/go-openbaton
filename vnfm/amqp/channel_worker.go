package amqp

import (
	"github.com/openbaton/go-openbaton/util"
	"github.com/openbaton/go-openbaton/vnfm/channel"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func (acnl *Channel) spawnWorkers() {
	acnl.wg.Add(acnl.numOfWorkers)
	for i := 0; i < acnl.numOfWorkers; i++ {
		go acnl.worker(i)
	}
}

func (acnl *Channel) worker(id int) {
	tag := util.FuncName()

	acnl.l.WithFields(log.Fields{
		"tag":       tag,
		"worker-id": id,
	}).Debug("AMQP worker starting")

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
		cnl, err := acnl.getAMQPChan()
		if err != nil {
			// fetch a new one
			continue WorkerLoop
		}

		if cnl == nil {
			break WorkerLoop
		}

		errChan := cnl.NotifyClose(make(chan *amqp.Error))

		acnl.l.WithFields(log.Fields{
			"tag": tag,
		}).Debug("new AMQP channel received")

	ServeLoop:
		for {
			select {
			// Updates the status. If it becomes Running, the next loop will accept incoming jobs again
			case status = <-acnl.statusChan:
				continue ServeLoop // if the status becomes "quitting", the channel will be closed and an error will be received on the errChan.

			// if anything is received on errChan, it means that the chan we were used
			// was closed by something; try to get another one.
			// If the connection was cleanly closed, the getAMQPChan request above will return nil, and the worker will
			// shutdown.
			case <-errChan:

			case exc := <-work():
				if exc == nil {
					continue ServeLoop
				}

				if exc.replyChan != nil { // RPC request
					resp, err := acnl.rpc(cnl, exc.queue, exc.msg)

					exc.replyChan <- response{resp, err}
				} else { //send only
					if err := acnl.publish(cnl, exc.queue, exc.msg); err != nil {
						acnl.l.WithError(err).WithFields(log.Fields{
							"tag":       "worker-amqp",
							"worker-id": id,
						}).Error("publish failed")
					}
				}
			}
		}
	}

	acnl.l.WithFields(log.Fields{
		"tag":       "worker-amqp",
		"worker-id": id,
	}).Debug("AMQP worker stopping")

	acnl.wg.Done()
}
