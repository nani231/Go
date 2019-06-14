package watcher

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

const (
	FTP_HOST      = "localhost"
	FTP_PORT      = "21"
	FTP_USER      = "narendra"
	FTP_PASSWORD  = "narendra"
	FTP_DIRECOTRY = "E://watch//"
)

var localPath = "E://temp//"

func fileTransfer(input, file string) {
	client, err := ftp.Dial(fmt.Sprintf("%s:%s", FTP_HOST, FTP_PORT))
	if err != nil {
		fmt.Println(err)
	}

	// Log in the server
	err = client.Login(FTP_USER, FTP_PASSWORD)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := client.Retr(input)
	if err != nil {
		fmt.Println(err)
	}

	// Read the file
	var srcFile []byte
	_, err = reader.Read(srcFile)
	if err != nil {
		fmt.Println(err)
	}

	// Create the destination file
	dstFile, err := os.Create(localPath + file)
	if err != nil {
		fmt.Println(err)
	}
	defer dstFile.Close()

	// Copy the file
	_, err = io.Copy(dstFile, reader)
	if err != nil {
		fmt.Println(err)
	}
	dstFile.Close()
	fmt.Println(input, " File transfered successfully")

	// Delete Client File
	delFile := client.Delete(input)
	_ = client.Delete(strings.Replace(input, ".zip", "_COMPLETE.TXT", -1))
	fmt.Println(input, "Deleted with", delFile)
	fmt.Println("Signal File Deleted ")
}

func fileStatus() map[int]string {
	client, err := ftp.Dial(fmt.Sprintf("%s:%s", FTP_HOST, FTP_PORT))
	err = client.Login(FTP_USER, FTP_PASSWORD)
	if err != nil {
		fmt.Println(err)
	}
	// Log in the server
	err = client.Login(FTP_USER, FTP_PASSWORD)
	if err != nil {
		fmt.Println(err)
	}

	reader, _ := client.NameList(FTP_DIRECOTRY)

	err1 := client.Quit()
	if err1 != nil {
		fmt.Println(err1)
	}

	return find(reader, "COMPLETE")
}

func find(a []string, sub string) map[int]string {
	arr := make(map[int]string)
	for i := 0; i < len(a); i++ {
		if strings.Contains(a[i], sub) {
			arr[i] = a[i]
		}
	}

	return arr
}

func foreverloop() {
	for {
		list := fileStatus()
		if len(list) > 0 {
			for _, value := range list {
				filePath := strings.Replace(value, "_COMPLETE.TXT", ".zip", -1)
				fileArray := strings.Split(filePath, "/")
				fileName := fileArray[len(fileArray)-1]
				fileTransfer(filePath, fileName)
			}
		} else {
			fmt.Println("Waiting for files")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	go foreverloop()
	fmt.Scanln()
}
