/*
 * @Author: alfred_hong 1911948020@qq.com
 * @Date: 2022-05-14 18:36:33
 * @LastEditors: alfred_hong 1911948020@qq.com
 * @LastEditTime: 2022-05-14 18:38:00
 * @FilePath: /go代码/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Print("hello '/'") })
	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) { fmt.Print("hello '/1'") })
	
	http.HandleFunc("/2", func(w http.ResponseWriter, r *http.Request) { fmt.Print("hello '/2'") })
}
