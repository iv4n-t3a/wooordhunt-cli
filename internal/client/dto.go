package client

type Tips struct {
	Word string `json:"w"`
	Tips string `json:"t"`
}

type TipsList struct {
	Tips []Tips `json:"tips"`
}

type Example struct {
	Example     string
	Translation string
}

type SimilarWord struct {
	Word    string
	Meaning string
}

type WordForm struct {
	word     string
	formType string
}

type WordInfo struct {
  Word string
	Meaning      string
	WordType     string
	Phrases      []string
	Examples     []string
	SimilarWords []string
	WordForms    []string
}
