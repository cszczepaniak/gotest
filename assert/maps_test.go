package assert

import "testing"

func TestMapHasKey(t *testing.T) {
	t.Run("existing key passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		m := map[string]int{"a": 1, "b": 2}
		MapHasKey(f, m, "a")
		if f.fatalCalled {
			t.Errorf("MapHasKey should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("missing key fails", func(t *testing.T) {
		f := &fakeTB{TB: t}
		m := map[string]int{"a": 1, "b": 2}
		MapHasKey(f, m, "c")
		if !f.fatalCalled {
			t.Error("MapHasKey with missing key should have called Fatal")
		}
		if f.fatalMsg == "" || len(f.fatalMsg) < 10 {
			t.Errorf("Fatal message should mention the key, got: %q", f.fatalMsg)
		}
	})
}

func TestMapLen(t *testing.T) {
	t.Run("matching length passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		MapLen(f, map[string]int{"a": 1, "b": 2}, 2)
		if f.fatalCalled {
			t.Errorf("MapLen should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("empty map length 0 passes", func(t *testing.T) {
		f := &fakeTB{TB: t}
		MapLen(f, map[string]int{}, 0)
		if f.fatalCalled {
			t.Errorf("MapLen should not have called Fatal, got: %s", f.fatalMsg)
		}
	})

	t.Run("wrong length fails", func(t *testing.T) {
		f := &fakeTB{TB: t}
		MapLen(f, map[string]int{"a": 1, "b": 2}, 3)
		if !f.fatalCalled {
			t.Error("MapLen with wrong length should have called Fatal")
		}
	})
}
