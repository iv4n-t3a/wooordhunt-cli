package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/iv4n-t3a/wooordhunt-cli/config"
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

func (c Client) GetTips(word string) (*TipsList, error) {
	url := fmt.Sprintf("http://wooordhunt.ru/openscripts/forjs/get_tips.php?abc=%s", word)
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Server returned code %d", resp.StatusCode))
	}

	var tips TipsList
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tips)
	return &tips, err
}

func (c Client) GetWord(word string) (*WordInfo, error) {
	url := fmt.Sprintf("http://wooordhunt.ru/word/%s", word)
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Server returned code %d", resp.StatusCode))
	}

	doc, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	res := HtmlToWordInfo(string(doc), word)
	return &res, nil
}
