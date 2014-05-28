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
            weightedMove result;
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
    Queue *partitions, *cycles;
    Color color;
    int i, j;
    for (color = Red; color <= Purple; color++) {
        Mask colorMask = ColorMask(board, color);
        partitions = NewQueue();
        Partition(colorMask, partitions);
        for (i = 0; i < partitions->size; i++) {
            Mask partition = (Mask)partitions->values[i];
            cycles = NewQueue();
            Cycles(partition, colorMask, cycles);
            for (j = 0; j < cycles->size; j++) {
                Mask cycle = (Mask)cycles->values[j];
                Push(moves, (void*)ENCODE_MOVE(cycle, color, 1));
            }
            FreeQueue(cycles);
        }
        FreeQueue(partitions);
    }
}
