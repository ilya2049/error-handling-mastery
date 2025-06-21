#include <errno.h>

typedef int errno_t;

#define ADMIN 777
#define MIN_MEMORY_BLOCK 1024


errno_t allocate(int user_id, size_t size, void **mem)
{
if (user_id != ADMIN) {
        *mem = NULL;

        return EPERM;
    }

    if (size < MIN_MEMORY_BLOCK) {
        *mem = NULL;
        
        return EDOM;
    }

    void *p = malloc(size);
    if (p == NULL) {
        *mem = NULL;

        return ENOMEM;
    }

    *mem = p;

    return 0;
}
