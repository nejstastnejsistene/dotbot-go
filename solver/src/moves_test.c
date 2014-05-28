#include <stdio.h>

#include "board.h"
#include "cycles.h"
#include "mask.h"
#include "moves.h"

int TestPathOnRandomBoards() {
    int i;
	for (i = 0; i < 1000; i++) {
        Board board = {{0}};
        FillEmpty(board);
		Move move = ChooseMove(board, -1);
		Queue *path = ConstructPath(move);
		if (CYCLIC(move)) {
			if (path->size != Count(PATH(move))+1) {
				fprintf(stderr, "length of path should be numDots + 1:\n");
                PrintMask(PATH(move));
                return -1;
			}
			if (path->values[0] != path->values[path->size-1]) {
				fprintf(stderr, "path does not connect to itself\n");
                return -1;
			}
		} else if (path->size != Count(PATH(move))) {
			fprintf(stderr, "length of path should be the number of dots:\n");
            PrintMask(PATH(move));
            return -1;
		}
	}
    return 0;
}

int main() {
    init();
    int failed = 0;
    if (TestPathOnRandomBoards() < 0) {
        fprintf(stderr, "TestPathOnRandomBoards: FAILED\n");
        failed = 1;
    }
    return failed;
}
