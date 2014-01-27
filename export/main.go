package main

import (
	"fmt"
	"runtime"
)

func main() {
	// mongoexport -c new_builds -h zach.mongohq.com --port 10081 -d travisarchive -u travisarchive -p ILoveOwen1028 --out foo.json
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)
}
