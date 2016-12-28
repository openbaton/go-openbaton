package config

import (
	"strings"
)

type Properties map[string]interface{}

func (p Properties) Section(key string) (section Properties, ok bool) {
	if val, ok := p.Value(key, nil); ok {
		ret, ok := val.(Properties)

		return ret, ok
	}

	return nil, false
}

func (p Properties) Value(key string, fb interface{}) (interface{}, bool) {
	keys := stack(strings.Split(key, "."))
	current := p

	for len(keys) > 1 {
		key := keys.Pop()

		iface, ok := current.Value(key, fb)
		if !ok {
			return fb, false
		}

		subMap, ok := iface.(map[string]interface{})
		if !ok {
			return fb, false
		}

		current = Properties(subMap)
	}

	if ret, ok := current[keys[0]]; ok {
		return ret, true
	}

	return fb, false
}

// ValueBool returns a key as a boolean value.
func (p Properties) ValueBool(key string, fb bool) (value, ok bool) {
	if val, ok := p.Value(key, nil); ok {
		ret, ok := val.(bool)

		return ret, ok
	}

	return fb, false
}

func (p Properties) ValueInt(key string, fb int) (value int, ok bool) {
	if val, ok := p.Value(key, nil); ok {
		ret, ok := val.(int)

		return ret, ok
	}

	return fb, false
}

func (p Properties) ValueString(key string, fb string) (value string, ok bool) {
	if val, ok := p.Value(key, nil); ok {
		ret, ok := val.(string)

		return ret, ok
	}

	return fb, false
}

type stack []string

func (s *stack) Empty() bool {
	return len(*s) == 0
}

func (s *stack) Pop() string {
	if s.Empty() {
		return ""
	}

	ret := (*s)[0]
	*s = (*s)[1:]

	return ret
}
