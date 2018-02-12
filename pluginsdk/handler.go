package pluginsdk

import (
	"errors"
	"encoding/json"
	"github.com/openbaton/go-openbaton/sdk"
	"github.com/streadway/amqp"
)

func handlePluginRequest(bytemsg []byte, handler sdk.Handler, allocate bool, connection *amqp.Connection) ([]byte, error) {
	var req request
	logger := sdk.GetLogger("handler-plugin-function", "DEBUG")
	if err := json.Unmarshal(bytemsg, &req); err != nil {
		logger.Error("message unmarshaling error")
		return nil, errors.New("message unmarshaling error")
	}

	switch h := handler.(type) {
	case HandlerVim:
		wk := &worker{
			l: logger,
			h: h,
		}
		result, err := wk.handle(req.MethodName, req.Parameters)
		var resp response
		if err != nil {
			switch err.(type) {

			case plugError:
				resp.Exception = err
			case sdk.DriverError:
				resp.Exception = err
			default:
				resp.Exception = plugError{err.Error()}
			}
		} else {
			resp.Answer = result
		}

		bResp, err := json.MarshalIndent(resp,"","  ")
		if err != nil {
			logger.Error("failure while serialising response")
			return nil, err
		}
		logger.Debugf("Returning %v", string(bResp))
		return bResp, nil
	default:
		logger.Errorf("Error, worker of wrong type")
		return nil, errors.New("worker of wrong type")
	}

}
