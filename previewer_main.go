package main

import (
	"os"
	"zshanjun/markdown/previewer"
)

//用法：
//1、编译  go build -o previewer2 previewer_main.go
//2、运行  ./previewer2 test.md
func main() {
	previewer2 := previewer.NewPreviewer(8089)
	previewer2.UseBasic()
	previewer2.Run(os.Args...)
}
