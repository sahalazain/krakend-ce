package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func post(address, path string, data, res interface{}) error {
	o, err := json.Marshal(data)
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(o)

	url := address + path

	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	rdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return errors.New(string(rdata))
	}

	return json.Unmarshal(rdata, res)
}
