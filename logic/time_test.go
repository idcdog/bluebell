package logic

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

func TestTime(t *testing.T) {
	var a = time.Now()
	fmt.Println(unsafe.Sizeof(a))
	var b int8 = 1
	fmt.Println(unsafe.Sizeof(b))
	var c string = ""
	fmt.Println(unsafe.Sizeof(c))
}
