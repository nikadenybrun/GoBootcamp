package window

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#include "application.h"
#include "window.h"
*/
import "C"
import "unsafe"

func InitApplication() {
	C.InitApplication()
}

func RunApplication() {
	C.RunApplication()
}

func WindowCreate(x, y, w, h int, title string) unsafe.Pointer {
	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))

	return C.Window_Create(C.int(x), C.int(y), C.int(w), C.int(h), t)
}

func MakeKeyAndOrderFront(ptr unsafe.Pointer) {
	C.Window_MakeKeyAndOrderFront(ptr)
}
