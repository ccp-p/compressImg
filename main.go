package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"test/auto"
	"time"
)

func generateIp() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://tinypng.com/backend/opt/shrink", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-forwarded-for", generateIp())

	files, err := auto.WatchFolder("D:\\project\\my_go_project\\test\\png")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("files:", files)

	for filePath := range files {
		// 如果是图片类型的才往下执行 png jpg jpeg 使用迭代器
		if !auto.IsImage(filePath) {
			continue
		}
		// 路径不存在则跳过
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}

		log.Println("Modified file:", filePath)
		// 添加延迟
		time.Sleep(1 * time.Second)
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(fileData))

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		location := resp.Header.Get("Location")
		//把图片的原名和location保存到结构体里 结构体叫imgObj 有url 何 fileName两个字段
		absPath, err := filepath.Abs(filePath)
		compress(location, absPath)

		err = resp.Body.Close()
		if err != nil {
			return
		}
		// 删除源文件
		err = os.Remove(filePath)
		if err != nil {
			log.Println("Error deleting file:", err)
			continue
		}

	}

	// handle resp here
}
func compress(location string, absPath string) {

	fmt.Println("Absolute path:", absPath)
	imgObj := map[string]string{"url": location, "fileName": absPath, "absPath": absPath}
	urls := []map[string]string{imgObj}
	fmt.Println(location)

	auto.Downloaded(urls)

}
