package exhaustive

import (
	"errors"
	"flag"
	"io"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestRepeatFlag(t *testing.T) {
	type testcase struct {
		caption string
		args    []string // input arguments
		want    []int
		err     string
	}

	var (
		even = func(v string) (int, error) {
			n, err := strconv.Atoi(v)
			if err != nil {
				return 0, err
			}
			if n%2 == 0 {
				return n, nil
			}
			return 0, errors.New("not even")
		}
		testcases = []testcase{
			{"multiple", []string{"-a", "10", "-a=100", "-a", "1000"}, []int{10, 100, 1000}, ""},
			{"single", []string{"-a", "10"}, []int{10}, ""},
			{"empty", nil, nil, ""},
			{"error", []string{"-a", "2", "-a", "3", "-a", "4"}, []int{2}, "not even"},
		}
	)

	for _, tt := range testcases {
		t.Run(tt.caption, func(t *testing.T) {
			fs := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
			fs.SetOutput(io.Discard)
			f := repeatFlag[int]{set: even}
			fs.Var(&f, "a", "...")

			err := fs.Parse(tt.args)
			switch {
			case err == nil && tt.err != "":
				t.Errorf("got nil error, want error suffix %q", tt.err)
			case err != nil && tt.err == "":
				t.Errorf("got error %q, want nil", err)
			case err != nil && !strings.HasSuffix(err.Error(), tt.err):
				t.Errorf("got error %q, want error suffix %q", err, tt.err)
			case !slices.Equal(f.vals, tt.want):
				t.Errorf("got: %v, want: %v", f.vals, tt.want)
			}
		})
	}

}
