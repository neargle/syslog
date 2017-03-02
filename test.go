package main

import (
	"fmt"
	//	"os"
	//	"strings"

	"./info"
)

func main() {
	ba := []string{"10.100.159.125", "127.0.0.1"}
	info.BlackList = ba
	dir := info.GetSysLog("windows", "all")
	fmt.Println(dir)
	fmt.Println( /*os.Hostname()*/ )
}
