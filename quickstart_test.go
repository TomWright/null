package null_test

import (
	"encoding/json"
	"fmt"
	"github.com/tomwright/null"
)

func ExampleQuickStart() {
	type data struct {
		PresentString null.String `json:"presentString"`
		BlankString   null.String `json:"blankString"`
	}

	inputBytes := []byte(`{"presentString":"asd","blankString":""}`)

	fmt.Printf("before:\n%s\n", string(inputBytes))

	var d data
	if err := json.Unmarshal(inputBytes, &d); err != nil {
		panic(err)
	}

	fmt.Printf("during:\n"+
		"PresentString.Valid:%v, PresentString.String:%s\n"+
		"BlankString.Valid:%v, BlankString.String:%s\n",
		d.PresentString.Valid, d.PresentString.String,
		d.BlankString.Valid, d.BlankString.String)

	out, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	fmt.Printf("after:\n%s\n", string(out))

	// Output:
	// before:
	// {"presentString":"asd","blankString":""}
	// during:
	// PresentString.Valid:true, PresentString.String:asd
	// BlankString.Valid:false, BlankString.String:
	// after:
	// {"presentString":"asd","blankString":null}
}
