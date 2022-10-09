package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("OK")
	path, _ := os.Getwd()
	fmt.Println("local path ", path)
}
