#ifndef CYCLES_H
#define CYCLES_H

#include "mask.h"

#define Square (3ULL | 3ULL << BoardSize)

// Returns the perimeter of a rows x cols cycle.
#define PERIMETER(rows, cols) (2 * ((rows) + (cols)) - 4)

// A database of all possible non-square cycles. The slice of
// cycles at db[rows][cols] is the list of all cycles of size
// rows x cols that doesn't contain a square. The logic behind
// this is that any cycle that has a square in it doesn't add
// anything to that same cycle with the square removed. E.g.
//
// X X X X                                           X X X
// X   X X will have the same affect on the board as X   X
// X X X X                                           X X X
//
// This is because any cycle will remove all dots of that color,
// so the only real differentiating factor between cycles is how
// many and which dots they encircle. This database only has data
// for sizes 3x3 to 6x6, because anything less than that either
// contains a square or no cycles at all. It is populated in init().
Queue *db[BoardSize+1][BoardSize+1];

void Cycles(Mask mask, Mask colorMask, Queue *cycles);
Mask findSquare(Mask mask, int r0, int c0, int r1, int c1);
void ConvexHull(Mask mask, int *r0, int *c0, int *r1, int *c1);
Mask Encircled(Mask mask);
void buildCandidateCycles(Queue *cycles, Mask cycle, int col, int prevStart, int prevEnd, int rows, int cols);
int isValidCycle(Mask mask, int rows, int cols);
void init();

#endif
