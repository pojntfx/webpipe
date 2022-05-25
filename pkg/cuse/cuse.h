#define FUSE_USE_VERSION 31

#include <cuse_lowlevel.h>

typedef struct fuse_conn_info fuse_conn_info;
typedef struct fuse_file_info fuse_file_info;
typedef struct fuse_pollhandle fuse_pollhandle;

int wbcuse_start(void *device, unsigned int major, unsigned int minor,
                 char *name, int argc, char **argv);