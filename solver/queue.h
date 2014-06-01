#ifndef QUEUE_H
#define QUEUE_H

#define QUEUE_INITIAL_CAPACITY 8

#define Item long long unsigned int

typedef struct Queue {
    int size;
    int capacity;
    Item *values;
} Queue;

Queue *NewQueue();
void FreeQueue(Queue *q);
void Push(Queue *q, Item value);
Item Pop(Queue *q);

#endif
