#include <stdio.h>
#include <stdlib.h>

#include "mask.h"

Mask RandMask() {
    Mask mask = 0;
    int row, col;
    for (row = 0; row < BoardSize; row++) {
		for (col = 0; col < BoardSize; col++) {
            if (rand() % 5 == 0) {
                mask = ADD(mask, row, col);
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
    int row, col, count, i, n;
	for (i = 0; i < 1000; i++) {
		mask = 0;
		count = 0;
        n = 0;
        while (n < NumDots) {
            UNINDEX(n, row, col);
            mask = ADD(mask, row, col);
            n += rand() % 3 + 1;
            count++;
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

int TestDFS() {
	Mask mask = 1;
	Queue *paths = NewQueue();
    DFS(mask, paths);
    if (paths->size != 1) {
		fprintf(stderr, "Expecting one path for partition of size one\n");
        FreeQueue(paths);
        return -1;
    }
    if ((Mask)paths->values[0] != mask) {
        fprintf(stderr, "Expecting singleton mask's DFS to yield itself\n");
        FreeQueue(paths);
        return -1;
    }
    FreeQueue(paths);
    mask = 0x10c1;
    Queue *expected = NewQueue();
    Push(expected, (void*)1);
    Push(expected, (void*)(1 | (1 << BoardSize)));
    Push(expected, (void*)(1 | (3 << BoardSize)));
    Push(expected, (void*)(1 | (1 << BoardSize) | (1 << 2 * BoardSize)));
    Push(expected, (void*)(2 << BoardSize));
    Push(expected, (void*)(3 << BoardSize));
    Push(expected, (void*)((3 << BoardSize) | (1 << 2 * BoardSize)));
    Push(expected, (void*)(1 << 2 * BoardSize));
    Push(expected, (void*)((1 << BoardSize) | (1 << 2 * BoardSize)));
    Push(expected, (void*)(1 << BoardSize));
    paths = NewQueue();
    DFS(mask, paths);
    if (paths->size != expected->size) {
		fprintf(stderr, "Actual and expected paths for DFS don't match\n");
        FreeQueue(expected);
        FreeQueue(paths);
        return -1;
    }
    int i, j;
    for (i = 0; i < paths->size; i++) {
        int found = 0;
        for (j = 0; j < expected->size; j++) {
            if (paths->values[i] == expected->values[j]) {
                found = 1;
                break;
            }
        }
        if (!found) {
            fprintf(stderr, "Actual and expected paths for DFS don't match\n");
            FreeQueue(expected);
            FreeQueue(paths);
            return -1;
        }
    }
    FreeQueue(expected);
    FreeQueue(paths);
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
    if (TestDFS() < 0) {
        fprintf(stderr, "TestDFS: FAILED\n");
        failed = 1;
    }
    return failed;
}
