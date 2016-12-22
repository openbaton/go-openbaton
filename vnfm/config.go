package vnfm

import (
	"strings"
)

type Properties map[string]interface{} 

func (p Properties) Section(key string) (section Properties, ok bool) {
    if val, ok := p.Value(key); ok {
        return val.(Properties)
    }

    return nil, false
}

func (p Properties) Value(key string) (interface{}, bool) {
    keys := stack(strings.Split(key, "."))
    current := p

    for len(keys) > 1 {
        key := keys.Pop()

        iface, ok := current.Value(key)
        if !ok {
            return nil, false
        }

        subMap, ok := iface.(map[string]interface{})
        if !ok {
            return nil, false
        }

        current = Properties(subMap)
    }

    return current[keys[0]]
}

func (p Properties) ValueBool(key string) (value,ok bool) {
    if val, ok := p.Value(key); ok {
        return val.(bool)
    }

    return nil, false
}

func (p Properties) ValueInt(key string) (value int, ok bool) {
    if val, ok := p.Value(key); ok {
        return val.(int)
    }

    return nil, false
}

func (p Properties) ValueString(key string) (value string, ok bool) {
    if val, ok := p.Value(key); ok {
        return val.(string)
    }

    return nil, false
}

type stack []string

func (s *stack) Empty() bool {
    return len(*s) == 0
}

func (s *stack) Pop() string {
    if s.Empty() {
        return ""
    }

    ret := *s[0]
    *s = *s[1:]

    return ret
}