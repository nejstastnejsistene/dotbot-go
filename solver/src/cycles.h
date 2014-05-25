#ifndef CYCLES_H
#define CYCLES_H

#include "mask.h"

#define PERIMETER(rows, cols) (2 * ((rows) + (cols)) - 4)

Queue *db[BoardSize+1][BoardSize+1];

Mask findSquare(Mask mask, int r0, int c0, int r1, int c1);
void buildCandidateCycles(Queue *cycles, Mask cycle, int col, int prevStart, int prevEnd, int rows, int cols);
int isValidCycle(Mask mask, int rows, int cols);

#endif
