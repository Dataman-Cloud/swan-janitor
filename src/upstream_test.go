package janitor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUpstream(t *testing.T) {
	u := NewUpstream()
	assert.NotNil(t, u)
}

func TestEqual(t *testing.T) {
	u1 := NewUpstream()
	u1.AppID = "foobar"

	u2 := NewUpstream()
	u2.AppID = "foobar"

	assert.True(t, u1.Equal(u2))
}

func TestAddTarget(t *testing.T) {
	u := NewUpstream()
	ta := &Target{
		TaskID: "foobar",
	}

	u.AddTarget(ta)
	assert.Equal(t, 1, len(u.Targets))

}

func TestContainsTarget(t *testing.T) {
	u := NewUpstream()
	ta := &Target{
		TaskID: "foobar",
	}

	u.AddTarget(ta)

	assert.True(t, u.ContainsTarget("foobar"))
}

func TestRemoveTarget(t *testing.T) {
	u := NewUpstream()
	ta := &Target{
		TaskID: "foobar",
	}

	u.AddTarget(ta)

	assert.True(t, u.ContainsTarget("foobar"))
	u.RemoveTarget(ta)
	assert.False(t, u.ContainsTarget("foobar"))
}

func TestUpdateTargetWeight(t *testing.T) {
	u := NewUpstream()
	ta := &Target{
		TaskID: "foobar",
		Weight: 100,
	}

	u.AddTarget(ta)
	assert.Equal(t, float64(100), u.GetTarget("foobar").Weight)

	u.UpdateTargetWeight("foobar", 101)
	assert.Equal(t, float64(101), u.GetTarget("foobar").Weight)
}

func TestNextTargetEntry(t *testing.T) {
	u := NewUpstream()
	ta := &Target{
		TaskID:   "foobar",
		Weight:   100,
		TaskIP:   "0.0.0.0",
		TaskPort: 1023,
	}

	ta1 := &Target{
		TaskID:   "foobar",
		Weight:   400,
		TaskIP:   "0.0.0.0",
		TaskPort: 1024,
	}

	u.AddTarget(ta)
	u.AddTarget(ta1)
	assert.Equal(t, 2, len(u.Targets))

	taMeet := 0
	for i := 0; i < 10000; i++ {
		if strings.HasSuffix(u.NextTargetEntry().Host, "1023") {
			taMeet += 1
		}
	}

	assert.True(t, taMeet < 2500)

}
