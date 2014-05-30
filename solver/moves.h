#ifndef MOVES_H
#define MOVES_H

#include "board.h"
#include "mask.h"
#include "queue.h"

#define MaxDepth    3
#define Cutoff      (NumDots / 2)
#define Decay       0.5
#define CycleWeight (1 / Decay)

typedef Mask Move;

#define CYCLIC_SHIFT                     NumDots
#define COLOR_SHIFT                      (CYCLIC_SHIFT + 1)
#define ENCODE_CYCLIC(cyclic)            (((Mask)(cyclic)) << CYCLIC_SHIFT)
#define ENCODE_COLOR(color)              (((Mask)(color)) << COLOR_SHIFT)
#define ENCODE_MOVE(path, color, cyclic) ((path) | ENCODE_COLOR(color) | ENCODE_CYCLIC(cyclic))
#define PATH(move)                       ((move) & AllDots)
#define COLOR(move)                      ((move) >> COLOR_SHIFT)
#define CYCLIC(move)                     (((move) >> CYCLIC_SHIFT) & 1)

typedef struct {
	float weight;
	int depth;
    Move move;
} weightedMove;

int MakeMove(Board board, Move move);
Move ChooseMove(Board board, int movesRemaining);
void chooseMove(Board board, Queue *moves, int numEmpty, int depth, int maxDepth, weightedMove *chosen);
void Moves(Board board, Queue *moves);
Queue *ConstructPath(Move move);
int constructPath(Mask mask, Queue *points, int prev);

#endif
