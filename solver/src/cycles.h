#ifndef CYCLES_H
#define CYCLES_H

#include "mask.h"

#define Square (3ULL | 3ULL << BoardSize)
#define PERIMETER(rows, cols) (2 * ((rows) + (cols)) - 4)

Queue *db[BoardSize+1][BoardSize+1];

Mask findSquare(Mask mask, int r0, int c0, int r1, int c1);
void ConvexHull(Mask mask, int *r0, int *c0, int *r1, int *c1);
Mask Encircled(Mask mask);
void buildCandidateCycles(Queue *cycles, Mask cycle, int col, int prevStart, int prevEnd, int rows, int cols);
int isValidCycle(Mask mask, int rows, int cols);

#endif
