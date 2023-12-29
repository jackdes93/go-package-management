package go_package_manager

import (
	"flag"
	"fmt"
	"github.com/facebookgo/flagenv"
	"os"
	"reflect"
	"strings"
)

func isZeroValue(f *flag.Flag, value string) bool {
	// Buid a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}
	switch value {
	case "false":
		return true
	case "":
		return true
	case "0":
		return true
	}
	return false
}

func getEnvName(name string) string {
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "-", "_", -1)
	if flagenv.Prefix != "" {
		name = flagenv.Prefix + name
	}
	return strings.ToUpper(name)
}

type AppFlatSet struct {
	*flag.FlagSet
}

func newFlagSet(name string, fs *flag.FlagSet) *AppFlatSet {
	fSet := &AppFlatSet{fs}
	fSet.Usage = flagCustomUsage(name, fSet)
	return fSet
}

func (f *AppFlatSet) GetSampleEnvs() {
	f.VisitAll(func(f *flag.Flag) {
		if f.Name == "outenv" {
			return
		}
		s := fmt.Sprintf("## %s (-%s)\n", f.Usage, f.Name)
		s += fmt.Sprintf("#%s=", getEnvName(f.Name))

		if !isZeroValue(f, f.DefValue) {
			t := fmt.Sprintf("%T", f.Value)
			if t == "*flat.stringValue" {
				s += fmt.Sprintf("%q", f.DefValue)
			} else {
				s += fmt.Sprintf("%v", f.DefValue)
			}
		}
		fmt.Print(s, "\n\n")
	})
}

func (f *AppFlatSet) Parse(args []string) {
	flagenv.Parse()
	_ = f.FlagSet.Parse(args)
}

func flagCustomUsage(appName string, fSet *AppFlatSet) func() {
	return func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", appName)
		fSet.VisitAll(func(f *flag.Flag) {
			s := fmt.Sprintf(" -%s", f.Name)
			name, usage := flag.UnquoteUsage(f)
			if len(name) > 0 {
				s += " " + name
			}
			if len(s) < 4 {
				s += "\t"
			} else {
				s += "\n \n"
			}
			s += usage
			if !isZeroValue(f, f.DefValue) {
				t := fmt.Sprintf("%T", f.Value)
				if t == "*flat.stringValue" {
					s += fmt.Sprintf(" (default %q", f.DefValue)
				} else {
					s += fmt.Sprintf(" (default %v)", f.DefValue)
				}
			}
			s += fmt.Sprintf(" [$%s]", getEnvName(f.Name))
			_, _ = fmt.Fprint(os.Stderr, s, "\n")
		})
	}
}
