package scrapper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/vincenciusgeraldo/fetcher/internal/compressor"
	"github.com/vincenciusgeraldo/fetcher/internal/downloader"
	"github.com/vincenciusgeraldo/fetcher/internal/normalizer"
)

const MAX_ROUTINES = 3

type scrapeResponse struct {
	filename  string
	site      string
	numLinks  int
	images    int
	lastFetch time.Time
}

func SaveWebPages(urls []string, meta bool, local bool) {
	var wg sync.WaitGroup
	reqChan := make(chan string)
	resChan := make(chan scrapeResponse)

	for i := 0; i < MAX_ROUTINES; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			scrapeWebsite(reqChan, resChan, local)
		}()
	}

	for _, url := range urls {
		reqChan <- url
	}
	close(reqChan)

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for res := range resChan {
		if meta {
			fmt.Println("====================================")
			fmt.Println("site : " + res.site)
			fmt.Println(fmt.Sprintf("num_links : %d", res.numLinks))
			fmt.Println(fmt.Sprintf("images : %d", res.images))
			fmt.Println("last_fetch : " + res.lastFetch.Local().Format(time.RFC1123))
			fmt.Println("====================================")
		}
	}
}

func scrapeWebsite(urls chan string, resp chan scrapeResponse, local bool) {
	for url := range urls {
		col := colly.NewCollector()

		scrapeCallback(col, url, resp, local)

		col.Visit(url)
	}
}

func scrapeCallback(col *colly.Collector, url string, resp chan scrapeResponse, local bool) {
	var host, path, filename string
	links, images := 0, 0
	col.OnResponse(func(r *colly.Response) {
		host = r.Request.URL.Host
		path = r.Request.URL.Path
		filename = host + "_" + strings.ReplaceAll(path, "/", "_")
		if path == "" || path == "/" {
			filename = host
		}

		if local {
			if err := os.MkdirAll(filename, 0770); err != nil {
				log.Fatal(err)
			}
		}
	})

	col.OnHTML("a[href]", func(h *colly.HTMLElement) {
		links += 1
	})

	var assetUrls []string
	col.OnHTML("img[src]", func(h *colly.HTMLElement) {
		images += 1
		if local {
			assetUrls = append(assetUrls, scrapeImage(h, filename, host)...)
		}
	})

	col.OnHTML("link[href]", func(h *colly.HTMLElement) {
		if local {
			assetUrls = append(assetUrls, scrapeCSS(h, filename, host)...)
		}
	})

	col.OnHTML("script[src]", func(h *colly.HTMLElement) {
		if local {
			assetUrls = append(assetUrls, scrapeJS(h, filename, host)...)
		}
	})

	col.OnScraped(func(r *colly.Response) {
		res := scrapeResponse{
			site:     filepath.Join(host, path),
			filename: filename,
			numLinks: links,
			images:   images,
		}
		parseResponse(r, res, resp, local)

		if local {
			normalizer.NormalizeAssetUrl(filename, assetUrls)
			compressor.CompressFiles(filename)
			os.RemoveAll(filename)
		}
	})
}

func scrapeImage(h *colly.HTMLElement, folder string, host string) []string {
	paths := []string{}
	url := h.Attr("src")
	if !strings.Contains(url, "data") && !strings.Contains(url, "https") {
		paths = append(paths, url)
	}

	imgUrls := strings.Split(h.Attr("srcset"), ",")
	for _, imgUrl := range imgUrls {
		if imgUrl == "" || strings.Contains(imgUrl, "data") {
			continue
		}
		path := strings.Split(imgUrl, " ")
		url = path[0]
		if len(path) >= 2 {
			url = path[len(path)-2]
		}
		sanitizedPath := strings.ReplaceAll(url, "\n", "")
		paths = append(paths, strings.ReplaceAll(sanitizedPath, " ", ""))
	}

	for _, path := range paths {
		if err := downloader.DownloadAsset(host, folder, path); err != nil {
			fmt.Println(err)
		}
	}
	return paths
}

func scrapeCSS(h *colly.HTMLElement, folder, host string) []string {
	paths := []string{}
	url := h.Attr("href")
	if !strings.Contains(url, "https") {
		paths = append(paths, url)
	}

	for _, path := range paths {
		if err := downloader.DownloadAsset(host, folder, path); err != nil {
			fmt.Println(err)
		}
	}

	return paths
}

func scrapeJS(h *colly.HTMLElement, folder, host string) []string {
	paths := []string{}
	url := h.Attr("src")
	if !strings.Contains(url, "https") {
		paths = append(paths, url)
	}

	for _, path := range paths {
		if err := downloader.DownloadAsset(host, folder, path); err != nil {
			fmt.Println(err)
		}
	}

	return paths
}

func parseResponse(r *colly.Response, res scrapeResponse, resp chan scrapeResponse, local bool) {
	name := res.filename + ".html"
	compressedName := res.filename + ".zip"
	if local {
		name = filepath.Join(res.filename, name)
	}

	var f *os.File
	var err error
	if local {
		f, err = os.Open(compressedName)
	} else {
		f, err = os.Open(name)
	}

	if err != nil {
		if f, err = os.Create(name); err != nil {
			log.Fatal(err)
		}
		res.lastFetch = time.Now()
	} else {
		stat, _ := f.Stat()
		res.lastFetch = stat.ModTime()
	}
	defer f.Close()

	if err := r.Save(name); err != nil {
		log.Fatal(err)
	}

	resp <- res
}