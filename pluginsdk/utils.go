package pluginsdk

import (
	"github.com/openbaton/go-openbaton/catalogue"
	"errors"
	"encoding/json"
)

// Obtain a Docker vim instance from a interfaced struct
func GetDockerVimInstance(vimInstance interface{}) (*catalogue.DockerVimInstance, error) {
	switch t := vimInstance.(type) {
	case *catalogue.DockerVimInstance:
		return t, nil
	default:
		return nil, errors.New("Not Received Docker Vim Instance");
	}
}

// Obtain a Openstack vim instance from a interfaced struct
func GetOpenstackVimInstance(vimInstance interface{}) (*catalogue.OpenstackVimInstance, error) {
	switch t := vimInstance.(type) {
	case *catalogue.OpenstackVimInstance:
		return t, nil
	default:
		return nil, errors.New("Not Received Docker Vim Instance");
	}
}

// Obtain a Generic vim instance from a interfaced struct
func GetBaseVimInstance(vimInstance interface{}) (*catalogue.BaseVimInstance, error) {
	switch t := vimInstance.(type) {
	case *catalogue.BaseVimInstance:
		return t, nil
	default:
		return nil, errors.New("Not Received Docker Vim Instance");
	}
}

func GetVimInstance(jsonArg json.RawMessage, argValue map[string]interface{}) interface{} {
	if argValue["type"] == "docker" {
		var ret = &catalogue.DockerVimInstance{}
		json.Unmarshal(jsonArg, ret)
		return ret
	} else if argValue["type"] == "openstack" {
		var ret = &catalogue.OpenstackVimInstance{}
		json.Unmarshal(jsonArg, ret)
		return ret
	} else {
		var ret = &catalogue.BaseVimInstance{}
		json.Unmarshal(jsonArg, ret)
		return ret
	}
}
