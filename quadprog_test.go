package quadprog

import (
	"fmt"
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

	fmt.Print(sol, err)
}
