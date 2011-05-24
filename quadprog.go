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
import "unsafe"

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