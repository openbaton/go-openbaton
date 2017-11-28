package pluginsdk

import (
	"fmt"
	"reflect"
	"encoding/json"
	"encoding/base64"
	"github.com/op/go-logging"
)

var (
	ErrProtocolFail error = plugError{"protocol error, NFVO and plugin are out of sync"}
)

type worker struct {
	l *logging.Logger
	h HandlerVim
}

func (w worker) handle(fname string, args []json.RawMessage) (interface{}, error) {
	if len(args) < 1 {
		return nil, plugError{"expecting at least one VimInstance"}
	}

	fValue, err := w.matchFunc(fname, args)
	if err != nil {
		return nil, err
	}


	fType := fValue.Type()

	// matchFunc ensures the length of args is correct
	callArgs := make([]reflect.Value, len(args))
	for i, jsonArg := range args {
		argType := fType.In(i)

		kind := argType.Kind()
		if kind == reflect.Interface && i == 0 { // special case: VimInstance at first position
			argValue := map[string]interface{}{}
			if err := json.Unmarshal(jsonArg, &argValue); err != nil {
				return nil, err
			}

			//callArgs[i] = &catalogue.DockerVimInstance{}
			callArgs[i] = reflect.ValueOf(GetVimInstance(jsonArg, argValue))
		} else if kind == reflect.Slice && argType.Elem().Kind() == reflect.Uint8 { // special case: argument is an array of bytes
			var baseStr string
			if err := json.Unmarshal(jsonArg, &baseStr); err != nil {
				return nil, err
			}

			b, err := base64.StdEncoding.DecodeString(baseStr)
			if err != nil {
				return nil, plugError{"base64 decoding failed"}
			}

			callArgs[i] = reflect.ValueOf(b)
		} else {
			// create a new pointer to the arg type, and deserialise into it
			// its JSON
			argValue := reflect.New(argType)
			if err := json.Unmarshal(jsonArg, argValue.Interface()); err != nil {
				return nil, err
			}

			callArgs[i] = reflect.ValueOf(argValue.Elem().Interface())
		}
	}

	retVals := fValue.Call(callArgs)
	var errVal, respVal reflect.Value
	switch len(retVals) {
	// function returning just an error
	case 1:
		errVal = retVals[0]

	case 2:
		respVal, errVal = retVals[0], retVals[1]

	default:
		return nil, plugError{"function returns an unsupported number of values"}
	}

	if errIface := errVal.Interface(); errIface != nil {
		err, ok := errIface.(error)
		if !ok {
			return nil, plugError{"unsupported function - no error value"}
		}

		return nil, err
	}

	if len(retVals) == 2 {
		if !respVal.IsValid() {
			return nil, plugError{"broken response from driver"}
		}

		return respVal.Interface(), nil
	}

	return nil, nil
}

func (w worker) matchFunc(fname string, args []json.RawMessage) (reflect.Value, error) {
	var fVal reflect.Value

	switch fname {
	case "addFlavor":
		fVal = reflect.ValueOf(w.h.AddFlavour)

		// overloaded function
	case "addImage":
		// we need to check if the last argument is a byte array or a string
		fVal = reflect.ValueOf(w.h.AddImage)
		fType := fVal.Type()
		if len(args) != fType.NumIn() {
			break // will make the same check below and fail there
		}

		lastArg := args[len(args)-1]
		var str string
		if err := json.Unmarshal(lastArg, &str); err != nil {
			return fVal, err
		}

		// if base64 deserialisation fails, then this is an URL string
		if _, e := base64.StdEncoding.DecodeString(str); e != nil {
			fVal = reflect.ValueOf(w.h.AddImageFromURL)
		}

	case "copyImage":
		fVal = reflect.ValueOf(w.h.CopyImage)

	case "createNetwork":
		fVal = reflect.ValueOf(w.h.CreateNetwork)

	case "createSubnet":
		fVal = reflect.ValueOf(w.h.CreateSubnet)

	case "deleteFlavor":
		fVal = reflect.ValueOf(w.h.DeleteFlavour)

	case "deleteImage":
		fVal = reflect.ValueOf(w.h.DeleteImage)

	case "deleteNetwork":
		fVal = reflect.ValueOf(w.h.DeleteNetwork)

	case "deleteServerByIdAndWait":
		fVal = reflect.ValueOf(w.h.DeleteServerByIDAndWait)

	case "deleteSubnet":
		fVal = reflect.ValueOf(w.h.DeleteSubnet)

	case "getNetworkById":
		fVal = reflect.ValueOf(w.h.NetworkByID)

	case "getQuota":
		fVal = reflect.ValueOf(w.h.Quota)

	case "getSubnetsExtIds":
		fVal = reflect.ValueOf(w.h.SubnetsExtIDs)

	case "getType":
		fVal = reflect.ValueOf(w.h.Type)

	case "launchInstance":
		fVal = reflect.ValueOf(w.h.LaunchInstance)

		// overloaded function
	case "launchInstanceAndWait":
		fVal = reflect.ValueOf(w.h.LaunchInstanceAndWait)

		// check for overloaded functions
		if len(args) == reflect.TypeOf(w.h.LaunchInstanceAndWaitWithIPs).NumIn() {
			fVal = reflect.ValueOf(w.h.LaunchInstanceAndWaitWithIPs)
		}

	case "listFlavors":
		fVal = reflect.ValueOf(w.h.ListFlavours)

	case "refresh":
		fVal = reflect.ValueOf(w.h.Refresh)

	case "listImages":
		fVal = reflect.ValueOf(w.h.ListImages)

	case "listNetworks":
		fVal = reflect.ValueOf(w.h.ListNetworks)

	case "listServer":
		fVal = reflect.ValueOf(w.h.ListServer)

	case "updateFlavor":
		fVal = reflect.ValueOf(w.h.UpdateFlavour)

	case "updateImage":
		fVal = reflect.ValueOf(w.h.UpdateImage)

	case "updateNetwork":
		fVal = reflect.ValueOf(w.h.UpdateNetwork)

	case "updateSubnet":
		fVal = reflect.ValueOf(w.h.UpdateSubnet)

	default:
		return fVal, ErrProtocolFail
	}

	fType := fVal.Type()
	if len(args) != fType.NumIn() {
		return fVal, plugError{fmt.Sprintf("wrong number of arguments (%d) for function %s", len(args), fType.String())}
	}

	return fVal, nil
}
