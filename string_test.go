package null_test

import (
	"testing"

	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/tomwright/null"
	"reflect"
)

func ExampleNewString() {
	a := null.NewString("asd")
	b := null.NewString("")

	fmt.Printf("a: valid: `%v`, string: `%s`\n", a.Valid, a.String)
	fmt.Printf("b: valid: `%v`, string: `%s`\n", b.Valid, b.String)

	// Output:
	// a: valid: `true`, string: `asd`
	// b: valid: `false`, string: ``
}

func ExampleString_MarshalJSON() {
	a := null.NewString("asd")
	b := null.NewString("")

	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)

	fmt.Printf("a: %s\n", string(aBytes))
	fmt.Printf("b: %s", string(bBytes))

	// Output:
	// a: "asd"
	// b: null
}

func ExampleString_UnmarshalJSON() {
	var (
		a null.String
		b null.String
		c null.String
		d null.String
	)

	json.Unmarshal([]byte(`"asd"`), &a)
	json.Unmarshal([]byte(`null`), &b)
	json.Unmarshal([]byte(`""`), &c)
	json.Unmarshal(nil, &d)

	fmt.Printf("a: valid: `%v`, string: `%s`\n", a.Valid, a.String)
	fmt.Printf("b: valid: `%v`, string: `%s`\n", b.Valid, b.String)
	fmt.Printf("c: valid: `%v`, string: `%s`\n", c.Valid, c.String)
	fmt.Printf("d: valid: `%v`, string: `%s`\n", d.Valid, d.String)

	// Output:
	// a: valid: `true`, string: `asd`
	// b: valid: `false`, string: ``
	// c: valid: `false`, string: ``
	// d: valid: `false`, string: ``
}

func TestNewString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc  string
		value string
		exp   null.String
	}{
		{
			"blank value",
			"",
			null.String{String: "", Valid: false},
		},
		{
			"filled value",
			"a",
			null.String{String: "a", Valid: true},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			out := null.NewString(tc.value)

			if exp, got := tc.exp, out; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected `%v`, got `%v`", exp, got)
			}
		})
	}
}

func TestString_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		exp  driver.Value
		str  null.String
	}{
		{
			"nil value",
			nil,
			null.String{String: "", Valid: false},
		},
		{
			"blank value",
			nil,
			null.String{String: "", Valid: false},
		},
		{
			"filled value",
			"a",
			null.String{String: "a", Valid: true},
		},
		{
			"blank valid value",
			"",
			null.String{String: "", Valid: true},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			val, err := tc.str.Value()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			if exp, got := tc.exp, val; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected `%v`, got `%v`", exp, got)
			}
		})
	}
}

func TestString_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc  string
		value interface{}
		exp   null.String
	}{
		{
			"nil value",
			nil,
			null.String{String: "", Valid: false},
		},
		{
			"blank value",
			"",
			null.String{String: "", Valid: false},
		},
		{
			"filled value",
			"a",
			null.String{String: "a", Valid: true},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			out := &null.String{}
			if err := out.Scan(tc.value); err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			if exp, got := tc.exp, *out; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected `%v`, got `%v`", exp, got)
			}
		})
	}
}

func TestString_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc    string
		expJson []byte
		str     null.String
	}{
		{
			"invalid blank string",
			[]byte(`null`),
			null.String{String: "", Valid: false},
		},
		{
			"invalid filled string",
			[]byte(`null`),
			null.String{String: "a", Valid: false},
		},
		{
			"filled value",
			[]byte(`"a"`),
			null.String{String: "a", Valid: true},
		},
		{
			"blank valid value",
			[]byte(`""`),
			null.String{String: "", Valid: true},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			out, err := tc.str.MarshalJSON()
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			if exp, got := tc.expJson, out; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected `%v`, got `%v`", string(exp), string(got))
			}
		})
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc   string
		jsonIn []byte
		exp    null.String
	}{
		{
			"null value",
			[]byte(`null`),
			null.String{String: "", Valid: false},
		},
		{
			"filled string",
			[]byte(`"a"`),
			null.String{String: "a", Valid: true},
		},
		{
			"blank string",
			[]byte(`""`),
			null.String{String: "", Valid: false},
		},
		{
			"nil bytes",
			nil,
			null.String{String: "", Valid: false},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			var out null.String
			if err := (&out).UnmarshalJSON(tc.jsonIn); err != nil {
				t.Errorf("unexpected error: %s", err)
				return
			}

			if exp, got := tc.exp, out; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected `%v`, got `%v`", exp, got)
			}
		})
	}

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()

		var out null.String
		if err := (&out).UnmarshalJSON([]byte(``)); err != nil {
			if err.Error() != "unexpected end of JSON input" {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			t.Errorf("expected error but got none")
		}
	})
}
