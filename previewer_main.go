package main

import (
	"os"
	"zshanjun/markdown/previewer"
)

func main() {
	previewer2 := previewer.NewPreviewer(8089)
	previewer2.UseBasic()
	previewer2.Run(os.Args...)
}
