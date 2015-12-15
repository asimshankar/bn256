# Set USE_C= to disable use of the optimized assembly implementation and use
# the portable C implementation instead.

QHASM_FOR_LIB=-DQHASM
ifdef USE_C
	QHASM_FOR_LIB=
endif

CPP=g++
CPPFLAGS=-g -Wall -DCHECK

CC=gcc
CFLAGS=-std=c99 -O3 -fomit-frame-pointer
LFLAGS=-lm

all: as check c

c: bilintest-c \
	 speedtest-c

as: bilintest-as \
		speedtest-as \
		test_curvepoint_multiscalar-as \
		test_twistpoint_multiscalar-as

check: bilintest-check \
	     speedtest-check \
			 test_curvepoint_multiscalar-check \
			 test_twistpoint_multiscalar-check

bilintest-check: bilintest.c linefunction.c optate.c fpe.c fp2e.c fp6e.c fp12e.c curvepoint_fp.c twistpoint_fp2.c final_expo.c scalar.c parameters.c mul.c mydouble.c
	$(CPP) $(CPPFLAGS) -DNTESTS=20 -o $@ $^

bilintest-c: bilintest.c linefunction.c optate.c fpe.c fp2e.c fp6e.c fp12e.c curvepoint_fp.c twistpoint_fp2.c final_expo.c scalar.c parameters.c mul.c mydouble.c
	$(CC) $(CFLAGS) -DNTESTS=1000 -o $@ $^ $(LFLAGS)

bilintest-as: bilintest.c linefunction.c optate.c fpe.c fp2e.c fp6e.c fp12e.c curvepoint_fp.c twistpoint_fp2.c final_expo.c scalar.c parameters.c mul.c mydouble.c asfunctions.a
	$(CC) $(CFLAGS) -DQHASM -DNTESTS=1000000 -o $@ $^ $(LFLAGS)

speedtest-check: speedtest.c linefunction.c optate.c fpe.c fp2e.c fp6e.c fp12e.c curvepoint_fp.c twistpoint_fp2.c final_expo.c scalar.c parameters.c mul.c mydouble.c
	$(CPP) $(CPPFLAGS) -o $@ $^

speedtest-c: speedtest.c linefunction.c optate.c fpe.c fp2e.c fp6e.c fp12e.c curvepoint_fp.c twistpoint_fp2.c final_expo.c scalar.c parameters.c mul.c mydouble.c
	$(CC) $(CFLAGS) -o $@ $^ $(LFLAGS)

speedtest-as: speedtest.c linefunction.c optate.c fpe.c fp2e.c fp6e.c fp12e.c curvepoint_fp.c twistpoint_fp2.c final_expo.c scalar.c parameters.c mul.c mydouble.c asfunctions.a
	$(CC) $(CFLAGS) -DQHASM -o $@ $^ $(LFLAGS)

test_curvepoint_multiscalar-check: test_curvepoint_multiscalar.c fpe.c scalar.c parameters.c curvepoint_fp.c mul.c mydouble.c \
																		curvepoint_fp_multiscalar.c heap_rootreplaced.s index_heap.c scalar_sub_nored.s \
																		asfunctions.a
	$(CPP) $(CPPFLAGS) -o $@ $^

test_curvepoint_multiscalar-as: test_curvepoint_multiscalar.c fpe.c scalar.c parameters.c curvepoint_fp.c mul.c mydouble.c \
																curvepoint_fp_multiscalar.c heap_rootreplaced.s index_heap.c scalar_sub_nored.s \
																asfunctions.a
	$(CC) $(CFLAGS) -DQHASM -o $@ $^ $(LFLAGS)

test_twistpoint_multiscalar-check: test_twistpoint_multiscalar.c fpe.c fp2e.c scalar.c parameters.c twistpoint_fp2.c mul.c mydouble.c \
																	twistpoint_fp2_multiscalar.c heap_rootreplaced.s index_heap.c scalar_sub_nored.s \
																	asfunctions.a
	$(CPP) $(CPPFLAGS) -o $@ $^

test_twistpoint_multiscalar-as: test_twistpoint_multiscalar.c fpe.c fp2e.c scalar.c parameters.c twistpoint_fp2.c mul.c mydouble.c \
																	twistpoint_fp2_multiscalar.c heap_rootreplaced.s index_heap.c scalar_sub_nored.s \
																	asfunctions.a
	$(CC) $(CFLAGS) -DQHASM -o $@ $^ $(LFLAGS)


%.o: %.s
	$(CC) $(CFLAGS) -c -o $@ $^

%.o: %.c
	$(CC) $(CFLAGS) $(QHASM_FOR_LIB) -c -o $@ $^

asfunctions.a: fp2e_add2.o fp2e_sub2.o \
	fp2e_double2.o fp2e_triple2.o fp2e_neg2.o \
	fp2e_mul.o fp2e_mul_fpe.o fp2e_short_coeffred.o \
	fp2e_add.o fp2e_sub.o fp2e_parallel_coeffmul.o fp2e_mulxi.o\
	fp2e_double.o fp2e_triple.o fp2e_neg.o fp2e_conjugate.o \
	fpe_mul.o fp2e_square.o \
	consts.o
	rm -f asfunctions.a
	ar cr asfunctions.a $^

ASOBJECTS=fp2e_add2.o \
          fp2e_sub2.o \
          fp2e_double2.o \
          fp2e_triple2.o \
          fp2e_neg2.o \
          fp2e_mul.o \
          fp2e_mul_fpe.o \
          fp2e_short_coeffred.o \
          fp2e_add.o \
          fp2e_sub.o \
          fp2e_parallel_coeffmul.o \
          fp2e_mulxi.o \
          fp2e_double.o \
          fp2e_triple.o \
          fp2e_neg.o \
          fp2e_conjugate.o \
          fpe_mul.o \
          fp2e_square.o \
          consts.o
COBJECTS=linefunction.o \
         optate.o \
         fpe.o \
         fp2e.o \
         fp6e.o \
         fp12e.o \
         curvepoint_fp.o \
         twistpoint_fp2.o \
         final_expo.o \
         scalar.o \
         parameters.o \
         mul.o \
         mydouble.o
LIBOBJECTS=$(COBJECTS)
ifndef USE_C
	LIBOBJECTS=$(COBJECTS) $(ASOBJECTS)
endif
libdclxvi.a: $(LIBOBJECTS)
	rm -f $@
	ar cr $@ $^

.PHONY: clean

clean:
	-rm bilintest-check 
	-rm speedtest-check 
	-rm bilintest-c
	-rm speedtest-c
	-rm bilintest-as 
	-rm speedtest-as 
	-rm test_curvepoint_multiscalar-as
	-rm test_curvepoint_multiscalar-check
	-rm *.o
	-rm asfunctions.a
	-rm libdclxvi.a
