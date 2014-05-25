#include <stdio.h>
#include <stdlib.h>

#include "mask.h"

int Count(Mask mask) {
    int count = 0;
    while (mask) {
        mask ^= (mask & -mask);
        count++;
    }
    return count;
}

int CountNeighbors(Mask mask, int row, int col) {
	int count = 0;
	if (CONTAINS(mask, row-1, col)) {
		count++;
	}
	if (CONTAINS(mask, row+1, col)) {
		count++;
	}
	if (CONTAINS(mask, row, col-1)) {
		count++;
	}
	if (CONTAINS(mask, row, col+1)) {
		count++;
	}
    return count;
}

void Partition(Mask mask, Queue *q) {
    int row, col;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
            if (CONTAINS(mask, row, col)) {
                Push(q, (void*)buildPartition(&mask, 0, row, col));
            }
        }
    }
}

Mask buildPartition(Mask *mask, Mask p, int row, int col) {
    *mask = REMOVE(*mask, row, col);
    p = ADD(p, row, col);

    if (CONTAINS(*mask, row-1, col)) {
        p = buildPartition(mask, p, row-1, col);
    }
    if (CONTAINS(*mask, row+1, col)) {
        p = buildPartition(mask, p, row+1, col);
    }
    if (CONTAINS(*mask, row, col-1)) {
        p = buildPartition(mask, p, row, col-1);
    }
    if (CONTAINS(*mask, row, col+1)) {
        p = buildPartition(mask, p, row, col+1);
    }

	return p;
}
