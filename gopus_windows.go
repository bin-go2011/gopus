package gopus

import (
	"fmt"
	"os"
	"path"
	"unsafe"

	"golang.org/x/sys/windows"
)

var loadedDLL *windows.DLL

const opusPath = "opus.dll"

var (
	opusGetVersionStringProc *windows.Proc
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := path.Join(dir, "lib", opusPath)
	fmt.Println(path)
	loadedDLL = windows.MustLoadDLL(path)
	opusGetVersionStringProc = loadedDLL.MustFindProc("opus_get_version_string")
}

func byteSliceToString(bval []byte) string {
	for i := range bval {
		if bval[i] == 0 {
			return string(bval[:i])
		}
	}
	return string(bval[:])
}

// bytePtrToString returns a string copied from pointer to a null terminated byte array
// WARNING: ONLY SAFE WITH IF r POINTS TO C MEMORY!
// govet will complain about this function for the reason stated above
func bytePtrToString(r uintptr) string {
	if r == 0 {
		return ""
	}
	bval := (*[1 << 30]byte)(unsafe.Pointer(r))
	return byteSliceToString(bval[:])
}

func opusGetVersionString() string {
	r1, _, _ := opusGetVersionStringProc.Call()
	return bytePtrToString(r1)
}
