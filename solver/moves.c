#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "board.h"
#include "cycles.h"
#include "mask.h"
#include "moves.h"

int MakeMove(Board board, Move move) {
    Mask dots = PATH(move);
    if (CYCLIC(move)) {
        dots |= Encircled(dots);
    }
    int row, col, score = 0;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
            if (CONTAINS(dots, row, col) || (CYCLIC(move) && board[col][row] == COLOR(move))) {
                Shrink(board, row, col);
                score++;
            }
        }
    }
    return score;
}

Move ChooseMove(Board board, int movesRemaining) {
	Queue *moves = NewQueue();
    Moves(board, moves);
	int maxDepth = movesRemaining;
	if (maxDepth <= 0 || maxDepth > MaxDepth) {
		maxDepth = MaxDepth;
	}
    weightedMove result = {0, 0, 0};
	chooseMove(board, moves, 0, 1, maxDepth, &result);
    FreeQueue(moves);
    return result.move;
}

void chooseMove(Board board, Queue *moves, int numEmpty, int depth, int maxDepth, weightedMove *chosen) {
    chosen->depth = depth;
    Board newBoard;
    int i, score, newNumEmpty, weight, deepest;
	for (i = 0; i < moves->size; i++) {
        Move move = (Move)moves->values[i];
		if (depth > 1 && Count(PATH(move)) == 1) {
			continue;
		}
        memcpy(newBoard, board, sizeof(newBoard));
		score = MakeMove(newBoard, move);
		newNumEmpty = numEmpty + score;
		weight = (float)score;
		deepest = depth;
		if (CYCLIC(move)) {
			weight *= CycleWeight;
		}
		if (numEmpty < Cutoff && depth < maxDepth) {
            Queue *newMoves = NewQueue();
            Moves(newBoard, newMoves);
            weightedMove result = {0, 0, 0};
			chooseMove(newBoard, newMoves, newNumEmpty, depth+1, maxDepth, &result);
            FreeQueue(newMoves);
			weight += Decay * result.weight;
			deepest = result.depth;
		}
		if (depth == 1) {
			weight /= (float)deepest;
		}
		if (weight > chosen->weight) {
            chosen->weight = weight;
            chosen->depth = deepest;
            chosen->move = move;
		}
	}
}

void Moves(Board board, Queue *moves) {
    Queue *partitions, *cycles, *paths;
    Color color;
    int i, j, k;
    for (color = Red; color <= Purple; color++) {
        Mask colorMask = ColorMask(board, color);
        partitions = NewQueue();
        Partition(colorMask, partitions);
        for (i = 0; i < partitions->size; i++) {
            Mask partition = (Mask)partitions->values[i];
            cycles = NewQueue();
            Cycles(partition, colorMask, cycles);
            if (cycles->size) {
                for (j = 0; j < cycles->size; j++) {
                    Mask cycle = (Mask)cycles->values[j];
                    Push(moves, ENCODE_MOVE(cycle, color, 1));
                }
            } else {
                paths = NewQueue();
                DFS(partition, paths);
                for (k = 0; k < paths->size; k++) {
                    Push(moves, ENCODE_MOVE((Mask)paths->values[k], color, 0));
                }
                FreeQueue(paths);
            }
            FreeQueue(cycles);
        }
        FreeQueue(partitions);
    }
}

Queue *ConstructPath(Move move) {
    int row, col;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
            if (CONTAINS(PATH(move), row, col)) {

                Queue *points = NewQueue();
                Item point = INDEX(row, col);
                Push(points, point);

                Mask mask = REMOVE(PATH(move), row, col);
				if (mask == 0) {
					return points;
				}
                // Success is indicated by a zero return value.
                if (constructPath(mask, points, point) == 0) {
                    if (CYCLIC(move)) {
                        Push(points, points->values[0]);
                    }
                    return points;
                } else {
                    FreeQueue(points);
                }
            }
        }
    }
    fprintf(stderr, "solver: unable to construct path:\n");
    PrintMask(PATH(move));
    exit(1);
}

int constructPath(Mask mask, Queue *points, int prev) {
    int prevRow, prevCol, row, col;
    UNINDEX(prev, prevRow, prevCol);
    if        (CONTAINS(mask, prevRow-1, prevCol)) {
        row = prevRow - 1;
        col = prevCol;
    } else if (CONTAINS(mask, prevRow+1, prevCol)) {
        row = prevRow + 1;
        col = prevCol;
    } else if (CONTAINS(mask, prevRow, prevCol-1)) {
        row = prevRow;
        col = prevCol - 1;
    } else if (CONTAINS(mask, prevRow, prevCol+1)) {
        row = prevRow;
        col = prevCol + 1;
    } else {
        return -1;
	}
    mask = REMOVE(mask, row, col);
    long long point = INDEX(row, col);
    Push(points, point);
	if (mask != 0) {
		int err = constructPath(mask, points, point);
        if (err != 0) {
            return err;
		}
	}
    return 0;
}
