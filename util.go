package main

import (
	"encoding/json"
	"github.com/ddliu/go-httpclient"
	"io/ioutil"
	"os/user"
	"path"
)

func ReadBody(res *httpclient.Response) (map[string]interface{}, error) {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	// Declared an empty interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	err = json.Unmarshal([]byte(bodyString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ReadArrBody(res *httpclient.Response) ([]map[string]interface{}, error) {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)
	// Declared an empty interface
	var result []map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	err = json.Unmarshal([]byte(bodyString), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func LoadConfig() (*Configuration, error) {

	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	file, err := ioutil.ReadFile(path.Join(u.HomeDir, ".gitlo"))
	if err != nil {
		return nil, err
	}

	data := Configuration{}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func SaveConfig(configuration Configuration) error {
	u, err := user.Current()
	if err != nil {
		return err
	}

	file, err := json.MarshalIndent(configuration, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(u.HomeDir, ".gitlo"), file, 0644)
}

