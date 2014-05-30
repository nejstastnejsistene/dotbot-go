#include <stdio.h>
#include <stdlib.h>

#include "cycles.h"

void Cycles(Mask mask, Mask colorMask, Queue *cycles) {
    if (!MATCHES(colorMask, mask)) {
        fprintf(stderr, "mask is not contained within colorMask\n");
        exit(1);
    }
    int numDots = Count(mask);
    if (numDots < 4) {
        return;
    }
    int r0, c0, r1, c1;
    ConvexHull(mask, &r0, &c0, &r1, &c1);
    int numRows = r1 - r0 + 1;
    int numCols = c1 - c0 + 1;
    if (numRows < 2 || numCols < 2) {
        return;
    }

    Queue *seen = NewQueue();
    int seenSquare = 0;

    Mask cycle, result;
    int rows, cols, i, r, c, j, notSeen;
	for (rows = 3; rows <= numRows; rows++) {
		for (cols = 3; cols <= numCols; cols++) {
			if (numDots >= PERIMETER(rows, cols)-1) {
                for (i = 0; i < db[rows][cols]->size; i++) {
                    for (r = r0; r <= r1-rows+1; r++) {
                        for (c = c0; c <= c1-cols+1; c++) {
                            cycle = (Mask)db[rows][cols]->values[i] << INDEX(r, c);
                            if (MATCHES(mask, cycle)) {
                                result = colorMask | Encircled(cycle);
                                if (result == colorMask) {
                                    seenSquare = 1;
                                }
                                notSeen = 1;
                                for (j = 0; j < seen->size; j++) {
                                    if ((Mask)seen->values[j] == result) {
                                        notSeen = 0;
                                        break;
                                    }
                                }
                                if (notSeen) {
                                    Push(seen, (void*)result);
                                    Push(cycles, (void*)cycle);
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    FreeQueue(seen);

    if (!seenSquare) {
        Mask square = findSquare(mask, r0, c0, r1, c1);
        if (square != 0) {
            Push(cycles, (void*)square);
        }
    }
}

Mask findSquare(Mask mask, int r0, int c0, int r1, int c1) {
    if (!INBOUNDS(r0, c0) || !INBOUNDS(r1, c1)) {
        fprintf(stderr, "Out of bounds convex hull\n");
        exit(1);
    }
    Mask square;
    int r, c;
    for (r = r0; r < r1; r++) {
        for (c = c0; c < c1; c++) {
            square = Square << INDEX(r, c);
            if (MATCHES(mask, square)) {
                return square;
            }
        }
    }
    return 0;
}

void ConvexHull(Mask mask, int *r0, int *c0, int *r1, int *c1) {
    if (mask == 0) {
        *r0 = *c0 = *r1 = *c1 = 0;
        return;
    }
    *r0 = BoardSize;
    *c0 = BoardSize;
    *r1 = 0;
    *c1 = 0;
    int row, col;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
            if (CONTAINS(mask, row, col)) {
                if (row < *r0) {
                    *r0 = row;
                }
                if (row > *r1) {
                    *r1 = row;
                }
                if (col < *c0) {
                    *c0 = col;
                }
                if (col > *c1) {
                    *c1 = col;
                }
            }
        }
    }
    return;
}

Mask Encircled(Mask mask) {
    Mask h, v;
    int r, c;
    h = v = mask;
	for (r = 0; r < BoardSize; r++) {
		for (c = 0; c < BoardSize && !CONTAINS(h, r, c); c++) {
			h = ADD(h, r, c);
		}
		for (c = BoardSize - 1; c >= 0 && !CONTAINS(h, r, c); c--) {
            h = ADD(h, r, c);
		}
	}
	for (c = 0; c < BoardSize; c++) {
		for (r = 0; r < BoardSize && !CONTAINS(v, r, c); r++) {
            v = ADD(v, r, c);
		}
		for (r = BoardSize - 1; r >= 0 && !CONTAINS(v, r, c); r--) {
			v = ADD(v, r, c);
		}
	}
	return ~(h | v) & AllDots;
}

void buildCandidateCycles(Queue *cycles, Mask cycle, int col, int prevStart, int prevEnd, int rows, int cols) {
    Mask newCycle;
    int start, end, row;
    for (start = 0; start < rows-2; start++) {
        for (end = start + 2; end < rows; end++) {
            newCycle = cycle;
            newCycle = ADD(newCycle, start, col);
            newCycle = ADD(newCycle, end, col);
            if (col == 0 || col == cols-1) {
                for (row = start + 1; row < end; row++) {
                    newCycle = ADD(newCycle, row, col);
                }
            }
            if (col > 0) {
                for (row = start + 1; row <= prevStart; row++) {
                    newCycle = ADD(newCycle, row, col);
                }
                for (row = prevStart + 1; row <= start; row++) {
                    newCycle = ADD(newCycle, row, col-1);
                }
                for (row = prevEnd; row < end; row++) {
                    newCycle = ADD(newCycle, row, col);
                }
                for (row = end; row < prevEnd; row++) {
                    newCycle = ADD(newCycle, row, col-1);
                }
            }

            if (col+1 == cols) {
                Push(cycles, (void*)newCycle);
            } else {
                buildCandidateCycles(cycles, newCycle, col+1, start, end, rows, cols);
            }
        }
    }
}

int isValidCycle(Mask mask, int rows, int cols) {
    if (findSquare(mask, 0, 0, rows-1, cols-1)) {
        return 0;
    }
    int row, col, n, numDots;
    numDots = PERIMETER(rows, cols);
    for (row = 0; row < rows; row++) {
        for (col = 0; col < cols; col++) {
            if (CONTAINS(mask, row, col)) {
                n = CountNeighbors(mask, row, col);
                if (n == 4) {
                    numDots--;
                } else if (n < 2) {
                    return 0;
                }
            }
        }
    }
    if (Count(mask) != numDots) {
        return 0;
    }
    return 1;
}

void init() {
    Queue *candidates;
    Mask cycle;
    int rows, cols;
    for (rows = 3; rows <= BoardSize; rows++) {
        for (cols = 3; cols <= BoardSize; cols++) {
            db[rows][cols] = NewQueue();
            candidates = NewQueue();
            buildCandidateCycles(candidates, 0, 0, -1, -1, rows, cols);
            while (candidates->size) {
                cycle = (Mask)Pop(candidates);
                if (isValidCycle(cycle, rows, cols)) {
                    Push(db[rows][cols], (void*)cycle);
                }
            }
            FreeQueue(candidates);
        }
    }
}
