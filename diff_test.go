package properties

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	events := make([]DiffEvent, 0)

	l, _ := LoadString("k1=v1\nk2=v2\nk4=v4")
	r, _ := LoadString("k1=v10\nk3=v3\nk4=v4")

	Diff(l, r, func(event DiffEvent) {
		events = append(events, event)
	})

	assert.Equal(t, []DiffEvent{
		{ChangeType: Modified, Key: "k1", LeftValue: "v1", RightValue: "v10"},
		{ChangeType: Added, Key: "k3", LeftValue: "", RightValue: "v3"},
		{ChangeType: Same, Key: "k4", LeftValue: "v4", RightValue: "v4"},
		{ChangeType: Removed, Key: "k2", LeftValue: "v2", RightValue: ""},
	}, events)
}
