package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			address := fmt.Sprintf("scanme.nmap.org:%d", j)
			conn, err := net.Dial("tcp", address)
			if err != nil {

				return
			}

			err = conn.Close()
			if err != nil {
				return
			}

			fmt.Println("%d oppen\n", i)
		}(i)
		wg.Wait()
	}
}
