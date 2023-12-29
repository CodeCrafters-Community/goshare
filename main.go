package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/user"
	"strings"
	"sync"
)

//go:embed html/*
var multi embed.FS

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

var (
	port         string
	view         string
	viewHelpText string
	path         string
	htmlMap      map[string]string
)

func init() {
	htmlMap = make(map[string]string)

	files, err := multi.ReadDir("html")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	var fileName string
	for _, file := range files {
		if file.IsDir() {
			fmt.Println(file.Name() + " is a directory")
		} else {
			//分割字符取第一个
			fileName = strings.Split(file.Name(), ".")[0]
			viewHelpText = viewHelpText + fileName + "\t"
			htmlMap[fileName] = file.Name()
		}
	}
	viewHelpText = viewHelpText + "\ndisplay *.html"
}

func main() {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = os.MkdirAll(currentUser.HomeDir+"/Desktop/upfile", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&port, "p", "9000", "port example: -p 9000;端口 例子: -p 9000")
	flag.StringVar(&view, "v", "simple", viewHelpText)
	flag.StringVar(&path, "P", currentUser.HomeDir+"/Desktop/upfile", "download/upload directory 下载/上传目录")
	fmt.Println("默认上传/下载路径：\t" + currentUser.HomeDir + "/Desktop/upfile")
	flag.Parse()
	fmt.Println("请访问下面的链接:")
	showIp()
	http.HandleFunc("/", uploadFileHandler)
	http.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(path))))
	http.ListenAndServe(":"+port, nil)
}
func nameToView(view string) string {
	content, err := multi.ReadFile("html/" + view)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, nameToView(htmlMap[view]))
	switch r.Method {
	case "POST":
		upload(w, r)
	}
}

func writeFile(wg *sync.WaitGroup, file *os.File, multiFile multipart.File) {
	defer wg.Done()
	io.Copy(file, multiFile)
}

func upload(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseMultipartForm(1 << 30)
	var wg sync.WaitGroup
	files := r.MultipartForm.File["file"]
	for _, file := range files {
		create, err := os.Create(path + "/" + file.Filename)
		check(err)
		wg.Add(1)
		open, err := file.Open()
		check(err)
		go writeFile(&wg, create, open)
		wg.Wait()
		fmt.Println(" upload successfully:" + path + "/" + file.Filename)
		w.Write([]byte(file.Filename + "&nbsp;<span style='color:#00DB00'>SUCCESS</span></br>"))
	}
}

func showIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println("http://" + ipNet.IP.String() + ":" + port)
			}
		}
	}
}
