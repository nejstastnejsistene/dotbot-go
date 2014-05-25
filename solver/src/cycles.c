#include "cycles.h"

Mask findSquare(Mask mask, int r0, int c0, int r1, int c1) {
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
            candidates = NewQueue();
            db[rows][cols] = NewQueue();
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
