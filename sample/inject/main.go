package main

import pdfinject "../../pdfinject"

func main() {

	// Create the form values.
	form := map[string]interface{}{
		"location": "Hello",
		"function": "World",
		"reason": "Kim",
		"headcountAddition": "Yes",
		"directors.0": "Yo",
		"kazekim_checkbox": "Yes",
	}

	pdfInject := pdfinject.New()
	_, err := pdfInject.FillWithDestFile(form, "./sample/inject/sample.pdf", "./sample/inject/filled.pdf")
	if err != nil {
		panic(err)
	}

}
