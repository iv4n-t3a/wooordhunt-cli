package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iv4n-t3a/wooordhunt-cli/config"
)

const (
  GET_TIPS_METHOD = "GET"
  GET_TIPS_ENDPOINT = "http://wooordhunt.ru/openscripts/forjs/get_tips.php?abc=%s"

  GET_WORD_METHOD = "GET"
  GET_WORD_ENDPOINT = "http://wooordhunt.ru/word/%s"
)

type Client struct {
	httpClient http.Client
	conf       config.Config
}

func NewClient(conf config.Config) (Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: conf.Insecure},
	}
	httpClient := &http.Client{Transport: tr}
	return Client{httpClient: *httpClient}, nil
}

func (c Client) makeRequest(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.conf.UserAgent)

	return c.httpClient.Do(req)
}

func (c Client) GetTips(word string) (*TipsList, error) {
	url := fmt.Sprintf(GET_TIPS_ENDPOINT, word)
	resp, err := c.makeRequest(GET_TIPS_METHOD, url, nil)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned code %d", resp.StatusCode)
	}

	var tips TipsList
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tips)
	return &tips, err
}

func (c Client) GetWord(word string) (*WordInfo, error) {
	url := fmt.Sprintf("http://wooordhunt.ru/word/%s", word)
	resp, err := c.makeRequest(GET_WORD_METHOD, url, nil)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server returned code %d", resp.StatusCode)
	}

	doc, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	res := HtmlToWordInfo(string(doc), word)
	return &res, nil
}
