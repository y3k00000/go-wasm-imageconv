package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"syscall/js"
	"time"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func add(a, b int) int {
	return a + b
}

func addJS(this js.Value, args []js.Value) interface{} {
	return add(args[1].Int(), args[1].Int())
}

func toPng(imageBytes []byte) ([]byte, error) {
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
	case "image/png":
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode jpeg")
		}

		img = resize.Resize(480, 0, img, resize.Lanczos3)

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, errors.Wrap(err, "unable to encode png")
		}

		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unable to convert %#v to png", contentType)
}

func toPngJS(this js.Value, args []js.Value) interface{} {
	content := make([]byte, args[0].Length())
	js.CopyBytesToGo(content, args[0])
	result, err := toPng(content)
	if err != nil {
		panic(err)
	}
	output := js.Global().Get("Uint8Array").New(len(result))
	js.CopyBytesToJS(output, result)
	return output
}

func main() {
	fmt.Println("it works!")
	js.Global().Set("add", js.FuncOf(addJS))
	js.Global().Set("toPng", js.FuncOf(toPngJS))
	go func() {
		for {
			fmt.Println("Yoyo")
			time.Sleep(5 * time.Second)
		}
	}()
	waitC := make(chan (int), 1)
	<-waitC
}
