#include <stdio.h>

#include "cycles.h"
#include "test.h"

int TestCycles() {
    Mask mask;
    Queue *cycles = NewQueue();

    mask = 0x71441f05f;
	Cycles(mask, mask, cycles);
    if (cycles->size != 0) {
        fprintf(stderr, "Unexpected cycle\n");
        FreeQueue(cycles);
        return -1;
    }

    mask = 0xc8dc8dc8d;
	Cycles(mask, mask, cycles);
    if (cycles->size != 0) {
        fprintf(stderr, "Unexpected cycle\n");
        FreeQueue(cycles);
        return -1;
    }

    mask = 0x76797b57f;
	Cycles(mask, mask, cycles);
    if (cycles->size != 0) {
        fprintf(stderr, "Unexpected cycle\n");
        FreeQueue(cycles);
        return -1;
    }

    mask = 0x61c51f147;
	Cycles(mask, mask, cycles);
    int seen[4] = {0};
    if (cycles->size != 4) {
	    fprintf(stderr, "Actual and expected cycles don't match: expected 4 cycles, not %d\n", cycles->size);
        return -1;
    }
    int index, i;
    for (i = 0; i < cycles->size; i++) {
        switch ((Mask)cycles->values[i]) {
        case 0x7147:
            index = 0;
            break;
        case 0x1c51c000:
            index = 1;
            break;
        case 0x7147|0x1c51c000:
            index = 2;
            break;
        case Square << INDEX(3, 4):
            index = 3;
            break;
        default:
            fprintf(stderr, "Actual and expected cycles don't match: unexpected cycle\n");
            return -1;
        }
    }
    if (seen[index]) {
	    fprintf(stderr, "Actual and expected cycles don't match: duplicate cycles\n");
        return -1;
    } else {
        seen[index] = 1;
    }

    FreeQueue(cycles);
    return 0;
}

int TestAllCycles() {
    int failed = 0;
    if (TestCycles() < 0) {
        fprintf(stderr, "TestCycles: FAILED\n");
        failed = 1;
    }
    return failed;
}
