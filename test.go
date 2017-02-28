package main

import (
	"fmt"

	"./info"
)

func main() {
	dir := info.GetSysLog("windows", 10, "2017/02/10-01:12:22")
	fmt.Println(dir)
}
