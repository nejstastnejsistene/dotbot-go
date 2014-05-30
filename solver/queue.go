package solver

// #include "queue.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type queue C.Queue

func newQueue() *queue {
	return (*queue)(C.NewQueue())
}

func (q *queue) free() {
	C.FreeQueue((*C.Queue)(q))
}

func (q *queue) slice() (slice []Mask) {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sliceHeader.Cap = int(q.capacity)
	sliceHeader.Len = int(q.size)
	sliceHeader.Data = uintptr(unsafe.Pointer(q.values))
	return
}
