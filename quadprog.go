package quadprog

/*
extern void qpgen1(
	double *dmat, double *dvec, int *fddmat, int *n,
	double *sol, double *lagr, double *crval,
	double *amat, int *iamat, double *bvec, int *fdamat, int *q,
	int *meq, int *iact, int *nact, int *iter,
	double *work, int *ierr
);

extern void qpgen2_(
	double *dmat, double *dvec, int *fddmat, int *n,
	double *sol, double *lagr, double *crval,
	double *amat, double *bvec, int *fdamat, int *q,
	int *meq, int *iact, int *nact, int *iter,
	double *work, int *ierr
);

extern void aind(
	int *ind, int *m, int *q, int *n, int *ok
);

#cgo windows LDFLAGS: lib/windows/amd64/aind.o lib/windows/amd64/daxpy.o lib/windows/amd64/ddot.o lib/windows/amd64/dpofa.o lib/windows/amd64/dscal.o lib/windows/amd64/solve.QP.o lib/windows/amd64/util.o
*/
import "C"
import (
	"fmt"
	. "github.com/badgerodon/lalg"
	"unsafe"
)

func mkd(arr []float64) *C.double {
	return (*C.double)(unsafe.Pointer(&arr[0]))
}

func mki(i *int) *C.int {
	return (*C.int)(unsafe.Pointer(i))
}

func run(dmat []float64,
	dvec []float64,
	fddmat *int,
	n *int,
	sol []float64,
	lagr []float64,
	crval *float64,
	amat []float64,
	bvec []float64,
	fdamat *int,
	q *int,
	meq *int,
	iact []int,
	nact *int,
	iter []int,
	work []float64,
	ierr *int) {

	C.qpgen2_(mkd(dmat),
		mkd(dvec),
		mki(fddmat),
		mki(n),
		mkd(sol),
		mkd(lagr),
		(*C.double)(unsafe.Pointer(crval)),
		mkd(amat),
		mkd(bvec),
		mki(fdamat),
		mki(q),
		mki(meq),
		(*C.int)(unsafe.Pointer(&iact[0])),
		mki(nact),
		(*C.int)(unsafe.Pointer(&iter[0])),
		mkd(work),
		mki(ierr),
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/*
From the original documentation:

this routine uses the Goldfarb/Idnani algorithm to solve the following
minimization problem:

   minimize  -d^T x + 1/2 *  x^T D x
   where   A1^T x  = b1
           A2^T x >= b2

the matrix D is assumed to be positive definite.  Especially, w.l.o.g. D is
assumed to be symmetric.
*/
func Solve(D Matrix, d Vector, A1 Matrix, b1 Vector, A2 Matrix, b2 Vector) (Vector, error) {
	// The algorithm expects a combined "A" matrix. So build it
	A := NewMatrix(A1.Rows+A2.Rows, max(A2.Cols, A1.Cols))
	b := NewVector(len(b1) + len(b2))
	meq := len(b1)

	for i := 0; i < A1.Rows; i++ {
		for j := 0; j < A1.Cols; j++ {
			A.Set(i, j, A1.Get(i, j))
		}
	}

	for i := 0; i < A2.Rows; i++ {
		for j := 0; j < A2.Cols; j++ {
			A.Set(meq+i, j, A2.Get(i, j))
		}
	}

	for i := 0; i < len(b1); i++ {
		b[i] = b1[i]
	}

	for i := 0; i < len(b2); i++ {
		b[meq+i] = b2[i]
	}

	error := 0
	n := D.Rows
	q := A.Rows

	if n != D.Cols {
		return nil, fmt.Errorf("The D Matrix is not symmetric")
	}
	if n != len(d) {
		return nil, fmt.Errorf("The D Matrix and the D Vector are incompatible")
	}
	if n != A.Cols {
		return nil, fmt.Errorf("The A Matrix and the D Vector are incompatible")
	}
	if q != len(b) {
		return nil, fmt.Errorf("The A Matrix and the B Vector are incompatible")
	}

	iact := make([]int, q)
	nact := 0
	r := n
	if r > q {
		r = q
	}
	sol := NewVector(n)
	lagr := NewVector(q)
	crval := float64(0)
	work := NewVector(2*n + r*(r+5)/2 + 2*q + 1)
	iter := make([]int, 2)

	run(D.Elements,
		d,
		&n,
		&n,
		sol,
		lagr,
		&crval,
		A.Elements,
		b,
		&n,
		&q,
		&meq,
		iact,
		&nact,
		iter,
		work,
		&error,
	)

	if error == 1 {
		return sol, fmt.Errorf("The supplied constraints are inconsistent")
	} else if error == 2 {
		return sol, fmt.Errorf("The D matrix must be positive-definite in order for this algorithm to work.")
	}

	return sol, nil
}
