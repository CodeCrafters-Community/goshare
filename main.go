package main

import (
	"embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"strings"
)

//go:embed html/*
var multi  embed.FS

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
var(
	port string
	view string
	viewHelpText string
	path string
	htmlMap map[string]string
)

func init(){
	htmlMap =  make(map[string]string)

	files, err := multi.ReadDir("html")
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    var fileName string
    for _, file := range files {
        if file.IsDir() {
            fmt.Println(file.Name()+" is a directory")
        } else {
			//分割字符取第一个
			fileName = strings.Split(file.Name(),".")[0]
			viewHelpText = viewHelpText + fileName + "\t"
            htmlMap[fileName] = file.Name()
        }
    }
	viewHelpText = viewHelpText + "\ndisplay *.html"
	fmt.Println(os.Getwd())
}


func main() {
	currentUser, err := user.Current()
    if err != nil {
       log.Fatalf(err.Error())
    }
	
	err = os.MkdirAll(currentUser.HomeDir+"/Desktop/upfile",os.ModePerm)
	if err!=nil{
		log.Fatal(err)
	}
	flag.StringVar(&port, "p", "9000", "port example: -p 9000;端口 例子: -p 9000")
	flag.StringVar(&view, "v", "simple",viewHelpText)
	flag.StringVar(&path, "P", currentUser.HomeDir+"/Desktop/upfile","download/upload directory 下载/上传目录")

	flag.Parse()
	fmt.Println("请访问下面的链接:")
	showip()
	http.HandleFunc("/", uploadFileHandler)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(path))))
	http.ListenAndServe(":"+port, nil)
}
func nameToView(view string)string{
	content ,err := multi.ReadFile("html/"+view)
	if err != nil {
        panic(err)
    }
    return string(content)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, nameToView(htmlMap[view]))
	uploadOne(w, r)
}

func uploadOne(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		file, handler, err := r.FormFile("file_uploads") //name的字段
		check(err)
		defer file.Close()
		fileBytes, err := ioutil.ReadAll(file)
		check(err)
		newFile, err := os.Create(path + "/" + handler.Filename)
		check(err)
		defer newFile.Close()
		if _, err := newFile.Write(fileBytes); err != nil {
			check(err)
			return
		}
		fmt.Println(" upload successfully:" + path + "/" + handler.Filename)
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
			cur, err := os.Create(path + "/" + fileItem.Filename)
			fmt.Println("上传地址:path/")
			check(err)
			defer cur.Close()
			if _, err := cur.Write(fileBytes); err != nil {
				check(err)
				return
			}
			fmt.Println(" upload successfully:" + path + "/" + fileItem.Filename)
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
				fmt.Println(ipnet.IP.String() + ":" + port)
			}
		}
	}
}
