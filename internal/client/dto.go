package client

type Tips struct {
	Word string `json:"w"`
	Tips string `json:"t"`
}

type TipsList struct {
	Tips []Tips `json:"tips"`
}

type WordInfo struct {
	Text string
}
