//go:build wasmjs
package wasm

import (
	"fmt"
	"syscall/js"

	"github.com/whosonfirst/go-whosonfirst/v4/validate"
)

func ValidateFunc(opts *validate.Options) js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		geojson_data := args[0].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			go func() {

				err := validate.ValidateWithOptions([]byte(geojson_data), opts)

				if err != nil {
					reject.Invoke(fmt.Sprintf("Failed to export data, %v", err))
					return
				}

				resolve.Invoke()
			}()

			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}

