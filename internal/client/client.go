package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/anaskhan96/soup"

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
	resp, err := soup.GetWithClient(url, &c.httpClient)

	if err != nil {
		return WordInfo{}, err
	}

	doc := soup.HTMLParse(resp)

	return WordInfo{
		Word:         word,
		Meaning:      doc.Find("div", "class", "t_inline_en").FullText(),
		WordType:     doc.Find("h4", "class", "pos_item").FullText(),
		Phrases:      ParseList(doc.Find("div", "class", "phrases").HTML()),
		SimilarWords: ParseList(doc.Find("div", "class", "similar_words").HTML()),
		WordForms:    ParseList(doc.Find("div", "class", "word_form_block").HTML()),
	}, nil
}

func ParseList(list string) (res []string) {
	listHTML := strings.Split(list, "<br/>")

	for i := range listHTML {
		data := soup.HTMLParse(listHTML[i]).FullText()
		data = strings.TrimSpace(data)

		if len(data) != 0 {
			res = append(res, data)
		}
	}

	return
}
