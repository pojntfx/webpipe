#define FUSE_USE_VERSION 31

#include <cuse_lowlevel.h>
#include <stddef.h>
#include <stdio.h>
#include <string.h>

void _wbcuse_init_done(void *userdata) { wbcuse_init_done(userdata); }

void _wbcuse_destroy(void *userdata) { wbcuse_destroy(userdata); }

void _wbcuse_open(fuse_req_t req, struct fuse_file_info *fi) {
  wbcuse_open(req, fi);
}

void _wbcuse_read(fuse_req_t req, size_t size, off_t off,
                  struct fuse_file_info *fi) {
  wbcuse_read(req, size, off, fi);
}

void _wbcuse_write(fuse_req_t req, const char *buf, size_t size, off_t off,
                   struct fuse_file_info *fi) {
  wbcuse_write(req, buf, size, off, fi);
}

void _wbcuse_flush(fuse_req_t req, struct fuse_file_info *fi) {
  wbcuse_flush(req, fi);
}

void _wbcuse_release(fuse_req_t req, struct fuse_file_info *fi) {
  wbcuse_release(req, fi);
}

void _wbcuse_fsync(fuse_req_t req, int datasync, struct fuse_file_info *fi) {
  wbcuse_fsync(req, datasync, fi);
}

void _wbcuse_ioctl(fuse_req_t req, int cmd, void *arg,
                   struct fuse_file_info *fi, unsigned int flags,
                   const void *in_buf, size_t in_bufsz, size_t out_bufsz) {
  wbcuse_ioctl(req, cmd, arg, fi, flags, in_buf, in_bufsz, out_bufsz);
}

void _wbcuse_poll(fuse_req_t req, struct fuse_file_info *fi,
                  struct fuse_pollhandle *ph) {
  wbcuse_poll(req, fi, ph);
}

int wbcuse_start(int registry_id, int argc, char **argv) {
  struct cuse_info ci;
  memset(&ci, 0, sizeof(ci));

  const char *dev_info_argv[] = {"DEVNAME=wbcuse"};
  ci.dev_major = 69;
  ci.dev_minor = 69;
  ci.dev_info_argc = 1;
  ci.dev_info_argv = dev_info_argv;
  ci.flags = CUSE_UNRESTRICTED_IOCTL;

  auto void _wbcuse_init(void *userdata, struct fuse_conn_info *conn) {
    wbcuse_init(registry_id, userdata, conn);
  }

  struct cuse_lowlevel_ops clop;
  memset(&clop, 0, sizeof(clop));
  clop.init = _wbcuse_init;
  clop.init_done = _wbcuse_init_done;
  clop.destroy = _wbcuse_destroy;
  clop.open = _wbcuse_open;
  clop.read = _wbcuse_read;
  clop.write = _wbcuse_write;
  clop.flush = _wbcuse_flush;
  clop.release = _wbcuse_release;
  clop.fsync = _wbcuse_fsync;
  clop.ioctl = _wbcuse_ioctl;
  clop.poll = _wbcuse_poll;

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
