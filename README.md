# pdfinject-go
Inject Value into Adobe PDF Form with GoLang. You can create FDF Data with values to fill in the PDF Form you create from Adobe Acrobat Pro DC. 

This library is support for text field, radio button, check box and also support for pdf overlay on pdf function 

Inspiration from : https://github.com/desertbit/fillpdf


## Usage:

Define you form with map[string]interface{}

input your source file name replace in "sample.pdf" and define your output file name replace in "filled.pdf"

```go

    form := map[string]interface{}{
		"location": "Hello",
		"function": "World",
		"reason": "Kim",
		"headcountAddition": "Yes",
		"directors.0": "Yo",
		"kazekim_checkbox": "Yes",
	}

	pdfInject := pdfinject.New()
	err := pdfInject.FillWithDestFile(form, "sample.pdf", "filled.pdf")
	if err != nil {
		panic(err)
	}

```
You can see example in sample directory

## Merge PDF Files
input working directory, output file and input files.
```go
	err := pdfinject.MergePDF(
		"./sample/merge/pdf/",
		"./out.pdf",
		"./a1.pdf",
		"./a2.pdf",
		"./a3.pdf",
		"./a4.pdf",
	)
```
