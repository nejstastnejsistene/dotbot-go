#include <stdio.h>

#include "mask.h"

Mask RandMask() {
    Mask mask;
    int row, col;
    for (row = 0; row < BoardSize; row++) {
		for (col = 0; col < BoardSize; col++) {
            if (rand() % 5 == 0) {
                mask |= DOTMASK(row, col);
            }
        }
    }
    return mask;
}

int TestInBounds() {
    int row, col, i;
    for (row = 0; row < BoardSize; row++) {
		for (col = 0; col < BoardSize; col++) {
			if (!INBOUNDS(row, col)) {
                fprintf(stderr, "Unexpected out of bounds: (%d, %d)\n", row, col);
                return -1;
			}
		}
	}
	for (i = 0; i < 1000; i++) {
		row = rand() + BoardSize;
		col = rand() + BoardSize;
		if (INBOUNDS(row, col) || INBOUNDS(-row, -col)) {
            fprintf(stderr, "Expected out of bounds: (%d, %d)\n", row, col);
            return -1;
		}
	}
    return 0;
}

int TestMatches() {
    Mask original, mask;
    int row, col, i, j, n;
	for (i = 0; i < 1000; i++) {
		original = mask = RandMask();
		n = rand() % NumDots;
		for (j = 0; j < n; j++) {
			row = rand() % BoardSize;
			col = rand() % BoardSize;
			mask = ADD(mask, row, col);
		}
		if (!MATCHES(mask, original)) {
			fprintf(stderr, "Masks do not match\n");
            return -1;
		}
	}
    return 0;
}

int TestAddRemove() {
    Mask mask;
    int row, col, i;
	for (i = 0; i < 1000; i++) {
		mask = 0;
		row = rand() % BoardSize;
		col = rand() % BoardSize;
		mask = ADD(mask, row, col);
		if (!CONTAINS(mask, row, col)) {
            fprintf(stderr, "Mask should contain: (%d, %d)\n", row, col);
            return -1;
		}
		mask = REMOVE(mask, row, col);
		if(CONTAINS(mask, row, col)) {
            fprintf(stderr, "Mask shouldn't contain: (%d, %d)\n", row, col);
            return -1;
		}
	}
    return 0;
}

int TestCount() {
    Mask mask;
    int row, col, count, i, j, n;
	for (i = 0; i < 1000; i++) {
		mask = 0;
		count = 0;
        int n;
        while (n < NumDots) {
            UNINDEX(n, row, col);
            mask = ADD(mask, row, col);
            count += rand() % 3 + 1;
        }
		if (Count(mask) != count) {
			fprintf(stderr, "Count() is incorrect\n");
            return -1;
		}
	}
    return 0;
}

int TestPartition() {
    Queue *partitions;
    Mask mask, p;
    int row, col, count, i;
	for (i = 0; i < 1000; i++) {
		partitions = NewQueue();
		mask = RandMask();
		Partition(mask, partitions);
		count = 0;
        while (partitions->size > 0) {
            p = (Mask)Pop(partitions);
			count += Count(p);
			for (row = 0; row < BoardSize; row++) {
				for (col = 0; col < BoardSize; col++) {
					if (CONTAINS(p, row, col) && !CONTAINS(mask, row, col)) {
						fprintf(stderr, "Partition contains dot not in original\n");
                        FreeQueue(partitions);
                        return -1;
					}
				}
			}
		}
		if (count != Count(mask)) {
			fprintf(stderr, "Total number of dots is incorrect\n");
            FreeQueue(partitions);
            return -1;
		}
        FreeQueue(partitions);
	}
    return 0;
}

int main() {
    int failed = 0;
    if (TestInBounds() < 0) {
        fprintf(stderr, "TestInBounds: FAILED\n");
        failed = 1;
    }
    if (TestMatches() < 0) {
        fprintf(stderr, "TestMatches: FAILED\n");
        failed = 1;
    }
    if (TestAddRemove() < 0) {
        fprintf(stderr, "TestAddRemove: FAILED\n");
        failed = 1;
    }
    if (TestCount() < 0) {
        fprintf(stderr, "TestCount: FAILED\n");
        failed = 1;
    }
    if (TestPartition() < 0) {
        fprintf(stderr, "TestPartition: FAILED\n");
        failed = 1;
    }
    return failed;
}


