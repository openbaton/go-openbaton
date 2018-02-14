package pluginsdk

import (
	"github.com/openbaton/go-openbaton/catalogue"
	"errors"
	"encoding/json"
	"reflect"
)

// Obtain a Docker vim instance from a interfaced struct
func GetDockerVimInstance(vimInstance interface{}) (*catalogue.DockerVimInstance, error) {
	switch t := vimInstance.(type) {
	case *catalogue.DockerVimInstance:
		return t, nil
	default:
		return nil, errors.New("not Received Docker Vim Instance");
	}
}

// Obtain a Openstack vim instance from a interfaced struct
func GetOpenstackVimInstance(vimInstance interface{}) (*catalogue.OpenstackVimInstance, error) {
	switch t := vimInstance.(type) {
	case *catalogue.OpenstackVimInstance:
		return t, nil
	default:
		return nil, errors.New("not Received Docker Vim Instance");
	}
}

// Obtain a Generic vim instance from a interfaced struct
func GetBaseVimInstance(vimInstance interface{}) (*catalogue.BaseVimInstance, error) {
	switch t := vimInstance.(type) {
	case *catalogue.BaseVimInstance:
		return t, nil
	default:
		return nil, errors.New("not Received Docker Vim Instance");
	}
}

//Unmarshal the json raw message to the right struct type for the Vim Instance types
func GetVimInstance(jsonArg json.RawMessage, argValue map[string]interface{}) interface{} {
	if argValue["type"] == "docker" {
		var ret = &catalogue.DockerVimInstance{}
		json.Unmarshal(jsonArg, ret)
		return ret
	} else if argValue["type"] == "openstack" {
		var ret = &catalogue.OpenstackVimInstance{}
		json.Unmarshal(jsonArg, ret)
		return ret
	} else if argValue["type"] == "kubernetes" {
		var ret = &catalogue.KubernetesVimInstance{}
		json.Unmarshal(jsonArg, ret)
		return ret
	} else {
		var ret = &catalogue.BaseVimInstance{}
		json.Unmarshal(jsonArg, ret)
		return ret
	}
}

//Unmarshal the json raw message to the right struct type for the Network and Image types
func GetConcrete(jsonArg json.RawMessage, destType interface{}) reflect.Value {
	switch destType.(type) {
	case catalogue.DockerNetwork:
		ret := &catalogue.DockerNetwork{}
		json.Unmarshal(jsonArg, ret)
		return reflect.ValueOf(ret)
	case catalogue.DockerImage:
		ret := &catalogue.DockerImage{}
		json.Unmarshal(jsonArg, ret)
		return reflect.ValueOf(ret)
	case catalogue.BaseNfvImage:
		ret := &catalogue.BaseNfvImage{}
		json.Unmarshal(jsonArg, ret)
		return reflect.ValueOf(ret)
	case catalogue.BaseNetwork:
		ret := &catalogue.BaseNetwork{}
		json.Unmarshal(jsonArg, ret)
		return reflect.ValueOf(ret)
	default:
		return reflect.Value{}
	}
}
