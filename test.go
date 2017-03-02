package main

import (
	"fmt"

	"./info"
)

func main() {
	dir := info.GetSysLog("linux", "02/12-00:00")
	fmt.Println(dir)
}