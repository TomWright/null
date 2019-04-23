package null_test

import (
	"testing"

	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/tomwright/null"
	"reflect"
	"time"
)

func ExampleNewTime() {
	a := null.NewTime(time.Time{})
	b := null.NewTime(time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC))

	fmt.Printf("a: valid: `%v`, time: `%s`\n", a.Valid(), a.Time.Format(time.RFC3339))
	fmt.Printf("b: valid: `%v`, time: `%s`\n", b.Valid(), b.Time.Format(time.RFC3339))

	// Output:
	// a: valid: `false`, time: `0001-01-01T00:00:00Z`
	// b: valid: `true`, time: `2019-01-01T12:00:00Z`
}

func ExampleTime_Scan() {
	var a, b, c null.Time

	_ = a.Scan(nil)
	_ = b.Scan(time.Time{})
	_ = c.Scan(time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC))

	fmt.Printf("a: valid: `%v`, time: `%s`\n", a.Valid(), a.Time.Format(time.RFC3339))
	fmt.Printf("b: valid: `%v`, time: `%s`\n", b.Valid(), b.Time.Format(time.RFC3339))
	fmt.Printf("c: valid: `%v`, time: `%s`\n", c.Valid(), c.Time.Format(time.RFC3339))

	// Output:
	// a: valid: `false`, time: `0001-01-01T00:00:00Z`
	// b: valid: `false`, time: `0001-01-01T00:00:00Z`
	// c: valid: `true`, time: `2019-01-01T12:00:00Z`
}

// ExampleTime_MarshalJSON shows how different values will be marshal'd to JSON.
// Notice that `a` starts off as an invalid value and is output as `null`, but
// after adding a second to the time it then becomes valid and is output as
// a time string.
func ExampleTime_MarshalJSON() {
	a := null.NewTime(time.Time{})
	b := null.NewTime(time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC))

	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)

	fmt.Printf("a: %s\n", string(aBytes))
	fmt.Printf("b: %s\n", string(bBytes))

	a.Time = a.Add(time.Second)
	aBytes, _ = json.Marshal(a)

	fmt.Printf("a after add: %s", string(aBytes))

	// Output:
	// a: null
	// b: "2019-01-01T12:00:00Z"
	// a after add: "0001-01-01T00:00:01Z"
}

func ExampleTime_UnmarshalJSON() {
	var a, b, c null.Time

	_ = json.Unmarshal([]byte(`"2019-01-01T12:00:00Z"`), &a)
	_ = json.Unmarshal([]byte(`null`), &b)
	_ = json.Unmarshal([]byte(`""`), &c)

	fmt.Printf("a: valid: `%v`, time: `%s`\n", a.Valid(), a.Time.Format(time.RFC3339))
	fmt.Printf("b: valid: `%v`, time: `%s`\n", b.Valid(), b.Time.Format(time.RFC3339))
	fmt.Printf("c: valid: `%v`, time: `%s`\n", c.Valid(), c.Time.Format(time.RFC3339))

	// Output:
	// a: valid: `true`, time: `2019-01-01T12:00:00Z`
	// b: valid: `false`, time: `0001-01-01T00:00:00Z`
	// c: valid: `false`, time: `0001-01-01T00:00:00Z`
}

func TestNewTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc  string
		value time.Time
		exp   null.Time
	}{
		{
			"blank value",
			time.Time{},
			null.NewTime(time.Time{}),
		},
		{
			"filled value",
			time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC),
			null.NewTime(time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC)),
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			out := null.NewTime(tc.value)

			if exp, got := tc.exp, out; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected `%v`, got `%v`", exp, got)
			}
		})
	}
}

func TestTime_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		exp  driver.Value
		val  null.Time
	}{
		{
			"blank value",
			nil,
			null.NewTime(time.Time{}),
		},
		{
			"valid value after processing",
			time.Date(0001, 01, 01, 00, 00, 01, 0, time.UTC),
			null.NewTime(time.Time{}.Add(time.Second)),
		},
		{
			"nil value after processing",
			nil,
			null.NewTime(time.Date(0001, 01, 01, 00, 00, 01, 0, time.UTC).Add(-time.Second)),
		},
		{
			"filled value",
			time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC),
			null.NewTime(time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC)),
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			val, err := tc.val.Value()
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

func TestTime_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		exp  bool
		val  func() null.Time
	}{
		{
			"blank value",
			false,
			func() null.Time {
				return null.NewTime(time.Time{})
			},
		},
		{
			"valid value",
			true,
			func() null.Time {
				return null.NewTime(time.Now())
			},
		},
		{
			"valid value after adding to invalid value",
			true,
			func() null.Time {
				x := null.NewTime(time.Time{})
				x.Time = x.Add(time.Second)
				return x
			},
		},
		{
			"invalid value after removing from valid value",
			false,
			func() null.Time {
				x := null.NewTime(time.Date(01, 01, 01, 0, 0, 1, 0, time.UTC))
				x.Time = x.Add(-time.Second)
				return x
			},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			val := tc.val()

			if exp, got := tc.exp, val.Valid(); exp != got {
				t.Errorf("expected `%v`, got `%v`", exp, got)
			}
		})
	}
}

func TestTime_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc  string
		value interface{}
		exp   null.Time
	}{
		{
			"blank value",
			time.Time{},
			null.Time{Time: time.Time{}},
		},
		{
			"filled value",
			time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC),
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC)},
		},
		{
			"filled value",
			nil,
			null.Time{Time: time.Time{}},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			out := &null.Time{}
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

func TestTime_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc    string
		expJson []byte
		str     null.Time
	}{
		{
			"invalid blank value",
			[]byte(`null`),
			null.NewTime(time.Time{}),
		},
		{
			"filled value",
			[]byte(`"2019-01-01T12:00:00Z"`),
			null.NewTime(time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC)),
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

func TestTime_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc   string
		jsonIn []byte
		exp    null.Time
	}{
		{
			"null value",
			[]byte(`null`),
			null.Time{Time: time.Time{}},
		},
		{
			"filled string",
			[]byte(`"2019-01-01T12:00:00Z"`),
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC)},
		},
		{
			"blank string",
			[]byte(`"0001-01-01T00:00:00Z"`),
			null.Time{Time: time.Time{}},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			var out null.Time
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

		var out null.Time
		if err := (&out).UnmarshalJSON([]byte(``)); err != nil {
			if err.Error() != "unexpected end of JSON input" {
				t.Errorf("unexpected error: %s", err)
			}
		} else {
			t.Errorf("expected error but got none")
		}
	})
}
