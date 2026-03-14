package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/m-theory-io/idgen"
)

//export DocID
func DocID(prefix *C.char, format *C.char) *C.char {
	goPrefix := ""
	goFormat := ""

	if prefix != nil {
		goPrefix = C.GoString(prefix)
	}
	if format != nil {
		goFormat = C.GoString(format)
	}

	id := idgen.DocID(goPrefix, goFormat)
	return C.CString(id)
}

//export FreeCString
func FreeCString(s *C.char) {
	if s == nil {
		return
	}
	C.free(unsafe.Pointer(s))
}

func main() {}
