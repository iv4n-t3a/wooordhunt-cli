package client

import (
	"github.com/anaskhan96/soup"
	"strings"
)

func HtmlToWordInfo(html string, word string) WordInfo {
	doc := soup.HTMLParse(html)

  if DetectLang(word) == En {
    return EnHtmlToWordInfo(doc, word)
  }
  return RuHtmlToWordInfo(doc, word)
}

func EnHtmlToWordInfo(doc soup.Root, word string) WordInfo {
	return WordInfo{
		Word:         word,
		Meaning:      FullTextOrNil(doc, "div", "class", "t_inline_en"),
		WordType:     FullTextOrNil(doc, "h4", "class", "pos_item"),
		Phrases:      ParseHtmlList(HtmlOrNil(doc, "div", "class", "phrases")),
		SimilarWords: ParseHtmlList(HtmlOrNil(doc, "div", "class", "similar_words")),
		WordForms:    ParseHtmlList(HtmlOrNil(doc, "div", "class", "word_form_block")),
	}
}

func RuHtmlToWordInfo(doc soup.Root, word string) WordInfo {
	return WordInfo{
		Word:         word,
		Meaning:      FullTextOrNil(doc, "div", "class", "t_inline_en"),
		WordType:     FullTextOrNil(doc, "h4", "class", "pos_item"),
		Phrases:      ParseHtmlList(HtmlOrNil(doc, "div", "class", "phrases")),
		SimilarWords: ParseHtmlList(HtmlOrNil(doc, "div", "class", "similar_words")),
		WordForms:    ParseHtmlList(HtmlOrNil(doc, "div", "class", "word_form_block")),
	}
}

func ParseHtmlList(list *string) (res []string) {
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

func FullTextOrNil(doc soup.Root, args ...string) *string {
	root := doc.Find(args...)
	if root.Error != nil {
		return nil
	}
	res := root.FullText()
	return &res
}

func HtmlOrNil(doc soup.Root, args ...string) *string {
	root := doc.Find(args...)
	if root.Error != nil {
		return nil
	}
	res := root.HTML()
	return &res
}
