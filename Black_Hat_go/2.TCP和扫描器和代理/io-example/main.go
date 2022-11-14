package main

import (
	"fmt"
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
	//实例化reader和writer
	var (
		reader FooReader
		writer FooWriter
	)

	//创建缓冲区已保存输入/输出
	input := make([]byte, 4096)

	//使用reader读取输入
	s, err := reader.Read(input)
	if err != nil {
		log.Fatalln("Unable to read data")
	}
	fmt.Printf("Read %d bytes from stdin\n", s)

	//使用writer写入输出
	s, err = writer.Write(input)
	if err != nil {
		log.Fatalln("Unable to write data")
	}
	fmt.Printf("Wrote %d bytes to stdout", s)
}
