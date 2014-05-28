CC = gcc
CFLAGS = -Wall -O3

OBJS = board.o cycles.o mask.o moves.o queue.o
TEST_OBJS = board_test.o cycles_test.o mask_test.o moves_test.o queue_test.o

test: test.c $(TEST_OBJS) $(OBJS)
	$(CC) $(CFLAGS) -o $@ $^
	./$@

clean:
	rm *.o *test
