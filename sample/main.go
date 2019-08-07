package main

import (
	"github.com/kazekim/pdfinject-go"
)

func main() {

	// Create the form values.
	form := pdfinject.Form{
		"location": "Hello",
		"function": "World",
		"reason": "Kim",
		"headcountAddition": "Yes",
		"directors.0": "Yo",
	}

	pdfInject := pdfinject.New()
	err := pdfInject.FillWithDestFile(form, "sample.pdf", "filled.pdf")
	if err != nil {
		panic(err)
	}

}
