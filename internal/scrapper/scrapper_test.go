package scrapper_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincenciusgeraldo/fetcher/internal/scrapper"
)

func TestSaveWebPages(t *testing.T) {
	tests := []struct {
		name          string
		input         []string
		metadata      bool
		local         bool
		expectedHosts []string
		expectedMeta  []string
	}{
		{
			name:          "save web pages with www",
			input:         []string{"https://www.google.com"},
			expectedHosts: []string{"www.google.com"},
		},
		{
			name:          "save web pages without www",
			input:         []string{"https://google.com"},
			expectedHosts: []string{"www.google.com"},
		},
		{
			name:          "save with metadata",
			metadata:      true,
			input:         []string{"https://google.com"},
			expectedHosts: []string{"www.google.com"},
			expectedMeta:  []string{"site : www.google.com\nnum_links : 18\nimages : 2\n"},
		},
		{
			name:          "save multiple pages",
			metadata:      true,
			input:         []string{"https://google.com", "https://autify.com"},
			expectedHosts: []string{"www.google.com", "autify.com"},
			expectedMeta: []string{
				"site : autify.com\nnum_links : 54\nimages : 96",
				"site : www.google.com\nnum_links : 18\nimages : 2",
			},
		},
		{
			name:          "save pages locally with asset",
			local:         true,
			input:         []string{"https://google.com"},
			expectedHosts: []string{"www.google.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := captureOutput(func() {
				scrapper.SaveWebPages(tt.input, tt.metadata, tt.local)
			})

			for _, meta := range tt.expectedMeta {
				assert.Contains(t, out, meta)
			}

			for _, host := range tt.expectedHosts {
				if tt.local {
					assert.DirExists(t, filepath.Join("./", host, "/"))
					assert.FileExists(t, filepath.Join("./", host, host+".html"))
					os.RemoveAll(filepath.Join("./", host, "/"))
				} else {
					assert.FileExists(t, filepath.Join("./", host+".html"))
					os.RemoveAll(filepath.Join("./", host+".html"))
				}
			}
		})
	}
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()

	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
