package janitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWeightLoadBalancer(t *testing.T) {
	wb := NewWeightLoadBalancer()
	assert.NotNil(t, wb)
}

func TestNewWeightLoadBalancerSeed1(t *testing.T) {
	wb := NewWeightLoadBalancer()

	t1 := &Target{Weight: 40}
	t2 := &Target{Weight: 60}
	targets := []*Target{
		t1, t2,
	}

	t1Times := 0
	for i := 0; i < 10000; i++ {
		if wb.Seed(targets) == t1 {
			t1Times += 1
		}
	}

	assert.True(t, t1Times < 4500)
}

func TestNewWeightLoadBalancerSeed2(t *testing.T) {
	wb := NewWeightLoadBalancer()

	t1 := &Target{Weight: 1}
	t2 := &Target{Weight: 100}
	targets := []*Target{
		t1, t2,
	}

	t1Times := 0
	for i := 0; i < 10000; i++ {
		if wb.Seed(targets) == t1 {
			t1Times += 1
		}
	}

	assert.True(t, t1Times < 200)
	assert.True(t, t1Times > 0)
}

func TestNewWeightLoadBalancerSeed3(t *testing.T) {
	wb := NewWeightLoadBalancer()

	t1 := &Target{Weight: 1}
	t2 := &Target{Weight: 100}
	targets := []*Target{
		t1, t2,
	}

	t1Times := 0
	for i := 0; i < 10000; i++ {
		if wb.Seed(targets) == t1 {
			t1Times += 1
		}
	}

	assert.True(t, t1Times < 200)
	assert.True(t, t1Times > 0)
}
