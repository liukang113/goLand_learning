package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print("hello word")
	if len(os.Args) > 1 {
		fmt.Println("满足参数要求，请求打印", os.Args[1])
	}
}
