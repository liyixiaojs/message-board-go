package main

import (
	"io"
	"errors"
	"net/http"
	"strings"
	"path"
	"os"
	"encoding/json"
	"io/ioutil"
	"time"
	"fmt"
	"strconv"
	"math/rand"
	"regexp"
)

// 随机数
func CreateCaptcha() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func saveImageMap(oriFilename string, filename string) {
	f, openFileErr := os.OpenFile("./imageMap.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if openFileErr != nil {
		webLogPrintln("[Error]", "访问错误" + openFileErr.Error())
			return
	}
	defer f.Close()
	f.Write([]byte(oriFilename + Splitter + filename + "\n"))
}

// 上传图片
func upload(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// 接受图片
	req.ParseMultipartForm(32 << 20)
	files := req.MultipartForm.File["image"]
	var data map[int]string
	data = make(map[int]string)
	for i,handler := range files {
		// fileName := handler.Filename
		// fileSize := handler.Size
		file, err := files[i].Open()
		if err!=nil{
			fmt.Println("打开文件失败")
		}
		defer file.Close()
		var time int64 = (time.Now().UnixNano() / 1000000)
		filename := CreateCaptcha()
		ext := strings.ToLower(path.Ext(handler.Filename))
		if ext != ".jpg" && ext != ".png" {
			webLogPrintln("[Error]", "图片格式错误" + err.Error())
			w.Write(getErrorMsg(errors.New("图片格式错误")))
			return
		}
		filename = filename + strconv.FormatInt(time, 10) + ext
		data[i] = "/uploaded/" + filename
		saveImageMap(handler.Filename, filename)
		os.Mkdir("./uploaded/", os.ModePerm)
		cur, err := os.Create("./uploaded/" + filename)
		defer cur.Close()
		if err != nil {
			webLogPrintln("[Error]", "上传失败 " + err.Error())
		}
		io.Copy(cur, file)

		// fmt.Println(filename)
		// tempFile, err := ioutil.TempFile("temp-images", "simple.*.png")
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// defer tempFile.Close()
		// fileBytes, err := ioutil.ReadAll(file)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// tempFile.Write(fileBytes)
	}
	// uploadFile, handle, err := req.ParseMultipartForm("image")
	// if err != nil {
	// 	webLogPrintln("[Error]", "上传失败" + err.Error())
	// 	w.Write(getErrorMsg(errors.New("上传失败")))
	// 	return
	// }
	// var time int64 = (time.Now().UnixNano() / 1000000)
	// filename := CreateCaptcha()
	// ext := strings.ToLower(path.Ext(handle.Filename))
	// if ext != ".jpg" && ext != ".png" {
	// 	webLogPrintln("[Error]", "图片格式错误" + err.Error())
	// 	w.Write(getErrorMsg(errors.New("图片格式错误")))
	// 	return
	// }
	// filename = filename + strconv.FormatInt(time, 10) + ext

	// saveImageMap(handle.Filename, filename)

	// 保存图片
	// os.Mkdir("./uploaded/", 0777)
	// saveFile, err := os.OpenFile("./uploaded/" + filename, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	webLogPrintln("[Error]", "更新失败" + err.Error())
	// 	w.Write(getErrorMsg(errors.New("更新失败")))
	// 	return
	// }

	// io.Copy(saveFile, uploadFile)
	// defer uploadFile.Close()
	// defer saveFile.Close()
	resData := Uploadp{
		Success: true,
		Data: data,
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

// 显示图片
func showPicHandle( w http.ResponseWriter, req *http.Request ) {
	file, err := os.Open("." + req.URL.Path)
	if err != nil {
		webLogPrintln("[Error]", "打开图片失败" + err.Error())
		w.Write(getErrorMsg(errors.New("打开图片失败")))
		return
	}

	defer file.Close()
	buff, err := ioutil.ReadAll(file)
	if err != nil {
		webLogPrintln("[Error]", "读取图片文件失败" + err.Error())
		w.Write(getErrorMsg(errors.New("读取图片文件失败")))
		return
	}
	w.Write(buff)
}

// 验证图片是否正在使用
func testRemoveImage(url string) bool {
	var res bool
	f, err := os.OpenFile("./index.txt", os.O_RDONLY, 0666)
	if err != nil {
		webLogPrintln("[Error]", "打开文件失败" + err.Error())
		return false
	}
	defer f.Close()
	temp := make([]byte, 1024 * 4)
	fileLen, _ := f.Read(temp)
	data := string(temp[:fileLen])
	kv := strings.Split(data, "\n")
	for i := 0; i < len(kv); i++ {
		items := strings.Split(kv[i], Splitter)
		if len(items) >= 4 {
			match, _ := regexp.MatchString("<img(.*)src=\"" + url, items[3])
			if match {
				res = true
			}
		}
	}
	return res
}

// 删除图片
func removeImg(url string) bool {
	if testRemoveImage(url) {
		return false
	}
	del := os.Remove("." + url)
	if del != nil {
		webLogPrintln("[Error]", "删除文件失败" + del.Error())
	}
	f, err := os.OpenFile("./imageMap.txt", os.O_RDONLY, 0666)
	if err != nil {
		// debugLog.SetPrefix("[Error]")
		// debugLog.Println("删除留言失败" + err.Error())
		webLogPrintln("[Error]", "打开文件失败" + err.Error())
		return false
	}
	defer f.Close()
	temp := make([]byte, 1024 * 4)
	fileLen, _ := f.Read(temp)
	data := string(temp[:fileLen])
	kv := strings.Split(data, "\n")
	var index int
	for i := 0; i < len(kv); i++ {
		items := strings.Split(kv[i], Splitter)
		if len(items) >= 2 && url == items[1] {
			index = i
			break
		}
	}
	if index > -1 {
		kv = append(kv[:index], kv[index+1:]...)
	}
	newData := strings.Join(kv, "\n")
	writeToFile("imageMap.txt", []byte(newData))
	return true
}

// 删除文本内匹配到的图片
func removeImgList(s string) {
	valid := regexp.MustCompile("<img(.*)src=\"(.*)\"")
	arr := valid.FindAllStringSubmatch(s, -1)
	for i := 0; i < len(arr); i += 1 {
		src := arr[i][2]
		fmt.Println(len(src))
		if len(src) > 0 {
			removeImg(src)
		}
	}
}

// 删除图片接口
func removeImage(w http.ResponseWriter, req *http.Request) {
	var user map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &user)
	if !removeImg(user["url"].(string)) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(getErrorMsg(errors.New("删除失败")))
		return
	}
	profile := Resp{
		Success: true,
	}
	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
