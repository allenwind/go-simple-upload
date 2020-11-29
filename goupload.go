package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

/*
application/x-www-form-urlencoded   表示在发送前编码所有字符（默认）
multipart/form-data	  不对字符编码。在使用包含文件上传控件的表单时，必须使用该值。
text/plain	  空格转换为 "+" 加号，但不对特殊字符编码。
*/

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL)
	if r.Method == "GET" {
		current := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(current, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.tpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 10)
		file, handler, err := r.FormFile("uploadfile") // get file handle
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v\r\n上传成功", handler.Header) // response
		f, err := os.OpenFile("./upload/"+handler.Filename,
			os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	http.HandleFunc("/", upload)
	log.Fatal(http.ListenAndServe(":8899", nil))
}
