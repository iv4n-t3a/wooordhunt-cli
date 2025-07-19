package client_test

import (
	"testing"

	"github.com/iv4n-t3a/wooordhunt-cli/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestDetectLang(t *testing.T) {
	ruString := "йцуке"
	lang := client.DetectLang(ruString)
  assert.Equal(t, lang, client.Ru, "Ru string not recognized")

  enString := "qwert"
	lang = client.DetectLang(enString)
  assert.Equal(t, lang, client.En, "En not recognized")
}
