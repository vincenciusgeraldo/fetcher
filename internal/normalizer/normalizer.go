package normalizer

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func NormalizeAssetUrl(filename string, urls []string) {
	filename = filepath.Join(filename, filename+".html")
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	output := input
	replaced := map[string]bool{}
	for _, url := range urls {
		assetPath := "assets" + url
		if strings.Contains(url, "https://") {
			assetPath = "assets/" + url
		}

		if _, ok := replaced[url]; !ok {
			output = bytes.Replace(output, []byte(url), []byte(assetPath), -1)
			replaced[url] = true
		}
	}

	if err = ioutil.WriteFile(filename, output, 0666); err != nil {
		log.Fatal(err)
	}
}
