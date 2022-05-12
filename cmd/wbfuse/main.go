package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"github.com/jacobsa/timeutil"
)

const (
	rootInode fuseops.InodeID = fuseops.RootInodeID + iota
	fileInode
)

type fs struct {
	fuseutil.NotImplementedFileSystem

	clock    timeutil.Clock
	filename string
	content  *strings.Reader
}

func updateTimestamps(clock timeutil.Clock, attr *fuseops.InodeAttributes) {
	now := clock.Now()

	attr.Atime = now
	attr.Mtime = now
	attr.Crtime = now
}

// General inode implementation
func (f *fs) GetInodeAttributes(
	ctx context.Context,
	op *fuseops.GetInodeAttributesOp) error {

	if op.Inode == rootInode {
		op.Attributes = fuseops.InodeAttributes{
			Nlink: 1,
			Mode:  os.ModePerm | os.ModeDir,
		}

		updateTimestamps(f.clock, &op.Attributes)

		return nil
	}

	if op.Inode == fileInode {
		op.Attributes = fuseops.InodeAttributes{
			Nlink: 1,
			Mode:  os.ModePerm,
			Size:  uint64(f.content.Len()),
		}

		updateTimestamps(f.clock, &op.Attributes)

		return nil
	}

	return fuse.ENOENT
}

func (f *fs) LookUpInode(
	ctx context.Context,
	op *fuseops.LookUpInodeOp) error {
	if op.Parent == rootInode {
		op.Entry.Child = fileInode
		op.Entry.Attributes = fuseops.InodeAttributes{
			Nlink: 1,
			Mode:  os.ModePerm,
			Size:  uint64(f.content.Len()),
		}

		return nil
	}

	return fuse.ENOENT
}

// Directories
func (f *fs) OpenDir(
	ctx context.Context,
	op *fuseops.OpenDirOp) error {
	if op.Inode == rootInode {
		return nil
	}

	return fuse.ENOENT
}

func (f *fs) ReadDir(
	ctx context.Context,
	op *fuseops.ReadDirOp) error {
	if op.Inode == rootInode {
		entries := []fuseutil.Dirent{
			{
				Offset: 1,
				Inode:  fileInode,
				Name:   f.filename,
				Type:   fuseutil.DT_File,
			},
		}

		if op.Offset > fuseops.DirOffset(len(entries)) {
			return fuse.EIO
		}

		entries = entries[op.Offset:]

		for _, e := range entries {
			n := fuseutil.WriteDirent(op.Dst[op.BytesRead:], e)
			if n == 0 {
				break
			}

			op.BytesRead += n
		}

		return nil
	}

	return fuse.ENOENT
}

func (f *fs) ReleaseDirHandle(
	ctx context.Context,
	op *fuseops.ReleaseDirHandleOp) error {
	return nil
}

// Files
func (f *fs) OpenFile(
	ctx context.Context,
	op *fuseops.OpenFileOp) error {
	if op.Inode == fileInode {
		return nil
	}

	return fuse.ENOENT
}

func (f *fs) ReadFile(
	ctx context.Context,
	op *fuseops.ReadFileOp) error {
	if op.Inode == fileInode {
		var err error
		op.BytesRead, err = f.content.ReadAt(op.Dst, op.Offset)

		if err == io.EOF {
			return nil
		}

		return err
	}

	return fuse.ENOENT
}

func (f *fs) FlushFile(
	ctx context.Context,
	op *fuseops.FlushFileOp) error {
	return nil
}

func (f *fs) ReleaseFileHandle(
	ctx context.Context,
	op *fuseops.ReleaseFileHandleOp) error {
	return nil
}

func main() {
	mountpoint := flag.String("mountpoint", "/tmp/wbfuse", "Where to mount the FUSE to")
	filename := flag.String("filename", "file.entangle", "Name of the file to mount")
	verbose := flag.Bool("verbose", false, "Whether to enable verbose logging")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	if err := os.MkdirAll(*mountpoint, os.ModePerm); err != nil {
		panic(err)
	}

	cfg := &fuse.MountConfig{}
	if *verbose {
		cfg.DebugLogger = log.New(os.Stderr, "fuse: ", 0)
	}

	mount, err := fuse.Mount(
		*mountpoint,
		fuseutil.NewFileSystemServer(
			&fs{
				clock:    timeutil.RealClock(),
				filename: *filename,
				content:  strings.NewReader("Hello, world!"),
			},
		),
		cfg,
	)
	if err != nil {
		panic(err)
	}

	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-s

		log.Println("Gracefully shutting down")

		go func() {
			<-s

			log.Println("Forcing shutdown")

			cancel()

			os.Exit(1)
		}()

		if err := fuse.Unmount(*mountpoint); err != nil {
			panic(err)
		}
	}()

	if err := mount.Join(ctx); err != nil {
		if err != context.Canceled {
			panic(err)
		}
	}
}
