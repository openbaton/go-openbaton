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

func (dh driverHandler) Handle(fname string, args []json.RawMessage) (resp interface{}, err error) {
	if len(args) < 1 {
		return nil, plugError{"expecting at least one VimInstance"}
	}

	var fn interface{}
	if fn, err = dh.matchFunc(fname, args); err != nil {
		fValue := reflect.ValueOf(fn)
		fType := fValue.Type()

		// matchFunc ensures the length of args is correct
		callArgs := make([]reflect.Value, len(args))
		for i, jsonArg := range args {
			argType := fType.In(i)
			// special case: argument is an array of bytes
			if argType.Kind() == reflect.Slice && argType.Elem().Kind() == reflect.Uint8 {
				var baseStr string
				if err = json.Unmarshal(jsonArg, &baseStr); err != nil {
					return
				}

				b, e := base64.StdEncoding.DecodeString(baseStr)
				if e != nil {
					return nil, plugError{"base64 decoding failed"}
				}

				callArgs[i] = reflect.ValueOf(b)
			} else {
				// create a new pointer to the arg type, and deserialise into it
				// its JSON
				argValue := reflect.New(argType)
				if err = json.Unmarshal(jsonArg, argValue.Interface()); err != nil {
					return
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

		e, ok := errVal.Interface().(error)
		if !ok {
			return nil, plugError{"unsupported function - no error value"}
		}

		err = e

		if respVal.IsValid() {
			resp = respVal.Interface()
		}
	}

	return
}

func (dh driverHandler) Type() string {
	return "vim-driver"
}

func (dh driverHandler) matchFunc(fname string, args []json.RawMessage) (f interface{}, err error) {
	switch fname {
	case "addFlavor":
		f = dh.Driver.AddFlavour

	// overloaded function
	case "addImage":
		// we need to check if the last argument is a byte array or a string
		f = dh.AddImage
		fType := reflect.TypeOf(f)
		if len(args) != fType.NumIn() {
			break // will make the same check below and fail there
		}

		lastArg := args[len(args)-1]
		var str string
		if err = json.Unmarshal(lastArg, &str); err != nil {
			return
		}

		// if base64 deserialisation fails, then this is an URL string
		if _, e := base64.StdEncoding.DecodeString(str); e != nil {
			f = dh.AddImageFromURL
		}

	case "copyImage":
		f = dh.CopyImage

	case "createNetwork":
		f = dh.CreateNetwork

	case "createSubnet":
		f = dh.CreateSubnet

	case "deleteFlavor":
		f = dh.DeleteFlavour

	case "deleteImage":
		f = dh.DeleteImage

	case "deleteNetwork":
		f = dh.DeleteNetwork

	case "deleteServerByIdAndWait":
		f = dh.DeleteServerByIDAndWait

	case "deleteSubnet":
		f = dh.DeleteSubnet

	case "getNetworkById":
		f = dh.GetNetworkByID

	case "getQuota":
		f = dh.GetQuota

	case "getSubnetsExtIds":
		f = dh.GetSubnetsExtIDs

	case "getType":
		f = dh.GetType

	case "launchInstance":
		f = dh.LaunchInstance

	// overloaded function
	case "launchInstanceAndWait":
		f = dh.LaunchInstanceAndWait

		// check for overloaded functions
		if len(args) == reflect.TypeOf(dh.LaunchInstanceAndWaitWithIPs).NumIn() {
			f = dh.LaunchInstanceAndWaitWithIPs
		}

	case "listFlavors":
		f = dh.ListFlavours

	case "listImages":
		f = dh.ListImages

	case "listNetworks":
		f = dh.ListNetworks

	case "listServer":
		f = dh.ListServer

	case "updateFlavor":
		f = dh.UpdateFlavour

	case "updateImage":
		f = dh.UpdateImage

	case "updateNetwork":
		f = dh.UpdateNetwork

	case "updateSubnet":
		f = dh.UpdateSubnet

	default:
		err = ErrProtocolFail
	}

	fType := reflect.TypeOf(f)
	if len(args) != fType.NumIn() {
		err = plugError{fmt.Sprintf("wrong number of arguments (%d) for function %s", len(args), fType.String())}
	}

	return
}
