package logkOption

import "time"

// GetString is helper to retrieve string value in Values by key
// Value type must be exact, as it use casting instead of converting to target value
func GetString(o *Options, k string) (string, bool) {
	s, ok := o.Values[k].(string)
	if !ok {
		return "", false
	}
	return s, true
}

func GetInt64(o *Options, k string) (int64, bool) {
	i, ok := o.Values[k].(int64)
	if !ok {
		return 0, false
	}
	return i, true
}

func GetTime(o *Options, k string) (time.Time, bool) {
	t, ok := o.Values[k].(time.Time)
	if !ok {
		return time.Time{}, false
	}
	return t, true
}

func GetError(o *Options, k string) error {
	e, ok := o.Values[k].(error)
	if !ok {
		return nil
	}
	return e
}
