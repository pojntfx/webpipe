#define FUSE_USE_VERSION 31

#include <cuse_lowlevel.h>

typedef struct fuse_conn_info fuse_conn_info;
typedef struct fuse_file_info fuse_file_info;
typedef struct fuse_pollhandle fuse_pollhandle;

int wbcuse_start(int argc, char **argv);