package properties

// ChangeType defines the type of chaging.
type ChangeType int

const (
	// Modified ...
	Modified ChangeType = iota
	// Added ...
	Added
	// Removed ...
	Removed
	// Same ...
	Same
)

// DiffEvent defines ChangeEvent for properties diff
type DiffEvent struct {
	ChangeType ChangeType
	Key        string
	LeftValue  string
	RightValue string
}

// Diff diffs l to r.
func Diff(l, r *Doc, f func(DiffEvent)) {
	lm := make(map[string]string)

	l.Foreach(func(v, k string) bool { lm[k] = v; return true })
	r.Foreach(func(v, k string) bool {
		if lv, ok := lm[k]; ok {
			typ := Same
			if v != lv {
				typ = Modified
			}
			f(DiffEvent{ChangeType: typ, Key: k, LeftValue: lv, RightValue: v})
			delete(lm, k)
		} else {
			f(DiffEvent{ChangeType: Added, Key: k, LeftValue: "", RightValue: v})
		}

		return true
	})

	for k, v := range lm {
		f(DiffEvent{ChangeType: Removed, Key: k, LeftValue: v, RightValue: ""})
	}
}
