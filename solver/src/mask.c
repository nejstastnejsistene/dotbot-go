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

void DFS(Mask mask, Queue *paths) {
    // seen[startPoint][endPoints]
    // This uniquely identifies a path through a non cyclic mask.
    int seen[NumDots][NumDots] = {{0}};
    int row, col;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
			if (CONTAINS(mask, row, col)) {
                int start = INDEX(row, col);
                if (CountNeighbors(mask, row, col) < 2) {
                    buildPaths(mask, paths, seen, start, row, col, 0);
                } else {
                    Push(paths, (void*)DOTMASK(row, col));
                    seen[start][start] = 1;
                }
            }
        }
    }
}

void buildPaths(Mask mask, Queue *paths, int seen[NumDots][NumDots], int startIndex, int row, int col, Mask path) {
	mask = REMOVE(mask, row, col);
    path = ADD(path, row, col);

    int currentIndex = INDEX(row, col);
	if (!seen[startIndex][currentIndex]) {
		seen[startIndex][currentIndex] = 1;
		seen[currentIndex][startIndex] = 1;
        Push(paths, (void*)path);
	}

    if (CONTAINS(mask, row-1, col)) {
        buildPaths(mask, paths, seen, startIndex, row-1, col, path);
    }
    if (CONTAINS(mask, row+1, col)) {
        buildPaths(mask, paths, seen, startIndex, row+1, col, path);
    }
    if (CONTAINS(mask, row, col-1)) {
        buildPaths(mask, paths, seen, startIndex, row, col-1, path);
    }
    if (CONTAINS(mask, row, col+1)) {
        buildPaths(mask, paths, seen, startIndex, row, col+1, path);
    }
}

void PrintMask(Mask mask) {
    int row, col;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
            if (CONTAINS(mask, row, col)) {
                printf(" *");
            } else {
                printf("  ");
            }
        }
        printf("\n");
    }
    printf("\n");
}
