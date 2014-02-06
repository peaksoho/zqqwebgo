// +build !go1.1

package zqqjsongo

import (
	"encoding/json"
	"errors"
	"regexp"
)

// Implements the json.Unmarshaler interface.
func (j *Json) UnmarshalJSON(p []byte) error {
	//filter comments start
	//peaksoho@163.com
	re, err := regexp.Compile("/\\*.*\\*/")
	p = re.ReplaceAll(p, []byte(""))
	re, err = regexp.Compile("//.*\\n")
	p = re.ReplaceAll(p, []byte("\n"))
	if err != nil {
		return err
	}
	//filter comments end
	return json.Unmarshal(p, &j.data)
}

// Float64 type asserts to `float64`
func (j *Json) Float64() (float64, error) {
	if i, ok := (j.data).(float64); ok {
		return i, nil
	}
	return -1, errors.New("type assertion to float64 failed")
}

// Int type asserts to `float64` then converts to `int`
func (j *Json) Int() (int, error) {
	if f, ok := (j.data).(float64); ok {
		return int(f), nil
	}
	return -1, errors.New("type assertion to float64 failed")
}

// Int type asserts to `float64` then converts to `int64`
func (j *Json) Int64() (int64, error) {
	if f, ok := (j.data).(float64); ok {
		return int64(f), nil
	}
	return -1, errors.New("type assertion to float64 failed")
}
