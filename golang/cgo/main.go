package main

import (
	//void SayHello(const char* s);
	"C"
)

// 使用自己的C函数

/*
#include <stdio.h>

	static void SayHello(const char* s) {
		puts(s)
	}
*/
func main() {
	println("hello, cgo")
	C.SayHello(C.CString("Hello, World\n"))
}
