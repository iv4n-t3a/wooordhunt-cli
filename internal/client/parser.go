package client

import (
	"github.com/anaskhan96/soup"
	"strings"
)

func htmlToWordInfo(html string, word string) WordInfo {
	doc := soup.HTMLParse(html)

	return WordInfo{
		Word:         word,
		Meaning:      fullTextOrNil(doc, "div", "class", "t_inline_en"),
		WordType:     fullTextOrNil(doc, "h4", "class", "pos_item"),
		Phrases:      parseHtmlList(htmlOrNil(doc, "div", "class", "phrases")),
		SimilarWords: parseHtmlList(htmlOrNil(doc, "div", "class", "similar_words")),
		WordForms:    parseHtmlList(htmlOrNil(doc, "div", "class", "word_form_block")),
	}
}

func parseHtmlList(list *string) (res []string) {
	if list == nil {
		return nil
	}

	listHTML := strings.Split(*list, "<br/>")

	for i := range listHTML {
		dataRoot := soup.HTMLParse(listHTML[i])
		if dataRoot.Error != nil {
      continue
		}
		data := dataRoot.FullText()
		data = strings.TrimSpace(data)

		if len(data) != 0 {
			res = append(res, data)
		}
	}

	return
}

func fullTextOrNil(doc soup.Root, args ...string) *string {
	root := doc.Find(args...)
	if root.Error != nil {
		return nil
	}
	res := root.FullText()
	return &res
}

func htmlOrNil(doc soup.Root, args ...string) *string {
	root := doc.Find(args...)
	if root.Error != nil {
		return nil
	}
	res := root.HTML()
	return &res
}
