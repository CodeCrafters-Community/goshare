package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestLoadHtml(t *testing.T){
	content ,err := ioutil.ReadFile("../html/simple.html")
	if err != nil {
        panic(err)
    }
    fmt.Println(string(content))
}

func TestDirMap(t *testing.T){
	htmlMap :=  make(map[string]string)

	files, err := ioutil.ReadDir("../html")
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
            htmlMap[fileName] = file.Name()
        }
    }
    fmt.Println(htmlMap)
}