package main

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"test/auto"
	"test/logger"
	"time"
)

type myService struct {
	stopChan chan struct{}
}

func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	go m.main(changes)
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				close(m.stopChan)
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				log.Printf("unexpected control request #%d", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func (m *myService) main(changes chan<- svc.Status) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://tinypng.com/backend/opt/shrink", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("x-forwarded-for", generateIp())

	// 默认的文件夹路径
	defaultPath := "C:\\Users\\83795\\Downloads"
	args := os.Args

	folderPath := defaultPath
	if len(args) > 1 {
		// 如果传入了参数，使用传入的文件夹路径
		//如果参数是--test则使用默认路径
		if args[1] != "--test" {
			folderPath = args[1]
		}

	}
	// Create a logs directory in the folderPath
	logsPath := filepath.Join(folderPath, "logs")
	if _, err := os.Stat(logsPath); os.IsNotExist(err) {
		err := os.MkdirAll(logsPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}
	}

	// Create a log file in the logs directory
	logFile, err := os.OpenFile(filepath.Join(logsPath, "log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	logger.InitLogger(logFile)

	files, err := auto.WatchFolder(folderPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("files:", files, "folderPath:", folderPath)
	logger.Logger.Println("files:", files, "folderPath:", folderPath)

	for {
		select {
		case <-m.stopChan:
			// 收到停止信号，结束运行
			return
		default:
			for filePath := range files {
				// 如果是图片类型的才往下执行 png jpg jpeg 使用迭代器
				if !auto.IsImage(filePath) {
					continue
				}
				// 路径不存在则跳过
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					continue
				}
				logger.Logger.Println("Modified file:", filePath)
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
					logger.Logger.Println("Error deleting file:", err)
					continue
				}
			}
		}
	}
}

func compress(location string, absPath string) {
	fmt.Println("Absolute path:", absPath)
	logger.Logger.Println("Absolute path:", absPath)
	imgObj := map[string]string{"url": location, "fileName": absPath, "absPath": absPath}
	urls := []map[string]string{imgObj}
	fmt.Println(location)
	logger.Logger.Println(location)
	auto.Downloaded(urls)
}

func generateIp() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--test" {
		println("Running in test mode")
		m := &myService{stopChan: make(chan struct{})}
		m.main(nil)
	} else {
		err := svc.Run("CompressImg", &myService{stopChan: make(chan struct{})})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Service stopped")
	}

}
