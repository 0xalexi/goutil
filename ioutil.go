package goutil

import "io/ioutil"

func ReadFileToString(fpath string) (string, error) {
	var buf []byte
	var err error
	buf, err = ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
