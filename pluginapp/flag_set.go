package pluginapp

import (
	"flag"
	"time"
)

// FlagSet 类似flag.FlagSet, 插件定义命令行参数时使用
type FlagSet struct {
	prefix  string
	flagset *flag.FlagSet
}

func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	return f.flagset.Bool(f.prefix+name, value, usage)
}

func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	f.flagset.BoolVar(p, f.prefix+name, value, usage)
}

func (f *FlagSet) Int(name string, value int, usage string) *int {
	return f.flagset.Int(f.prefix+name, value, usage)
}

func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.flagset.IntVar(p, f.prefix+name, value, usage)
}

func (f *FlagSet) Int64(name string, value int64, usage string) *int64 {
	return f.flagset.Int64(f.prefix+name, value, usage)
}

func (f *FlagSet) Int64Var(p *int64, name string, value int64, usage string) {
	f.flagset.Int64Var(p, f.prefix+name, value, usage)
}

func (f *FlagSet) Uint(name string, value uint, usage string) *uint {
	return f.flagset.Uint(f.prefix+name, value, usage)
}

func (f *FlagSet) UintVar(p *uint, name string, value uint, usage string) {
	f.flagset.UintVar(p, f.prefix+name, value, usage)
}

func (f *FlagSet) Uint64(name string, value uint64, usage string) *uint64 {
	return f.flagset.Uint64(f.prefix+name, value, usage)
}

func (f *FlagSet) Uint64Var(p *uint64, name string, value uint64, usage string) {
	f.flagset.Uint64Var(p, f.prefix+name, value, usage)
}

func (f *FlagSet) Float64(name string, value float64, usage string) *float64 {
	return f.flagset.Float64(f.prefix+name, value, usage)
}

func (f *FlagSet) Float64Var(p *float64, name string, value float64, usage string) {
	f.flagset.Float64Var(p, f.prefix+name, value, usage)
}

func (f *FlagSet) String(name string, value string, usage string) *string {
	return f.flagset.String(f.prefix+name, value, usage)
}

func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.flagset.StringVar(p, f.prefix+name, value, usage)
}

func (f *FlagSet) Duration(name string, value time.Duration, usage string) *time.Duration {
	return f.flagset.Duration(f.prefix+name, value, usage)
}

func (f *FlagSet) DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	f.flagset.DurationVar(p, f.prefix+name, value, usage)
}

func (f *FlagSet) Var(value flag.Value, name string, usage string) {
	f.flagset.Var(value, f.prefix+name, usage)
}
