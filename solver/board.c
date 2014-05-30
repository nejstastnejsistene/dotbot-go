#include <stdio.h>
#include <stdlib.h>

#include "board.h"
#include "mask.h"

Color RandomColor() {
    return rand() % (Purple - Red + 1) + Red;
}

void FillEmpty(Board board) {
    FillEmptyExcluding(board, Empty);
}

void FillEmptyExcluding(Board board, Color exclude) {
    Color color;
    int row, col;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
            if (board[col][row] == Empty) {
                do {
                    color = RandomColor();
                } while (color == exclude);
                board[col][row] = color;
            }
        }
    }
}

void Shrink(Board board, int row, int col) {
    while (row > 0) {
        board[col][row] = board[col][row-1];
        row--;
    }
    board[col][0] = Empty;
}

const Color colorCodes[] = {0, 0, 31, 33, 32, 36, 35};

void PrintBoard(Board board) {
    int row, col;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
            Color color = board[col][row];
            if (color == Empty) {
                printf("  ");
            } else if (color > 0 && color <= Purple) {
                printf(" \x1b[%dm\xe2\x97\x8f\x1b[0m", colorCodes[color]);
            }
        }
        printf("\n");
    }
    printf("\n");
}
