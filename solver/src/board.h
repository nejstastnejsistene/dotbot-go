#ifndef BOARD_H
#define BOARD_H

#define BoardSize 6

typedef enum {
	Empty,
	NotEmpty,
	Red,
	Yellow,
	Green,
	Blue,
	Purple,
} Color;

typedef Color Board[BoardSize][BoardSize];

Color RandomColor();
void FillEmpty(Board board);
void FillEmptyExcluding(Board board, Color exclude);
void Shrink(Board board, int row, int col);
void PrintBoard(Board board);

#endif
