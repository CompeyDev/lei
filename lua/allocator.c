#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>

typedef struct MemoryState MemoryState;

void* allocator(void* ud, void* ptr, size_t osize, size_t nsize) {
    MemoryState* mem = (MemoryState*)ud;

    // Convert void* to Go struct pointer using uintptr arithmetic
    // SAFETY: unsafe, must keep pointer alive in Go!
    int* usedMemory      = (int*)((uintptr_t)mem + 0);
    int* memoryLimit     = (int*)((uintptr_t)mem + sizeof(int));
    bool* ignoreLimit    = (bool*)((uintptr_t)mem + 2*sizeof(int));
    bool* limitReached   = (bool*)((uintptr_t)mem + 2*sizeof(int) + sizeof(bool));

    *limitReached = false;

    if (nsize == 0) {
        if (ptr != NULL) {
            free(ptr);
            *usedMemory -= (int)osize;
        }
        return NULL;
    }

    if (nsize > ((size_t)-1) >> 1) return NULL;

    int memDiff = ptr ? (int)nsize - (int)osize : (int)nsize;
    int newUsed = *usedMemory + memDiff;

    if (*memoryLimit > 0 && newUsed > *memoryLimit && !*ignoreLimit) {
        *limitReached = true;
        return NULL;
    }

    *usedMemory = newUsed;

    void* newPtr;
    if (ptr == NULL) {
        newPtr = malloc(nsize);
        if (!newPtr) abort();
    } else {
        newPtr = realloc(ptr, nsize);
        if (!newPtr) abort();
    }

    return newPtr;
}