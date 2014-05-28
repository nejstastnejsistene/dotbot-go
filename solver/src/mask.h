#ifndef MASK_H
#define MASK_H

#include "board.h"
#include "queue.h"

typedef long long unsigned int Mask;

#define NumDots (BoardSize * BoardSize)
#define AllDots ((1ULL << NumDots) - 1)

#define INBOUNDS(row, col)       (0 <= (row) && (row) < BoardSize && 0 <= (col) && (col) < BoardSize)
#define INDEX(row, col)          (BoardSize * (col) + (row))
#define UNINDEX(i, row, col)     {row = (i) % BoardSize; col = (i) / BoardSize;}
#define DOTMASK(row, col)        (1ULL << INDEX(row, col))
#define MATCHES(mask, pattern)   ((mask & pattern) == pattern)
#define CONTAINS(mask, row, col) (INBOUNDS(row, col) && MATCHES(mask, DOTMASK(row, col)))
#define ADD(mask, row, col)      (mask | DOTMASK(row, col))
#define REMOVE(mask, row, col)   (mask & ~DOTMASK(row, col))

int Count(Mask mask);
int CountNeighbors(Mask mask, int row, int col);
void Partition(Mask mask, Queue *q);
Mask buildPartition(Mask *mask, Mask p, int row, int col);
void DFS(Mask mask, Queue *paths);
void buildPaths(Mask mask, Queue *paths, int seen[NumDots][NumDots], int startIndex, int row, int col, Mask path);
void PrintMask(Mask mask);

#endif
