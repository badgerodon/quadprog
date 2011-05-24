include $(GOROOT)/src/Make.inc

#QP_LIB=lib/solve.a
QP_OFILES=lib/aind.o lib/daxpy.o lib/ddot.o lib/dpofa.o lib/dscal.o lib/util.o lib/solve.QP.o

TARG=github.com/badgerodon/quadprog
CGOFILES=quadprog.go
CGO_OFILES=$(QP_OFILES)
CGO_LDFLAGS=-lm -lgfortran

CLEANFILES+=$(QP_OFILES)

include $(GOROOT)/src/Make.pkg

$(QP_OFILES): 
	cd lib && gfortran -c *.f