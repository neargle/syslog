package main

import (
	"fmt"

	"./info"
)

func main() {
	dir := info.GetSysLog("windows", "all")
	fmt.Println(dir)
}
