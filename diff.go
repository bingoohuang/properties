package properties

type ChangeType int

const (
	Modified ChangeType = iota
	Added
	Removed
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
			if v != lv {
				f(DiffEvent{ChangeType: Modified, Key: k, LeftValue: lv, RightValue: v})
			}
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