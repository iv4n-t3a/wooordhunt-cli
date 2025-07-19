package client

import "strings"

type Lang int

const (
	En Lang = iota
	Ru
)

const enAlphabet = "abcdefghiklmnopqrstvxyz"

func DetectLang(s string) Lang {
	for i := range s {
		if strings.Contains(enAlphabet, string(s[i])) {
			return En
		}
	}
	return Ru
}
