package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	// "golang.org/x/net/html"

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

func (c Client) GetTips(word string) (TipsList, error) {
	url := fmt.Sprintf("http://wooordhunt.ru/openscripts/forjs/get_tips.php?abc=%s", word)
  resp, err := c.httpClient.Get(url)

  if err != nil {
    return TipsList{}, err
  }

	if resp.StatusCode != http.StatusOK {
		return TipsList{}, errors.New(fmt.Sprintf("Server returned code %d", resp.StatusCode))
	}

	var tips TipsList
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tips)
	return tips, err
}

func (c Client) GetWord(word string) (WordInfo, error) {
	url := fmt.Sprintf("http://wooordhunt.ru/word/%s", word)
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return WordInfo{}, err
	}

	// res := WordInfo{}
	// tkn := html.NewTokenizer(resp.Body)

	// for {
	// 	tt := tkn.Next()

	// 	switch {
	// 	case tt == html.ErrorToken:
	// 	case tt == html.StartTagToken:
	// 	case tt == html.TextToken:
	// 	}
	// }

	body, _ := io.ReadAll(resp.Body)
	return WordInfo{Text: string(body)}, nil
}
