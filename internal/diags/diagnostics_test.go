package diags

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	type diagFlat struct {
		Severity Severity
		Summary  string
		Detail   string
	}

	tests := map[string]struct {
		Cons func(Diagnostics) Diagnostics
		Want []diagFlat
	}{
		"nil": {
			func(diags Diagnostics) Diagnostics {
				return diags
			},
			nil,
		},
		"fmt.Errorf": {
			func(diags Diagnostics) Diagnostics {
				diags = diags.Append(fmt.Errorf("oh no bad"))
				return diags
			},
			[]diagFlat{
				{
					Severity: Error,
					Summary:  "oh no bad",
				},
			},
		},
		"errors.New": {
			func(diags Diagnostics) Diagnostics {
				diags = diags.Append(errors.New("oh no bad"))
				return diags
			},
			[]diagFlat{
				{
					Severity: Error,
					Summary:  "oh no bad",
				},
			},
		},
		"concat Diagnostics": {
			func(diags Diagnostics) Diagnostics {
				var moreDiags Diagnostics
				moreDiags = moreDiags.Append(errors.New("bad thing A"))
				moreDiags = moreDiags.Append(errors.New("bad thing B"))
				return diags.Append(moreDiags)
			},
			[]diagFlat{
				{
					Severity: Error,
					Summary:  "bad thing A",
				},
				{
					Severity: Error,
					Summary:  "bad thing B",
				},
			},
		},
		"single Diagnostic": {
			func(diags Diagnostics) Diagnostics {
				return diags.Append(SimpleWarning("Don't forget your toothbrush!"))
			},
			[]diagFlat{
				{
					Severity: Warning,
					Summary:  "Don't forget your toothbrush!",
				},
			},
		},
		"multiple appends": {
			func(diags Diagnostics) Diagnostics {
				diags = diags.Append(SimpleWarning("Don't forget your toothbrush!"))
				diags = diags.Append(fmt.Errorf("exploded"))
				return diags
			},
			[]diagFlat{
				{
					Severity: Warning,
					Summary:  "Don't forget your toothbrush!",
				},
				{
					Severity: Error,
					Summary:  "exploded",
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotDiags := test.Cons(nil)
			var got []diagFlat
			for _, item := range gotDiags {
				desc := item.Description()
				got = append(got, diagFlat{
					Severity: item.Severity(),
					Summary:  desc.Summary,
					Detail:   desc.Detail,
				})
			}

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("wrong result\ngot: %swant: %s", got, test.Want)
			}
		})
	}
}

func TestDiagnosticsErr(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var diags Diagnostics
		err := diags.Err()
		if err != nil {
			t.Errorf("got non-nil error %#v; want nil", err)
		}
	})
	t.Run("warning only", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(SimpleWarning("bad"))
		err := diags.Err()
		if err != nil {
			t.Errorf("got non-nil error %#v; want nil", err)
		}
	})
	t.Run("one error", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		err := diags.Err()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		if got, want := err.Error(), "didn't work"; got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
	})
	t.Run("two errors", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		diags = diags.Append(errors.New("didn't work either"))
		err := diags.Err()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		want := strings.TrimSpace(`
2 problems:

- didn't work
- didn't work either
`)
		if got := err.Error(); got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
	})
	t.Run("error and warning", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		diags = diags.Append(SimpleWarning("didn't work either"))
		err := diags.Err()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		// Since this "as error" mode is just a fallback for
		// non-diagnostics-aware situations like tests, we don't actually
		// distinguish warnings and errors here since the point is to just
		// get the messages rendered. User-facing code should be printing
		// each diagnostic separately, so won't enter this codepath,
		want := strings.TrimSpace(`
2 problems:

- didn't work
- didn't work either
`)
		if got := err.Error(); got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
	})
}

func TestDiagnosticsErrWithWarnings(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var diags Diagnostics
		err := diags.ErrWithWarnings()
		if err != nil {
			t.Errorf("got non-nil error %#v; want nil", err)
		}
	})
	t.Run("warning only", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(SimpleWarning("bad"))
		err := diags.ErrWithWarnings()
		if err == nil {
			t.Errorf("got nil error; want NonFatalError")
			return
		}
		if _, ok := err.(NonFatalError); !ok {
			t.Errorf("got %T; want NonFatalError", err)
		}
	})
	t.Run("one error", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		err := diags.ErrWithWarnings()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		if got, want := err.Error(), "didn't work"; got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
	})
	t.Run("two errors", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		diags = diags.Append(errors.New("didn't work either"))
		err := diags.ErrWithWarnings()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		want := strings.TrimSpace(`
2 problems:

- didn't work
- didn't work either
`)
		if got := err.Error(); got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
	})
	t.Run("error and warning", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		diags = diags.Append(SimpleWarning("didn't work either"))
		err := diags.ErrWithWarnings()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		// Since this "as error" mode is just a fallback for
		// non-diagnostics-aware situations like tests, we don't actually
		// distinguish warnings and errors here since the point is to just
		// get the messages rendered. User-facing code should be printing
		// each diagnostic separately, so won't enter this codepath,
		want := strings.TrimSpace(`
2 problems:

- didn't work
- didn't work either
`)
		if got := err.Error(); got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
	})
}

func TestDiagnosticsNonFatalErr(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var diags Diagnostics
		err := diags.NonFatalErr()
		if err != nil {
			t.Errorf("got non-nil error %#v; want nil", err)
		}
	})
	t.Run("warning only", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(SimpleWarning("bad"))
		err := diags.NonFatalErr()
		if err == nil {
			t.Errorf("got nil error; want NonFatalError")
			return
		}
		if _, ok := err.(NonFatalError); !ok {
			t.Errorf("got %T; want NonFatalError", err)
		}
	})
	t.Run("one error", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		err := diags.NonFatalErr()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		if got, want := err.Error(), "didn't work"; got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
		if _, ok := err.(NonFatalError); !ok {
			t.Errorf("got %T; want NonFatalError", err)
		}
	})
	t.Run("two errors", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		diags = diags.Append(errors.New("didn't work either"))
		err := diags.NonFatalErr()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		want := strings.TrimSpace(`
2 problems:

- didn't work
- didn't work either
`)
		if got := err.Error(); got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
		if _, ok := err.(NonFatalError); !ok {
			t.Errorf("got %T; want NonFatalError", err)
		}
	})
	t.Run("error and warning", func(t *testing.T) {
		var diags Diagnostics
		diags = diags.Append(errors.New("didn't work"))
		diags = diags.Append(SimpleWarning("didn't work either"))
		err := diags.NonFatalErr()
		if err == nil {
			t.Fatalf("got nil error %#v; want non-nil", err)
		}
		// Since this "as error" mode is just a fallback for
		// non-diagnostics-aware situations like tests, we don't actually
		// distinguish warnings and errors here since the point is to just
		// get the messages rendered. User-facing code should be printing
		// each diagnostic separately, so won't enter this codepath,
		want := strings.TrimSpace(`
2 problems:

- didn't work
- didn't work either
`)
		if got := err.Error(); got != want {
			t.Errorf("wrong error message\ngot:  %s\nwant: %s", got, want)
		}
		if _, ok := err.(NonFatalError); !ok {
			t.Errorf("got %T; want NonFatalError", err)
		}
	})
}
