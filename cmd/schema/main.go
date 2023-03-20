package main

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func main() {
	schemaLoader := gojsonschema.NewReferenceLoader("https://raw.githubusercontent.com/oxheadalpha/tezos-snapshot-metadata-schema/main/tezos-snapshot-metadata.schema.json")
	documentLoader := gojsonschema.NewReferenceLoader("http://localhost:8080/tezos-snapshots.json")

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}
