#include <stdio.h>

#include "cycles.h"
#include "test.h"

int main() {
    init();
    int failed = 0;
    if (TestAllBoard() < 0) {
        fprintf(stderr, "board_test: FAILED\n");
        failed = 1;
    }
    if (TestAllCycles() < 0) {
        fprintf(stderr, "cycles_test: FAILED\n");
        failed = 1;
    }
    if (TestAllMask() < 0) {
        fprintf(stderr, "mask_test: FAILED\n");
        failed = 1;
    }
    if (TestAllMoves() < 0) {
        fprintf(stderr, "moves_test: FAILED\n");
        failed = 1;
    }
    if (TestAllQueue() < 0) {
        fprintf(stderr, "queue_test: FAILED\n");
        failed = 1;
    }
    return failed;
}
