package main

import (
	"../../pdfinject"
	"fmt"
)

func main() {
	err := pdfinject.MergePDF(
		"./sample/merge/pdf/",
		"./out1234.pdf",
		"./a1.pdf",
		"./a2.pdf",
		"./a3.pdf",
		"./a4.pdf",
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Fin")
}

