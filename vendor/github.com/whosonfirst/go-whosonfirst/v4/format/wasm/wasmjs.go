//go:build wasmjs
package wasm

import (
	"fmt"
	"syscall/js"

	"github.com/whosonfirst/go-whosonfirst/v4/format"
)

func FormatFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		feature_str := args[0].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			feature_fmt, err := format.FormatBytes([]byte(feature_str))

			if err != nil {
				reject.Invoke(fmt.Printf("Failed to format feature, %v\n", err))
				return nil
			}

			resolve.Invoke(string(feature_fmt))
			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
