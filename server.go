package main

import "C"
import (
	"fmt"
	"github.com/IRA-Package-Manager/goremote/util"
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
	// The same as ipkg run_server
	directory := os.Args[1]
	if _, err := os.Stat(directory); err != nil {
		fmt.Fprintln(os.Stderr, "directory not exists")
	}
	remoteMux := util.NewRemoteMux(directory)
	http.ListenAndServe(":8000", remoteMux)
}
