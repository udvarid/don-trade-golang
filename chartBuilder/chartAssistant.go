package chart

import (
	"log"
	"os"
	"path/filepath"
)

func DeleteHtml() {
	folderPath := "./html"
	filePattern := "kline-*.html"

	fullPattern := filepath.Join(folderPath, filePattern)

	files, err := filepath.Glob(fullPattern)
	if err != nil {
		log.Fatalf("Error finding files: %v", err)
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			log.Printf("Error deleting file %s: %v", file, err)
		}
	}
}

func WaitUntilHtmlReady(sessionId string) {
	filePath := "./html/kline-" + sessionId + ".html"
	for {
		if _, err := os.Stat(filePath); err == nil {
			break
		}
	}
}

func DeleteSpecificHtml(sessionId string) {
	filePath := "./html/kline-" + sessionId + ".html"
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Error deleting file %s: %v", filePath, err)
	}
}
