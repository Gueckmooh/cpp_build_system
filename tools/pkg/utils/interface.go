package utils

func IsBool(v interface{}) bool {
	switch v.(type) {
	case bool:
		return true
	default:
		return false
	}
}
