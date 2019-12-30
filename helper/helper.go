package helper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Request(method string, url string, data []byte, header http.Header, statusCode int) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	req.Header = header

	if err != nil {
		return "", err
	}

	resp, _ := client.Do(req)

	if resp.StatusCode != statusCode {
		resp.Body.Close()
		return "", fmt.Errorf("status[%d]:%s", resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
