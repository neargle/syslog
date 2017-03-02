package main

import (
	"fmt"
	//	"os"
	//	"strings"

	"./info"
)

func main() {
	dir := info.GetSysLog("windows", "all")
	fmt.Println(dir)
	fmt.Println( /*os.Hostname()*/ )
}
