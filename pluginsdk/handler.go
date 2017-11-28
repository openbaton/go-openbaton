package pluginsdk

import (
	"log"
	"errors"
	"encoding/json"
	"github.com/openbaton/go-openbaton/sdk"
)

type driver struct {
	*log.Logger
}

func handlePluginRequest(bytemsg []byte, wk interface{}) ([]byte, error){
	var req request
	logger := sdk.GetLogger("handler-plugin-function", "DEBUG")
	if err := json.Unmarshal(bytemsg, &req); err != nil {
		logger.Error("message unmarshaling error")
		return nil, errors.New("message unmarshaling error")
	}
	switch t := wk.(type) {
		case *worker:
			result, err := t.handle(req.MethodName, req.Parameters)

			var resp response
			if err != nil {
				// The NFVO expects a Java Exception;
				// This type switch checks if the error is not one of the special
				// Java-compatible types already and wraps it.
				switch err.(type) {

				case plugError:
					resp.Exception = err
				case sdk.DriverError:
					resp.Exception = err

					// if the error is not a special plugin error, than wrap it:
					// the nfvo expects a Java exception.
				default:
					resp.Exception = plugError{err.Error()}
				}
			} else {
				resp.Answer = result
			}

			bResp, err := json.Marshal(resp)
			if err != nil {
				logger.Error("failure while serialising response")
				return nil, err
			}
			logger.Debugf("Returning %v", string(bResp))
			return bResp, nil
	default:
		logger.Errorf("Error, worker of wrong type")
		return nil, errors.New("Error, worker of wrong type")
	}

}