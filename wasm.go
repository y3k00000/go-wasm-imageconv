package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"runtime"
	"syscall/js" // syscall/js shows warning on vscode but why?
	"time"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func add(a, b int) int {
	return a + b
}

func addJS(this js.Value, args []js.Value) interface{} { // example 1 : add() to js global scope
	return add(args[1].Int(), args[1].Int())
}

func toPng(imageBytes []byte) ([]byte, error) { // convert jpg to png and output
	contentType := http.DetectContentType(imageBytes)

	switch contentType {
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

func toPngJS(this js.Value, args []js.Value) interface{} { // example 2 : toPng() to js global scope
	content := make([]byte, args[0].Length())
	js.CopyBytesToGo(content, args[0]) // get input image
	result, err := toPng(content)
	if err != nil {
		panic(err)
	}
	output := js.Global().Get("Uint8Array").New(len(result))
	js.CopyBytesToJS(output, result) // send output image
	return output
}

func statMemory() string { // https://golangcode.com/print-the-current-memory-usage/
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	stat := make(map[string]interface{})
	stat["m.Alloc"] = m.Alloc
	stat["m.TotalAlloc"] = m.TotalAlloc
	stat["m.Sys"] = m.Sys
	stat["m.NumGC"] = m.NumGC
	statBytes, err := json.Marshal(stat)
	if err != nil {
		panic(err)
	}
	return string(statBytes)
}

func statMemoryJS(this js.Value, args []js.Value) interface{} {
	return statMemory()
}

func main() {
	fmt.Println("it works!")
	js.Global().Set("add", js.FuncOf(addJS))
	js.Global().Set("toPng", js.FuncOf(toPngJS))
	js.Global().Set("statMemory", js.FuncOf(statMemoryJS))
	go func() { // goroutine works inside wasm too
		for {
			fmt.Println("From goroutine!!")
			time.Sleep(5 * time.Second)
		}
	}()
	waitC := make(chan (int), 1)
	<-waitC
}
