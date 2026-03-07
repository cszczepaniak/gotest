package assert

import (
	"cmp"
	"testing"
)

func TestElementsMatch(t *testing.T) {
	cmpInt := cmp.Compare[int]

	t.Run("same elements same order pass", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceElemsMatch(f, []int{1, 2, 3}, []int{1, 2, 3})
		SliceElemsMatchFunc(f, []int{1, 2, 3}, []int{1, 2, 3}, cmpInt)
		if f.fatalCalled {
			t.Errorf("ElementsMatch should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("matchfunc with different ordering works", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceElemsMatchFunc(f, []int{3, 1, 2}, []int{1, 2, 3}, func(a, b int) int {
			return -cmp.Compare(a, b)
		})
		if f.fatalCalled {
			t.Errorf("ElementsMatch should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("same elements different order pass", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceElemsMatch(f, []int{3, 1, 2}, []int{1, 2, 3})
		SliceElemsMatchFunc(f, []int{3, 1, 2}, []int{1, 2, 3}, cmpInt)
		if f.fatalCalled {
			t.Errorf("ElementsMatch should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("different elements fail", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceElemsMatch(f, []int{1, 2, 3}, []int{1, 2, 4})
		SliceElemsMatchFunc(f, []int{1, 2, 3}, []int{1, 2, 4}, cmpInt)
		if !f.fatalCalled {
			t.Error("ElementsMatch with different elements should have called Fatal")
		}
	})

	t.Run("different lengths fail", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceElemsMatch(f, []int{1, 2}, []int{1, 2, 3})
		SliceElemsMatchFunc(f, []int{1, 2}, []int{1, 2, 3}, cmpInt)
		if !f.fatalCalled {
			t.Error("ElementsMatch with different lengths should have called Fatal")
		}
	})
}

func TestSliceLen(t *testing.T) {
	t.Run("matching length passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceLen(f, []int{1, 2, 3}, 3)
		if f.fatalCalled {
			t.Errorf("SliceLen should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("empty slice length 0 passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceLen(f, []int{}, 0)
		if f.fatalCalled {
			t.Errorf("SliceLen should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("nil slice length 0 passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceLen(f, []int(nil), 0)
		if f.fatalCalled {
			t.Errorf("SliceLen should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("wrong length fails", func(t *testing.T) {
		f := &fakeTB{TB: t}
		SliceLen(f, []int{1, 2, 3}, 2)
		if !f.fatalCalled {
			t.Error("SliceLen with wrong length should have called Fatal")
		}
	})
}
