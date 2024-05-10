package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"test/auto"
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

	filePath := "./png/1.png"
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	req.Body = ioutil.NopCloser(bytes.NewReader(fileData))

	resp, err := client.Do(req)
	//print res
	fmt.Println(resp)
	if err != nil {
		panic(err)
	}
	location := resp.Header.Get("Location")
	//把图片的原名和location保存到结构体里 结构体叫imgObj 有url 何 fileName两个字段
	imgObj := map[string]string{"url": location, "fileName": "1.png"}
	urls := []map[string]string{imgObj}

	auto.Downloaded(urls)

	fmt.Println(location)

	defer resp.Body.Close()
	// handle resp here
}
