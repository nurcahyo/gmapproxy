package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func UnmarshalReader(input io.ReadCloser, cast interface{}) error {
	body, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, cast)
}
