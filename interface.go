package goutil

func GetInterfaceSlice(value interface{}) []interface{} {
	var islice []interface{}
	switch t := value.(type) {
	case string, float64, bool:
		return nil
	case []string:
		islice = make([]interface{}, len(t))
		for i := range t {
			islice[i] = t[i]
		}
	case []float64:
		islice = make([]interface{}, len(t))
		for i := range t {
			islice[i] = t[i]
		}
	case []bool:
		islice = make([]interface{}, len(t))
		for i := range t {
			islice[i] = t[i]
		}
	case []interface{}:
		for i := range t {
			switch tv := t[i].(type) {
			case string:
				t[i] = tv
			}
		}
		islice = t
	}
	return islice
}
