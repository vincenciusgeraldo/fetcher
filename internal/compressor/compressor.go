package compressor

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CompressFiles(filename string) {
	dest, err := os.Create(filename + ".zip")
    if err != nil {
        log.Fatal(err)
    }
    arch := zip.NewWriter(dest)
    err = filepath.Walk(filename, func(filePath string, info os.FileInfo, err error) error {
        if info.IsDir() {
            return nil
        }
        if err != nil {
            return err
        }
        relPath := strings.TrimPrefix(filePath, filepath.Dir(filename))
        zipFile, err := arch.Create(relPath)
        if err != nil {
            return err
        }
        fsFile, err := os.Open(filePath)
        if err != nil {
            return err
        }
        _, err = io.Copy(zipFile, fsFile)
        if err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
    err = arch.Close()
    if err != nil {
        log.Fatal(err)
    }
}
