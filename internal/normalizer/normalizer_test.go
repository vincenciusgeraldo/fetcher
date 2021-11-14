package normalizer_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vincenciusgeraldo/fetcher/internal/normalizer"
)

func TestNormalizeAssetUrl(t *testing.T) {
	tests := []struct {
		name     string
		urls     []string
		expected string
	}{
		{
			name:     "success replace urls",
			urls:     []string{"/test/a"},
			expected: "<img src=\"assets/test/a\"/><img src=\"/test/b\"/><img src=\"/test/c\"/>",
		},
		{
			name:     "success replace urls",
			urls:     []string{"/test/a", "/test/c"},
			expected: "<img src=\"assets/test/a\"/><img src=\"/test/b\"/><img src=\"assets/test/c\"/>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := []byte("<img src=\"/test/a\"/><img src=\"/test/b\"/><img src=\"/test/c\"/>")
			os.Mkdir("www.test.com", 0770)
			os.WriteFile("www.test.com/www.test.com.html", temp, 0644)

			normalizer.NormalizeAssetUrl("www.test.com", tt.urls)
			data, _ := os.ReadFile("www.test.com/www.test.com.html")

			assert.Equal(t, tt.expected, string(data))
			os.RemoveAll("./www.test.com/")
		})
	}
}
