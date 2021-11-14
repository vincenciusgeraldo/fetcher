package downloader_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincenciusgeraldo/fetcher/internal/downloader"
)

func TestNormalizeAssetUrl(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
		err      error
	}{
		{
			name: "success download image",
			url:  "/images/branding/googlelogo/1x/googlelogo_white_background_color_272x92dp.png",
		},
		{
			name: "success download image",
			url:  "/tes/blabla.jpgg",
			err: fmt.Errorf("Received non 200 response code"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := downloader.DownloadAsset("www.google.com", "www.google.com", tt.url)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.FileExists(t, "www.google.com/assets"+tt.url)
			}
			os.RemoveAll("./www.google.com/")
		})
	}
}
