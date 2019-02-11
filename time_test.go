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

	fmt.Printf("a: valid: `%v`, time: `%s`\n", a.Valid, a.Time.Format(time.RFC3339))
	fmt.Printf("b: valid: `%v`, time: `%s`\n", b.Valid, b.Time.Format(time.RFC3339))

	// Output:
	// a: valid: `false`, time: `0001-01-01T00:00:00Z`
	// b: valid: `true`, time: `2019-01-01T12:00:00Z`
}

func ExampleTimes_MarshalJSON() {
	a := null.NewTime(time.Time{})
	b := null.NewTime(time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC))

	aBytes, _ := json.Marshal(a)
	bBytes, _ := json.Marshal(b)

	fmt.Printf("a: %s\n", string(aBytes))
	fmt.Printf("b: %s", string(bBytes))

	// Output:
	// a: null
	// b: "2019-01-01T12:00:00Z"
}

func ExampleTime_UnmarshalJSON() {
	var (
		a null.Time
		b null.Time
		c null.Time
	)

	json.Unmarshal([]byte(`"2019-01-01T12:00:00Z"`), &a)
	json.Unmarshal([]byte(`null`), &b)
	json.Unmarshal([]byte(`""`), &c)

	fmt.Printf("a: valid: `%v`, time: `%s`\n", a.Valid, a.Time.Format(time.RFC3339))
	fmt.Printf("b: valid: `%v`, time: `%s`\n", b.Valid, b.Time.Format(time.RFC3339))
	fmt.Printf("c: valid: `%v`, time: `%s`\n", c.Valid, c.Time.Format(time.RFC3339))

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
			null.Time{Time: time.Time{}, Valid: false},
		},
		{
			"filled value",
			time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC),
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC), Valid: true},
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
			null.Time{Time: time.Time{}, Valid: false},
		},
		{
			"invalid value",
			nil,
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC), Valid: false},
		},
		{
			"filled value",
			time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC),
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC), Valid: true},
		},
		{
			"blank valid value",
			time.Time{},
			null.Time{Time: time.Time{}, Valid: true},
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
			null.Time{Time: time.Time{}, Valid: false},
		},
		{
			"filled value",
			time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC),
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC), Valid: true},
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
			null.Time{Time: time.Time{}, Valid: false},
		},
		{
			"invalid filled value",
			[]byte(`null`),
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC), Valid: false},
		},
		{
			"filled value",
			[]byte(`"2019-01-01T12:00:00Z"`),
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC), Valid: true},
		},
		{
			"blank valid value",
			[]byte(`"0001-01-01T00:00:00Z"`),
			null.Time{Time: time.Time{}, Valid: true},
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
			null.Time{Time: time.Time{}, Valid: false},
		},
		{
			"filled string",
			[]byte(`"2019-01-01T12:00:00Z"`),
			null.Time{Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC), Valid: true},
		},
		{
			"blank string",
			[]byte(`"0001-01-01T00:00:00Z"`),
			null.Time{Time: time.Time{}, Valid: false},
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

func TestTime_TimeOrZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		in   null.Time
		exp  time.Time
	}{
		{
			"valid value",
			null.Time{Valid: true, Time: time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC)},
			time.Date(2019, 01, 01, 12, 00, 00, 0, time.UTC),
		},
		{
			"invalid value",
			null.Time{Valid: false, Time: time.Time{}},
			time.Time{},
		},
	}

	for _, testCase := range tests {
		tc := testCase
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			got := tc.in.TimeOrZero()
			if !reflect.DeepEqual(tc.exp, got) {
				t.Errorf("expected `%v`, got `%v`", tc.exp, got)
			}
		})
	}
}
