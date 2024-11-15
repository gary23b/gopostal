package gopostal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func MakeRequestWithoutRedirects(req *http.Request, timeout time.Duration) (*Response, *http.Response, error) {
	client := &http.Client{}

	// Block redirects
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return fmt.Errorf("Redirect not allowed: %+v", via)
	}

	client.Timeout = timeout

	httpResp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	// convert to the GoPostal Response type as well. This will read the body as well.
	postalResp, err := NewResponse(httpResp)
	if err != nil {
		return nil, nil, err
	}

	return postalResp, httpResp, nil
}

func EncodeJson(in any) ([]byte, error) {
	jsonData, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func DecodeJson(jsonData []byte, out any) error {
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	err := decoder.Decode(out)
	if err != nil {
		return err
	}

	return nil
}

type Response struct {
	Status     int                     `json:"Status,omitempty"`
	Header     http.Header             `json:"Header,omitempty"`
	Body       []byte                  `json:"Body,omitempty"`
	BodyString string                  `json:"BodyString,omitempty"`
	Cookies    map[string]*http.Cookie `json:"Cookies,omitempty"`
}

func NewResponse(in *http.Response) (*Response, error) {
	if in == nil {
		return nil, fmt.Errorf("input is nil")
	}

	ret := &Response{}
	var err error

	ret.Status = in.StatusCode
	ret.Body, err = io.ReadAll(in.Body)
	in.Body.Close()
	if err != nil {
		return nil, err
	}

	ret.BodyString = string(ret.Body)

	ret.Header = in.Header

	ret.Cookies = make(map[string]*http.Cookie)
	cookies := in.Cookies()
	for _, c := range cookies {
		ret.Cookies[c.Name] = c
	}
	return ret, nil
}

func ReadResponseFromJson(filePath string) (*Response, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	ret := &Response{}
	err = DecodeJson(data, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *Response) SaveToJson(filePath string) error {
	if s == nil {
		return fmt.Errorf("input is nil")
	}

	jsonData, err := EncodeJson(s)
	if err != nil {
		return err
	}

	EnsureDir(filePath)

	err = os.WriteFile(filePath, jsonData, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func EnsureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}
