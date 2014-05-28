CC = gcc
CFLAGS = -Wall -O3

OBJS = board.o cycles.o mask.o moves.o queue.o

test: board_test cycles_test mask_test moves_test
	./board_test
	./cycles_test
	./mask_test
	./moves_test

board_test: board_test.c $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

cycles_test: cycles_test.c $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

mask_test: mask_test.c $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

moves_test: moves_test.c $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^

clean:
	rm *.o *_test
