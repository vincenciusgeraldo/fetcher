package downloader

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadAsset(host string, folder string, url string) error {
	assetPath := filepath.Join("./", folder, "assets", url)
	host = "https://" + host
	assetUrl := host + url

	response, err := http.Get(assetUrl)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	if err := os.MkdirAll(filepath.Dir(assetPath), 0770); err != nil {
		log.Fatal(err)
	}
	file, _ := os.Create(assetPath)
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}