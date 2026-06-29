//go:build wasmjs
package main

import (
	"log"
	"syscall/js"

	placetypes "github.com/whosonfirst/go-whosonfirst/v4/placetypes/wasm"
	format "github.com/whosonfirst/go-whosonfirst/v4/format/wasm"
	"github.com/whosonfirst/go-whosonfirst/v4/validate"			
	validate_wasm "github.com/whosonfirst/go-whosonfirst/v4/validate/wasm"
	export_wasm "github.com/whosonfirst/go-whosonfirst/v4/export/wasm"			
)

func main() {

	prep_func := export_wasm.PrepareWithoutTimestampsFunc()
	defer prep_func.Release()

	js.Global().Set("wof_prepare_feature", prep_func)
	
	placetypes_func := placetypes.PlacetypesFunc()
	defer placetypes_func.Release()

	isvalid_func := placetypes.IsValidPlacetypeFunc()
	defer isvalid_func.Release()

	children_func := placetypes.ChildrenFunc()
	defer children_func.Release()

	descendants_func := placetypes.DescendantsFunc()
	defer descendants_func.Release()

	ancestors_func := placetypes.AncestorsFunc()
	defer ancestors_func.Release()
	
	js.Global().Set("wof_placetypes", placetypes_func)
	js.Global().Set("wof_placetypes_is_valid", isvalid_func)
	js.Global().Set("wof_placetypes_children", children_func)
	js.Global().Set("wof_placetypes_descendants", descendants_func)
	js.Global().Set("wof_placetypes_ancestors", ancestors_func)				
	
	format_func := format.FormatFunc()
	defer format_func.Release()

	js.Global().Set("wof_format", format_func)

	opts := validate.DefaultValidateOptions()
	validate_func := validate_wasm.ValidateFunc(opts)

	defer validate_func.Release()

	js.Global().Set("wof_validate", validate_func)
	
	c := make(chan struct{}, 0)

	log.Println("wof-edit functions initialized")
	<-c
}

