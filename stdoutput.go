package asmiostat

import (
	"fmt"
)

func startStdoutput(src chan string) {
	go func() {
		for {
			//fmt.Println("-waiting for output")
			fmt.Printf("%s", <-src)
		}
	}()
}
