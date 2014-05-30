#ifndef QUEUE_H
#define QUEUE_H

#define QUEUE_INITIAL_CAPACITY 8

typedef struct Queue {
    int size;
    int capacity;
    void **values;
} Queue;

Queue *NewQueue();
void FreeQueue(Queue *q);
void Push(Queue *q, void *value);
void *Pop(Queue *q);

#endif
