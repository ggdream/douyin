package main

import (
	"fmt"
	"os"
)

func argv() (string, string) {
	if len(os.Args) == 1 {
		fmt.Println("参数不足")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		return os.Args[1], "./"
	}

	return os.Args[1], os.Args[2]
}

func main() {
	file, path := argv()
	if err := DouYin(file, path); err != nil {
		panic(err)
	}
}
