#include <stdio.h>
#include <stdlib.h>

#include "queue.h"

Queue *NewQueue() {
    Queue *q = malloc(sizeof(Queue));
    if (q < 0) {
        perror("NewQueue");
        exit(1);
    }
    q->size = 0;
    q->capacity = QUEUE_INITIAL_CAPACITY;
    q->values = malloc(sizeof(Item)*q->capacity);
    if (q->values < 0) {
        perror("NewQueue");
        exit(1);
    }
    return q;
}

void FreeQueue(Queue *q) {
    free(q->values);
    free(q);
}

void Push(Queue *q, Item value) {
    if (q->size == q->capacity) {
        q->capacity *= 2;
        q->values = realloc(q->values, sizeof(Item)*q->capacity);
        if (q->values < 0) {
            perror("Push");
            exit(1);
        }
    }
    q->values[q->size++] = value;
}

Item Pop(Queue *q) {
    if (q->size == 0) {
        fprintf(stderr, "queue: can't pop from empty queue\n");
        exit(1);
    }
    return q->values[--q->size];
}
