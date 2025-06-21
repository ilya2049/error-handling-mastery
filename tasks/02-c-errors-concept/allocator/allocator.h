#include <errno.h>

extern int errno;

#define ADMIN 777
#define MIN_MEMORY_BLOCK 1024

void *allocate(int user_id, size_t size)
{
    if (user_id != ADMIN) {
        errno = EPERM;
        
        return NULL;
    }

    if (size < MIN_MEMORY_BLOCK) {
        errno = EDOM;
        
        return NULL;
    }

    void *p = malloc(size);
    if (p == NULL) {
        errno = ENOMEM;

        return NULL;
    }

    return p;
}
