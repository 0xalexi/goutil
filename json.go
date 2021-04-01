package goutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func StructToJsonFile(v interface{}, fpath string) error {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fpath, b, 0644)
}

func JsonFileToStruct(fpath string, interfacePointer interface{}) error {
	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(file), interfacePointer)
}

func PrettySprint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprint("Json error:", err)
	}
	return string(b)
}
