package main

import (
	"net/http"
	// "text/template"
	"strconv"
	"os"
	"bufio"
	"encoding/json"
)

type SaveFormat struct {
	Id string `json:"id"`
	Time int64 `json:"time"`
	Name string `json:"name"`
	Content string `json:"content"`
}

type Resp struct {
	Success bool `json:"success"`
}

type ResError struct {
	Success bool `json:"success"`
	Msg string `json:"msg"`
}

type List []SaveFormat

type Listp struct {
	Success bool `json:"success"`
	Data List `json:"data"`
}

type Uploadp struct {
	Success bool `json:"success"`
	Data map[int]string `json:"data"`
}

func (s SaveFormat) toString() string {
	return s.Id + Splitter + strconv.FormatInt(s.Time, 10) + Splitter + s.Name + Splitter + s.Content
}

const Splitter = "<<<&&&>>>>>"

// 写入文件index.txt
func writeToFile(filename string, outPut []byte) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666)
	defer f.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	_, err = writer.Write(outPut)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

func getErrorMsg(err error) []byte {
	profile := ResError{
		Success: false,
		Msg: err.Error(),
	}
	msg, err := json.Marshal(profile)
	if err != nil {
		// debugLog.SetPrefix("[Error]")
		// debugLog.Println("访问错误" + err.Error())
		webLogPrintln("[Error]", "访问错误" + err.Error())
		return nil
	}
	return msg
}

func webInterface() {
	http.HandleFunc("/api/save", saveCallback)
	http.HandleFunc("/api/queryList", queryList)
	http.HandleFunc("/api/delete", deleteData)
	// http.HandleFunc("/error", errorCallback)
	http.HandleFunc("/api/upload", upload)
	http.HandleFunc("/uploaded/", showPicHandle)
	// http.HandleFunc("/removeImage", removeImage)
}

func webServer() {
	// fs := http.FileServer(http.Dir(Conf.WebPath))
  // http.Handle("/", fs)
	// http.HandleFunc("/web/", webServerCallback)
	// http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web"))))
	webInterface()
}