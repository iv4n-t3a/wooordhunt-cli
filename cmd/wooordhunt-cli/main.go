package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
)

var cli struct {
	Word     string `arg:""`
	Insecure bool   `short:"i" help:"disable ssl verification"`
	Verbose  bool   `short:"v"`
	RawPrint bool   `short:"r" help:"disable colored and formated output"`
}

type Response struct {
	Tips []struct {
		Word string `json:"w"`
		Tips string `json:"t"`
	} `json:"tips"`
}

func RawPrint(r Response) {
	for i := range r.Tips {
		fmt.Printf("%s: %s\n", r.Tips[i].Word, r.Tips[i].Tips)
	}
}

func PrettifyedPrint(r Response) {
	maxlen := 0
	for i := range r.Tips {
		maxlen = max(maxlen, len(r.Tips[i].Word))
	}

	cyan := color.New(color.FgCyan).SprintFunc()

	for i := range r.Tips {
		word := r.Tips[i].Word
		tips := r.Tips[i].Tips
		spacesCount := maxlen - len(word) + 1
		spaces := strings.Repeat(" ", spacesCount)
		fmt.Printf("%s: %s%s\n", cyan(word), spaces, tips)
	}
}

func main() {
	kong.Parse(&cli)

	if cli.Verbose {
		log.Print("Cli arguments are parsed")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cli.Insecure},
	}
	client := &http.Client{Transport: tr}

	if cli.Verbose {
		log.Print("Client is created")
	}

	url := fmt.Sprintf("http://wooordhunt.ru/openscripts/forjs/get_tips.php?abc=%s", cli.Word)

	if cli.Verbose {
		log.Printf("Url is created: %s", url)
	}

	resp, err := client.Get(url)

	if cli.Verbose {
		log.Print("Got response from site: ", resp)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server returned code %d", resp.StatusCode)
	}
	if cli.Verbose {
		log.Printf("Server returned code %d", resp.StatusCode)
	}

	var tips Response
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tips)

	if err != nil {
		log.Fatal(err)
	}
	if cli.Verbose {
		log.Print("Server response json is parsed: ", tips)
	}

	if cli.RawPrint {
		RawPrint(tips)
	} else {
		PrettifyedPrint(tips)
  }
}
