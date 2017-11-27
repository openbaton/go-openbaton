package pluginsdk

import (
	"github.com/openbaton/go-openbaton/catalogue"
	"errors"
	"encoding/json"
)

func GetDockerVimInstance(vimInstance interface{}) (*catalogue.DockerVimInstance, error) {
	switch t := vimInstance.(type) {
	case *catalogue.DockerVimInstance:
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
