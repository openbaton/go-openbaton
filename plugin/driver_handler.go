package plugin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	//"github.com/mcilloni/go-openbaton/catalogue"
	log "github.com/sirupsen/logrus"
)

var (
	ErrProtocolFail error = plugError{"protocol error, NFVO and plugin are out of sync"}
)

type driverHandler struct {
	Driver

	l *log.Logger
}

func (dh driverHandler) Handle(fname string, args []json.RawMessage) (interface{}, error) {
	if len(args) < 1 {
		return nil, plugError{"expecting at least one VimInstance"}
	}

	fValue, err := dh.matchFunc(fname, args)
	if err != nil {
		return nil, err
	}

	fType := fValue.Type()

	// matchFunc ensures the length of args is correct
	callArgs := make([]reflect.Value, len(args))
	for i, jsonArg := range args {
		argType := fType.In(i)
		// special case: argument is an array of bytes
		if argType.Kind() == reflect.Slice && argType.Elem().Kind() == reflect.Uint8 {
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

	err, ok := errVal.Interface().(error)
	if !ok {
		return nil, plugError{"unsupported function - no error value"}
	}

	if err != nil {
		return nil, err
	}

	if !respVal.IsValid() {
		return nil, plugError{"broken response from handler"}
	}

	return respVal.Interface(), nil
}

func (dh driverHandler) QueueTag() string {
	return "vim-drivers"
}

func (dh driverHandler) Type() string {
	return "vim-driver"
}

func (dh driverHandler) matchFunc(fname string, args []json.RawMessage) (reflect.Value, error) {
	var fVal reflect.Value

	switch fname {
	case "addFlavor":
		fVal = reflect.ValueOf(dh.Driver.AddFlavour)

	// overloaded function
	case "addImage":
		// we need to check if the last argument is a byte array or a string
		fVal = reflect.ValueOf(dh.AddImage)
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
			fVal = reflect.ValueOf(dh.AddImageFromURL)
		}

	case "copyImage":
		fVal = reflect.ValueOf(dh.CopyImage)

	case "createNetwork":
		fVal = reflect.ValueOf(dh.CreateNetwork)

	case "createSubnet":
		fVal = reflect.ValueOf(dh.CreateSubnet)

	case "deleteFlavor":
		fVal = reflect.ValueOf(dh.DeleteFlavour)

	case "deleteImage":
		fVal = reflect.ValueOf(dh.DeleteImage)

	case "deleteNetwork":
		fVal = reflect.ValueOf(dh.DeleteNetwork)

	case "deleteServerByIdAndWait":
		fVal = reflect.ValueOf(dh.DeleteServerByIDAndWait)

	case "deleteSubnet":
		fVal = reflect.ValueOf(dh.DeleteSubnet)

	case "getNetworkById":
		fVal = reflect.ValueOf(dh.NetworkByID)

	case "getQuota":
		fVal = reflect.ValueOf(dh.Quota)

	case "getSubnetsExtIds":
		fVal = reflect.ValueOf(dh.SubnetsExtIDs)

	case "getType":
		fVal = reflect.ValueOf(dh.Type)

	case "launchInstance":
		fVal = reflect.ValueOf(dh.LaunchInstance)

	// overloaded function
	case "launchInstanceAndWait":
		fVal = reflect.ValueOf(dh.LaunchInstanceAndWait)

		// check for overloaded functions
		if len(args) == reflect.TypeOf(dh.LaunchInstanceAndWaitWithIPs).NumIn() {
			fVal = reflect.ValueOf(dh.LaunchInstanceAndWaitWithIPs)
		}

	case "listFlavors":
		fVal = reflect.ValueOf(dh.ListFlavours)

	case "listImages":
		fVal = reflect.ValueOf(dh.ListImages)

	case "listNetworks":
		fVal = reflect.ValueOf(dh.ListNetworks)

	case "listServer":
		fVal = reflect.ValueOf(dh.ListServer)

	case "updateFlavor":
		fVal = reflect.ValueOf(dh.UpdateFlavour)

	case "updateImage":
		fVal = reflect.ValueOf(dh.UpdateImage)

	case "updateNetwork":
		fVal = reflect.ValueOf(dh.UpdateNetwork)

	case "updateSubnet":
		fVal = reflect.ValueOf(dh.UpdateSubnet)

	default:
		return fVal, ErrProtocolFail
	}

	fType := fVal.Type()
	if len(args) != fType.NumIn() {
		return fVal, plugError{fmt.Sprintf("wrong number of arguments (%d) for function %s", len(args), fType.String())}
	}

	return fVal, nil
}
