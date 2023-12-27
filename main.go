package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func main() {
	fmt.Println("请访问下面的链接:")
	showip()
	http.HandleFunc("/", uploadFileHandler)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("/home/m2/upfile/filesDir"))))
	http.ListenAndServe(":9000", nil)
}
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	/**/
	fmt.Fprintln(w, `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="viewport" charset="UTF-8" content="width=device-width, initial-scale=1.0">
    <title>m的二次方文件上传</title>
</head>
<body style="text-align: center;"> 
    <h1>m的二次方文件上传</h1>
    <br>
    <br>
    <form action="UploadFile.ashx" method="post" enctype="multipart/form-data">
    <input type="file" name="file_uploads" id="file_uploads" multiple/>
    <input type="submit" name="上传文件"/>
    </form>
        <br>
    <br>
        <br>
    <br>
    <a href="/file">文件下载</a>
</body>
</html>
        `)
	uploadOne(w, r)
}

func uploadOne(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, handler, err := r.FormFile("file_uploads") //name的字段
		check(err)
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		check(err)
		newFile, err := os.Create("/home/m2/upfile/filesDir/" + handler.Filename)
		check(err)
		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil {
			check(err)
			return
		}
		fmt.Println(" upload successfully:" + "/home/m2/upfile/filesDir/" + handler.Filename)
		w.Write([]byte("SUCCESS"))
	}
}
func uploadMore(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("======================进来了")
		//设置内存大小
		// r.ParseMultipartForm(32 << 20)
		//获取上传的文件组
		files := r.MultipartForm.File["file"]
		log.Println(files, "==================")
		for _, fileItem := range files {
			//打开上传文件
			file, err := fileItem.Open()
			fileBytes, err := ioutil.ReadAll(file)

			defer file.Close()
			check(err)

			//创建上传文件
			cur, err := os.Create("/home/m2/upfile/filesDir/" + fileItem.Filename)
			fmt.Println("上传地址:/home/m2/upfile/filesDir/")
			check(err)
			defer cur.Close()
			if _, err := cur.Write(fileBytes); err != nil {
				check(err)
				return
			}
			fmt.Println(" upload successfully:" + "/home/m2/upfile/filesDir/" + fileItem.Filename)
			w.Write([]byte("SUCCESS"))
		}
	}
}
func showip() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(ipnet.IP.String() + ":9000")
			}
		}
	}
}
