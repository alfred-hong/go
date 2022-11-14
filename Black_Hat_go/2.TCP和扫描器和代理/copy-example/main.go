package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// 定义了一个从标准输入stdin读取的io.Reader
type FooReader struct{}

// 从标准输入stdin读取数据
func (fooReader *FooReader) Read(b []byte) (int, error) {
	fmt.Print("in >")
	return os.Stdin.Read(b)
}

// 定一个写入标准输出stdout的io.Writer
type FooWriter struct{}

//将数据写入标准输出stdout
func (fooWriter *FooWriter) Write(b []byte) (int, error) {
	fmt.Print("out >")
	return os.Stdout.Write(b)
}

func main() {
	var (
		reader FooReader
		writer FooWriter
	)
	if _, err := io.Copy(&writer, &reader); err != nil {
		log.Fatalln("Unable to read/write date")
	}
}
