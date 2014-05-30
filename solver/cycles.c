#include <stdio.h>
#include <stdlib.h>

#include "cycles.h"

// Find all of the cycles in this mask. It will only yield
// unique masks, with respect to the effect that this mask
// would have to the board. It uses colorMask to help determine
// this. This means that if there are multiple squares, it will
// only return one of them, because all squares have the same
// effect on the board.
void Cycles(Mask mask, Mask colorMask, Queue *cycles) {
    if (!MATCHES(colorMask, mask)) {
        fprintf(stderr, "mask is not contained within colorMask\n");
        exit(1);
    }

    // Cycles have at least 4 dots.
    int numDots = Count(mask);
    if (numDots < 4) {
        return;
    }

    int r0, c0, r1, c1;
    ConvexHull(mask, &r0, &c0, &r1, &c1);

    // Cycles are at least 2x2.
    int numRows = r1 - r0 + 1;
    int numCols = c1 - c0 + 1;
    if (numRows < 2 || numCols < 2) {
        return;
    }

    Queue *seen = NewQueue();
    int seenSquare = 0;

    Mask cycle, result;
    int rows, cols, i, r, c, j, notSeen;
    // Compare all cycles from size 3x3 to numRows x numCols.
    for (rows = 3; rows <= numRows; rows++) {
        for (cols = 3; cols <= numCols; cols++) {
            // This prevents us from checking cycles that we don't
            // have enough dots to form. The perimeter is the number
            // of dots in any cycle, unless it crosses over itself.
            // In that case, there will be one less because a dot
            // is server as two corners. Example:
            //
            // X X X        This cycle is 5x5, so the perimeter
            // X   X        is 16. That center dot with four
            // X X X X X    neighbors rather than the typical two
            //     X   X    brings the actual number of dots to 15.
            //     X X X
            //
            if (numDots >= PERIMETER(rows, cols)-1) {
                for (i = 0; i < db[rows][cols]->size; i++) {
                    // Translate this pattern throughout the convex hull.
                    // As for the limits, imagine that the hull is the
                    // entire board, so r0 is 0 and r1 is 5. If rows is,
                    // for example, 4, we should check for the pattern
                    // starting at row 0, 1, and 2. 5 - 4 + 1 gives us
                    // the limit we need. Less than that and we can get
                    // false negatives, greater than that and we get false
                    // positives because the pattern will wrap around.
                    for (r = r0; r <= r1-rows+1; r++) {
                        for (c = c0; c <= c1-cols+1; c++) {
                            cycle = (Mask)db[rows][cols]->values[i] << INDEX(r, c);
                            if (MATCHES(mask, cycle)) {
                                result = colorMask | Encircled(cycle);
                                if (result == colorMask) {
                                    seenSquare = 1;
                                }
                                notSeen = 1;
                                for (j = 0; j < seen->size; j++) {
                                    if ((Mask)seen->values[j] == result) {
                                        notSeen = 0;
                                        break;
                                    }
                                }
                                if (notSeen) {
                                    Push(seen, (void*)result);
                                    Push(cycles, (void*)cycle);
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    FreeQueue(seen);

    // If there is a square that is not equivalent to any of the
    // previously seen cycles, yield the first one found.
    if (!seenSquare) {
        Mask square = findSquare(mask, r0, c0, r1, c1);
        if (square != 0) {
            Push(cycles, (void*)square);
        }
    }
}

// Returns the first square in this mask. Returns 0 if there is none.
Mask findSquare(Mask mask, int r0, int c0, int r1, int c1) {
    if (!INBOUNDS(r0, c0) || !INBOUNDS(r1, c1)) {
        fprintf(stderr, "Out of bounds convex hull\n");
        exit(1);
    }
    Mask square;
    int r, c;
    for (r = r0; r < r1; r++) {
        for (c = c0; c < c1; c++) {
            square = Square << INDEX(r, c);
            if (MATCHES(mask, square)) {
                return square;
            }
        }
    }
    return 0;
}

// Return the minimum and maximum rows and columns. In other words,
// return the convex hull, where (r0, c0), (r1, c1) are the
// top left and bottom right coordinates. Convex hull is a
// fancy math term to indicate the smallest convex set that
// contains all of a set of points.
// https://en.wikipedia.org/wiki/Convex_hull
void ConvexHull(Mask mask, int *r0, int *c0, int *r1, int *c1) {
    if (mask == 0) {
        *r0 = *c0 = *r1 = *c1 = 0;
        return;
    }
    *r0 = BoardSize;
    *c0 = BoardSize;
    *r1 = 0;
    *c1 = 0;
    int row, col;
    for (row = 0; row < BoardSize; row++) {
        for (col = 0; col < BoardSize; col++) {
            if (CONTAINS(mask, row, col)) {
                if (row < *r0) {
                    *r0 = row;
                }
                if (row > *r1) {
                    *r1 = row;
                }
                if (col < *c0) {
                    *c0 = col;
                }
                if (col > *c1) {
                    *c1 = col;
                }
            }
        }
    }
    return;
}

// Returns a Mask of any dots encircled by this cyclic Mask.
// Using this on masks that are cyclic or contain squares will have
// undefined/meaningless return values.
Mask Encircled(Mask mask) {
    Mask h, v;
    int r, c;
    h = v = mask;

    // This works by filling in the dots from each of the four directions,
    // outwards-in. Any dots left unfilled are encircled. This is done
    // with separately for the horizontal and vertical directions so that
    // they don't interfere with each other. This approach  would not work
    // for boards larger that 6x6 because it would become possible to create
    // concave cyclic paths.
    for (r = 0; r < BoardSize; r++) {
        for (c = 0; c < BoardSize && !CONTAINS(h, r, c); c++) {
            h = ADD(h, r, c);
        }
        for (c = BoardSize - 1; c >= 0 && !CONTAINS(h, r, c); c--) {
            h = ADD(h, r, c);
        }
    }
    for (c = 0; c < BoardSize; c++) {
        for (r = 0; r < BoardSize && !CONTAINS(v, r, c); r++) {
            v = ADD(v, r, c);
        }
        for (r = BoardSize - 1; r >= 0 && !CONTAINS(v, r, c); r--) {
            v = ADD(v, r, c);
        }
    }
    return ~(h | v) & AllDots;
}

// Build cycles from left to right, column by column. It does that
// by recognizing that all convex cycles without squares can be
// represented by the top and bottom dot for each column, or in
// other words, the start and end row for each column. For example,
//
// X X X X X X      X X X X X X
//                  X         X
// X X          =>  X X X     X
//     X                X X   X
//       X X X            X X X
//
// This algorithm recursively selects these start and end points
// and fills in the appropriate dots. This relatively small set
// of potential cycles is small enough that it can be filtered
// in a reasonable amount of time.
void buildCandidateCycles(Queue *cycles, Mask cycle, int col, int prevStart, int prevEnd, int rows, int cols) {
    Mask newCycle;
    int start, end, row;

    // Go through all the possible pairs of starts and ends. The starts
    // begin at 0 and can go until 2 less than the maximum size, go give
    // room for the end. The ends go from 2 more than the start until
    // the maximum size. The significance of this buffer of 2, is that it
    // is the smallest area that you can fit a corner into a cycle without
    // folding upon itself and creating a square.
    for (start = 0; start < rows-2; start++) {
        for (end = start + 2; end < rows; end++) {

            // Make a copy of the cycle and add the start and end points.
            newCycle = cycle;
            newCycle = ADD(newCycle, start, col);
            newCycle = ADD(newCycle, end, col);

            // For the first and last columsn, connect the dots
            // between the start and end rows.
            if (col == 0 || col == cols-1) {
                for (row = start + 1; row < end; row++) {
                    newCycle = ADD(newCycle, row, col);
                }
            }
            // This forms corners between the current and previous
            // start and end. For whichever of the starts is highest,
            // dots are placed below it in its column up until the other
            // start. For the ends, it is the same but from bottom to
            // top instead.
            //
            // Examples:
            //
            // X        X         X       X
            //   X      X X     X       X X
            //      =>              =>
            // X        X X       X     X X
            //   X        X     X       X
            //
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

            // Yield the generated cycle if on the last column, otherwise
            // recur to the next column.
            if (col+1 == cols) {
                Push(cycles, (void*)newCycle);
            } else {
                buildCandidateCycles(cycles, newCycle, col+1, start, end, rows, cols);
            }
        }
    }
}

// Check if a cycle fits our criterea for a unique cycle with no squares.
int isValidCycle(Mask mask, int rows, int cols) {
    if (findSquare(mask, 0, 0, rows-1, cols-1)) {
        return 0;
    }
    // Keep track of the number of dots it should have.
    int row, col, n, numDots;
    numDots = PERIMETER(rows, cols);
    for (row = 0; row < rows; row++) {
        for (col = 0; col < cols; col++) {
            if (CONTAINS(mask, row, col)) {
                n = CountNeighbors(mask, row, col);
                if (n == 4) {
                    // Four neighbors means the cycle
                    // crossed over itself, and needs one
                    // less corner dot.
                    numDots--;
                } else if (n < 2) {
                    // Zero neighbors means it's in isolation.
                    // One neighbor means it's an endpoint and
                    // is not necessary.
                    return 0;
                }
            }
        }
    }
    // Make sure it has the expected number of dots.
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
            db[rows][cols] = NewQueue();
            // Generate the potential cycles.
            candidates = NewQueue();
            buildCandidateCycles(candidates, 0, 0, -1, -1, rows, cols);
            while (candidates->size) {
                cycle = (Mask)Pop(candidates);
                // Filter out the invalid candidates.
                if (isValidCycle(cycle, rows, cols)) {
                    Push(db[rows][cols], (void*)cycle);
                }
            }
            FreeQueue(candidates);
        }
    }
}
