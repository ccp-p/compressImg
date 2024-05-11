package auto

import (
	"log"
	"testing"
)

// watchFolder
func TestWatchFolder(t *testing.T) {

	files, err := WatchFolder("D:\\project\\my_go_project\\test\\png")
	if err != nil {
		log.Fatal(err)
	}

	for filePath := range files {
		log.Println("Modified file:", filePath)
	}
}
