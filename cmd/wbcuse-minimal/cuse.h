#define FUSE_USE_VERSION 31

#include <cuse_lowlevel.h>

typedef void (*closure)();

typedef struct fuse_file_info fuse_file_info;

int wbcuse_start(int argc, char **argv);