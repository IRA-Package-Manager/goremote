package main

import "C"
import (
	"fmt"
	"ira/remote/util"
	"net/http"
	"os"
)

//export RemoteServe
func RemoteServe(host *C.char, port int, unsafeDirectory *C.char) (int, *C.char) {
	directory := C.GoString(unsafeDirectory)
	if _, err := os.Stat(directory); err != nil {
		return 401, C.CString("Directory not exists")
	}
	remoteMux := util.NewRemoteMux(directory)
	http.ListenAndServe(fmt.Sprintf("%s:%d", C.GoString(host), port), remoteMux)
	return 0, C.CString("")
}

func main() {
	// Used to compile C library
}
