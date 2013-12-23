package quadprog

import (
	. "github.com/badgerodon/lalg"
	"testing"
)

func Test(t *testing.T) {

	D := NewMatrix(4, 4)
	D.Set(0, 0, 1)
	D.Set(1, 1, 7)
	D.Set(2, 2, 1)
	D.Set(3, 3, 1)

	d := NewVector(4)
	d[0] = 0.25
	d[1] = 0.25
	d[2] = 0.25
	d[3] = 0.25

	A1 := NewIdentity(0)
	b1 := NewVector(0)

	A2 := NewIdentity(4)
	b2 := NewVector(4)
	b2[0] = 0
	b2[1] = 0
	b2[2] = 0
	b2[3] = 0

	sol, err := Solve(D, d, A1, b1, A2, b2)

	if err != nil {
		t.Errorf("expected no error got %v", err)
	}

	if len(sol) != 4 {
		t.Errorf("expected 4 results got %v", len(sol))
	}

	if sol[0] != 0.25 {
		t.Errorf("expected w0 to be 0.25, got %v", sol[0])
	}
	if sol[1] != 0.03571428571428571 {
		t.Errorf("expected w1 to be 0.03571428571428571, got %v", sol[1])
	}
	if sol[2] != 0.25 {
		t.Errorf("expected w2 to be 0.25, got %v", sol[2])
	}
	if sol[3] != 0.25 {
		t.Errorf("expected w3 to be 0.25, got %v", sol[3])
	}
}
