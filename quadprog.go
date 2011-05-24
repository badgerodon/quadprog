package quadprog

/*
#include <stdlib.h>
#include <math.h>
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
*/
import "C"
import (
	"fmt"
	"os"
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

type (
	Matrix struct {
		elements []float64
		rows, cols int
	}
	
	Vector []float64
)

func NewMatrix(rows, cols int) Matrix {
	return Matrix{make([]float64, rows*cols),rows,cols}
}
func NewIdentity(size int) Matrix {
	m := NewMatrix(size, size)
	for i := 0; i < size; i++ {
		m.Set(i, i, 1.0)
	}
	return m
}

func (this Matrix) Get(row, col int) float64 {
	return this.elements[row * this.cols + col]
}
func (this Matrix) Set(row, col int, value float64) {
	this.elements[row * this.cols + col] = value
}
func (this Matrix) String() string {
	str := ""
	for i := 0; i < this.rows; i++ {
		if i > 0 {		
			str += "\n"
		}
		for j := 0; j < this.cols; j++ {
			str += fmt.Sprintf("%6.2f", this.Get(i,j))
		}
	}
	return str
}

func NewVector(size int) Vector {
	return make([]float64, size)
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
func Solve(D Matrix, d Vector, A1 Matrix, b1 Vector, A2 Matrix, b2 Vector) (Vector, os.Error) {
	// The algorithm expects a combined "A" matrix. So build it
	A := NewMatrix(A1.rows + A2.rows, max(A2.cols, A1.cols))
	b := NewVector(len(b1) + len(b2))
	meq := len(b1)
	
	for i := 0; i < A1.rows; i++ {
		for j := 0; j < A1.cols; j++ {
			A.Set(i, j, A1.Get(i, j))
		}
	}
	
	for i := 0; i < A2.rows; i++ {
		for j := 0; j < A2.cols; j++ {
			A.Set(meq+i, j, A2.Get(i, j))
		}
	}
	
	for i := 0; i < len(b1); i++ {
		b[i] = b1[i]
	}
	
	for i := 0; i < len(b2); i++ {
		b[meq + i] = b2[i]
	}
	
	error := 0
	n := D.rows
	q := A.rows
	
	if n != D.cols {
		return nil, os.NewError("The D Matrix is not symmetric")
	}
	if n != len(d) {
		return nil, os.NewError("The D Matrix and the D Vector are incompatible")
	}
	if n != A.cols {
		return nil, os.NewError("The A Matrix and the D Vector are incompatible")
	}
	if q != len(b) {
		return nil, os.NewError("The A Matrix and the B Vector are incompatible")
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
	work := NewVector(2 * n + r * (r + 5) / 2 + 2 * q + 1)
	iter := make([]int, 2)
	
	run(D.elements,
		d,
		&n,
		&n,
		sol,
		lagr,
		&crval,
		A.elements,
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
		return sol, os.NewError("The supplied constraints are inconsistent")
	} else if error == 2 {
		return sol, os.NewError("The D matrix must be positive-definite in order for this algorithm to work.")
	}
	
	return sol, nil
}