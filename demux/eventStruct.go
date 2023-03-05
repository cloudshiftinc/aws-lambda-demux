package demux

func HasAttribute(e map[string]any, name string) bool {
	_, ok := e[name]
	return ok
}

func StringAttributeValue(e map[string]any, name string) *string {
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

func StringAttributeMatches(e map[string]any, name string, value string) bool {
	v := StringAttributeValue(e, name)
	if v == nil {
		return false
	}

	return *v == value
}
