CC = gcc

mask_test: mask_test.c board.o cycles.o mask.o moves.o queue.o
	$(CC) -o $@ $^

clean:
	rm *.o *_test
