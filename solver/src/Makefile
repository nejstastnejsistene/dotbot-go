CC = gcc
CFLAGS = -Wall

test: cycles_test mask_test
	./cycles_test
	./mask_test

cycles_test: cycles_test.c cycles.o mask.o queue.o
	$(CC) -o $@ $^

mask_test: mask_test.c mask.o queue.o
	$(CC) -o $@ $^

clean:
	rm *.o *_test
