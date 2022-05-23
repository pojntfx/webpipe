#define FUSE_USE_VERSION 31

#include "cuse.h"
#include <cuse_lowlevel.h>
#include <fuse.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>

void wbcuse_init(void *userdata, struct fuse_conn_info *conn) {
  printf("wbcuse_init\n");
  fflush(stdout);
}

void wbcuse_init_done(void *userdata) {
  printf("wbcuse_init_done\n");
  fflush(stdout);
}

void wbcuse_destroy(void *userdata) {
  printf("wbcuse_destroy\n");
  fflush(stdout);
}

// void wbcuse_open(fuse_req_t req, struct fuse_file_info *fi) {
//   printf("wbcuse_open\n");
//   fflush(stdout);
// }

void wbcuse_read(fuse_req_t req, size_t size, off_t off,
                 struct fuse_file_info *fi) {
  printf("wbcuse_read\n");
  fflush(stdout);
}

void wbcuse_write(fuse_req_t req, const char *buf, size_t size, off_t off,
                  struct fuse_file_info *fi) {
  printf("wbcuse_write\n");
  fflush(stdout);
}

void wbcuse_flush(fuse_req_t req, struct fuse_file_info *fi) {
  printf("wbcuse_flush\n");
  fflush(stdout);
}

void wbcuse_release(fuse_req_t req, struct fuse_file_info *fi) {
  printf("wbcuse_release\n");
  fflush(stdout);
}

void wbcuse_fsync(fuse_req_t req, int datasync, struct fuse_file_info *fi) {
  printf("wbcuse_fsync\n");
  fflush(stdout);
}

void wbcuse_ioctl(fuse_req_t req, int cmd, void *arg, struct fuse_file_info *fi,
                  unsigned int flags, const void *in_buf, size_t in_bufsz,
                  size_t out_bufsz) {
  printf("wbcuse_ioctl\n");
  fflush(stdout);
}

void wbcuse_poll(fuse_req_t req, struct fuse_file_info *fi,
                 struct fuse_pollhandle *ph) {
  printf("wbcuse_pollt\n");
  fflush(stdout);
}

int wbcuse_start(int argc, char **argv) {
  struct cuse_info ci;
  memset(&ci, 0, sizeof(ci));

  const char *dev_info_argv[] = {"DEVNAME=wbcuse"};
  ci.dev_major = 69;
  ci.dev_minor = 69;
  ci.dev_info_argc = 1;
  ci.dev_info_argv = dev_info_argv;
  ci.flags = CUSE_UNRESTRICTED_IOCTL;

  struct cuse_lowlevel_ops clop;
  memset(&clop, 0, sizeof(clop));
  clop.init = wbcuse_init;
  clop.init_done = wbcuse_init_done;
  clop.destroy = wbcuse_destroy;
  clop.open = wbcuse_open;
  clop.read = wbcuse_read;
  clop.write = wbcuse_write;
  clop.flush = wbcuse_flush;
  clop.release = wbcuse_release;
  clop.fsync = wbcuse_fsync;
  clop.ioctl = wbcuse_ioctl;
  clop.poll = wbcuse_poll;

  struct fuse_session *se;
  int multithreaded;
  int res;

  se = cuse_lowlevel_setup(argc, argv, &ci, &clop, &multithreaded, NULL);
  if (se == NULL) {
    return 1;
  }

  if (multithreaded) {
    res = fuse_session_loop_mt(se, 0);
  } else {
    res = fuse_session_loop(se);
  }

  cuse_lowlevel_teardown(se);

  if (res == -1) {
    return 1;
  }

  return 0;
}
