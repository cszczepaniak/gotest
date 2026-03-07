package assert

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// fakeTB implements testing.TB and records whether Fatal was called.
type fakeTB struct {
	testing.TB
	fatalCalled bool
	fatalMsg    string
}

func (f *fakeTB) Fatal(args ...any) {
	f.fatalCalled = true
	f.fatalMsg = fmt.Sprint(args...)
}

func (f *fakeTB) Fatalf(format string, args ...any) {
	f.fatalCalled = true
	f.fatalMsg = fmt.Sprintf(format, args...)
}

func (f *fakeTB) Helper() {}

func (f *fakeTB) reset() {
	f.fatalCalled = false
	f.fatalMsg = ""
}

func TestEqual(t *testing.T) {
	t.Run("equal values pass", func(t *testing.T) {
		f := &fakeTB{TB: t}
		Equal(f, 42, 42)
		if f.fatalCalled {
			t.Errorf("Equal(42, 42) should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("equal strings pass", func(t *testing.T) {
		f := &fakeTB{TB: t}
		Equal(f, "hello", "hello")
		if f.fatalCalled {
			t.Errorf("Equal(hello, hello) should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("unequal values fail", func(t *testing.T) {
		f := &fakeTB{TB: t}
		Equal(f, 42, 43)
		if !f.fatalCalled {
			t.Error("Equal(42, 43) should have called Fatal")
		}
		if f.fatalMsg == "" {
			t.Error("Fatal message should not be empty")
		}
	})

	t.Run("slices", func(t *testing.T) {
		f := &fakeTB{TB: t}
		Equal(f, []int{1, 2, 3}, []int{1, 2, 3})
		if f.fatalCalled {
			t.Errorf("Equal slices should not have called Fatal, got: %s", f.fatalMsg)
		}
		f.reset()
		Equal(f, []int{1, 2, 3}, []int{1, 2, 4})
		if !f.fatalCalled {
			t.Error("Equal different slices should have called Fatal")
		}
	})
}

func TestError(t *testing.T) {
	t.Run("non-nil error passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		Error(f, errors.New("something went wrong"))
		if f.fatalCalled {
			t.Errorf("Error with non-nil err should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("nil error fails", func(t *testing.T) {
		f := &fakeTB{TB: t}
		Error(f, nil)
		if !f.fatalCalled {
			t.Error("Error with nil should have called Fatal")
		}
	})
}

func TestNoError(t *testing.T) {
	t.Run("nil error passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		NoError(f, nil)
		if f.fatalCalled {
			t.Errorf("NoError with nil should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("non-nil error fails", func(t *testing.T) {
		f := &fakeTB{TB: t}
		err := errors.New("oops")
		NoError(f, err)
		if !f.fatalCalled {
			t.Error("NoError with non-nil err should have called Fatal")
		}
		if f.fatalMsg == "" || len(f.fatalMsg) < 5 {
			t.Errorf("Fatal message should include the error, got: %q", f.fatalMsg)
		}
	})
}

var errSentinel = errors.New("sentinel")

func TestErrorIs(t *testing.T) {
	t.Run("exact match passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		ErrorIs(f, errSentinel, errSentinel)
		if f.fatalCalled {
			t.Errorf("ErrorIs with same error should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("wrapped error matches passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		wrapped := fmt.Errorf("wrapped: %w", errSentinel)
		ErrorIs(f, wrapped, errSentinel)
		if f.fatalCalled {
			t.Errorf("ErrorIs with wrapped error should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("different error fails", func(t *testing.T) {
		f := &fakeTB{TB: t}
		other := errors.New("other")
		ErrorIs(f, other, errSentinel)
		if !f.fatalCalled {
			t.Error("ErrorIs with different error should have called Fatal")
		}
	})

	t.Run("nil error fails", func(t *testing.T) {
		f := &fakeTB{TB: t}
		ErrorIs(f, nil, errSentinel)
		if !f.fatalCalled {
			t.Error("ErrorIs with nil err should have called Fatal")
		}
	})

	t.Run("joined error containing target passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		joined := errors.Join(errors.New("first"), errSentinel, errors.New("third"))
		ErrorIs(f, joined, errSentinel)
		if f.fatalCalled {
			t.Errorf("ErrorIs with joined error containing target should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("joined error not containing target prints full chain", func(t *testing.T) {
		f := &fakeTB{TB: t}
		errA := errors.New("a")
		errB := errors.New("b")
		joined := errors.Join(errA, errB)
		ErrorIs(f, joined, errSentinel)
		if !f.fatalCalled {
			t.Fatal("ErrorIs with joined error missing target should have called Fatal")
		}
		// Message should show the joined structure: [1]: and [2]: for each branch
		if !strings.Contains(f.fatalMsg, "[1]:") || !strings.Contains(f.fatalMsg, "[2]:") {
			t.Errorf("Fatal message should show joined branches [1]: and [2]:, got:\n%s", f.fatalMsg)
		}
		if !strings.Contains(f.fatalMsg, "a (") || !strings.Contains(f.fatalMsg, "b (") {
			t.Errorf("Fatal message should show both joined errors, got:\n%s", f.fatalMsg)
		}
	})
}
