package pluginapp

import (
	"flag"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type TOptions struct {
	Bool     *bool
	Int      *int
	Uint     *uint
	Int64    *int64
	Uint64   *uint64
	Float64  *float64
	Str      *string
	Duration *time.Duration

	BoolVar     bool
	IntVar      int
	UintVar     uint
	Int64Var    int64
	Uint64Var   uint64
	Float64Var  float64
	StringVar   string
	DurationVar time.Duration
}

func (o TOptions) String() string {
	return fmt.Sprintf("Bool: %v, Int: %v, Uint: %v, Int64: %v, Uint64: %v, Float64: %v, String: %v, Duration: %v "+
		"BoolVar: %v IntVar: %v, UintVar: %v, Int64Var: %v, Uint64Var: %v, Float64Var: %v, StringVar: %v, DurationVar: %v",
		*o.Bool, *o.Int, *o.Uint, *o.Int64, *o.Uint64, *o.Float64, *o.Str, *o.Duration,
		o.BoolVar, o.IntVar, o.UintVar, o.Int64Var, o.UintVar, o.Float64Var, o.StringVar, o.DurationVar)
}

func TNewBool(v bool) *bool {
	return &v
}

func TNewInt(v int) *int {
	return &v
}

func TNewUint(v uint) *uint {
	return &v
}

func TNewInt64(v int64) *int64 {
	return &v
}

func TNewUint64(v uint64) *uint64 {
	return &v
}

func TNewFloat64(v float64) *float64 {
	return &v
}

func TNewString(v string) *string {
	return &v
}

func TNewDuration(v time.Duration) *time.Duration {
	return &v
}

func TestFlagSet(t *testing.T) {
	fs := FlagSet{
		prefix:  "test.",
		flagset: flag.NewFlagSet("test", flag.ExitOnError),
	}

	var o TOptions

	o.Bool = fs.Bool("Bool", false, "Bool usage")
	o.Int = fs.Int("Int", 0, "Int usage")
	o.Uint = fs.Uint("Uint", 0, "Uint usage")
	o.Int64 = fs.Int64("Int64", 0, "Int64 usage")
	o.Uint64 = fs.Uint64("Uint64", 0, "Uint64 usage")
	o.Float64 = fs.Float64("Float64", 0, "Float64 usage")
	o.Str = fs.String("String", "", "String usage")
	o.Duration = fs.Duration("Duration", 0, "Duration usage")

	fs.BoolVar(&o.BoolVar, "BoolVar", false, "BoolVar usage")
	fs.IntVar(&o.IntVar, "IntVar", 0, "IntVar usage")
	fs.UintVar(&o.UintVar, "UintVar", 0, "UintVar usage")
	fs.Int64Var(&o.Int64Var, "Int64Var", 0, "Int64Var usage")
	fs.Uint64Var(&o.Uint64Var, "Uint64Var", 0, "Uint64Var usage")
	fs.Float64Var(&o.Float64Var, "Float64Var", 0, "Float64Var usage")
	fs.StringVar(&o.StringVar, "StringVar", "", "StringVar usage")
	fs.DurationVar(&o.DurationVar, "DurationVar", 0, "DurationVar usage")

	args := []string{
		//"-h",
		"-test.Bool",
		"-test.Int", "1",
		"-test.Uint", "2",
		"-test.Int64", "3",
		"-test.Uint64", "4",
		"-test.Float64", "5.6",
		"-test.String", "string",
		"-test.Duration", "1m11s",
		"-test.IntVar", "1",
		"-test.UintVar", "2",
		"-test.Int64Var", "3",
		"-test.Uint64Var", "4",
		"-test.Float64Var", "5.6",
		"-test.StringVar", "string",
		"-test.DurationVar", "1m11s",
	}
	want := TOptions{
		Bool:        TNewBool(true),
		Int:         TNewInt(1),
		Uint:        TNewUint(2),
		Int64:       TNewInt64(3),
		Uint64:      TNewUint64(4),
		Float64:     TNewFloat64(5.6),
		Str:         TNewString("string"),
		Duration:    TNewDuration(71 * time.Second),
		IntVar:      1,
		UintVar:     2,
		Int64Var:    3,
		Uint64Var:   4,
		Float64Var:  5.6,
		StringVar:   "string",
		DurationVar: 71 * time.Second,
	}
	fs.flagset.Parse(args)

	if got := o; !reflect.DeepEqual(got, want) {
		t.Errorf("options: got %v, want %v", got, want)
	} else {
		t.Logf("options: got %v", got)
	}
}
