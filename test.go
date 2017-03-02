package main

import (
	"fmt"

	"./info"
)

func main() {
	dir := info.GetSysLog("linux", "all")
	fmt.Println(dir)
}
