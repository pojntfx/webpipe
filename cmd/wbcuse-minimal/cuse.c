#define FUSE_USE_VERSION 31

#include <cuse_lowlevel.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>

void wbcuse_open(fuse_req_t req, struct fuse_file_info *fi) {
  printf("onopen");
}

void wbcuse_read(fuse_req_t req, size_t size, off_t off,
                 struct fuse_file_info *fi) {
  printf("onread");
}

void wbcuse_write(fuse_req_t req, const char *buf, size_t size, off_t off,
                  struct fuse_file_info *fi) {
  printf("onwrite");
}

void wbcuse_ioctl(fuse_req_t req, int cmd, void *arg, struct fuse_file_info *fi,
                  unsigned flags, const void *in_buf, size_t in_bufsz,
                  size_t out_bufsz) {
  printf("onioctl");
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

  clop.open = wbcuse_open;
  clop.read = wbcuse_read;
  clop.write = wbcuse_write;
  clop.ioctl = wbcuse_ioctl;

  return cuse_lowlevel_main(argc, argv, &ci, &clop, NULL);
}