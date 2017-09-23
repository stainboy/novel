package main

import (
	"fmt"
	"os"

	"./pkg/util"
	"./pkg/zongheng"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: novel [id]")
		os.Exit(1)
	}
	util.Source(fmt.Sprintf("%s.metadata", os.Args[1]))
	if err := zongheng.NewClawer().Process(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
