package demux

type EventStruct map[string]any
type RawEvent EventStruct

type RawResourceContext EventStruct

func (e EventStruct) HasAttribute(name string) bool {
	_, ok := e[name]
	return ok
}

func (e EventStruct) StringAttributeValue(name string) *string {
	v, ok := e[name]
	if !ok {
		return nil
	}

	vStr, ok := v.(string)
	if !ok {
		return nil
	}

	return &vStr
}

func (e EventStruct) StringAttributeMatches(name string, value string) bool {
	v := e.StringAttributeValue(name)
	if v == nil {
		return false
	}

	return *v == value
}
