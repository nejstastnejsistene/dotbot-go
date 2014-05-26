CC = gcc
CFLAGS = -Wall

cycles_test: cycles_test.c cycles.o mask.o queue.o
	$(CC) -o $@ $^

mask_test: mask_test.c mask.o queue.o
	$(CC) -o $@ $^

clean:
	rm *.o *_test
