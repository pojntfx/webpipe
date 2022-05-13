package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
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

type readWriterAt interface {
	io.ReaderAt
	io.WriterAt
	Size() (int64, error)
}

type fileWithSize struct {
	*os.File
}

func (f *fileWithSize) Size() (int64, error) {
	stat, err := f.Stat()
	if err != nil {
		return -1, err
	}

	return stat.Size(), nil
}

type fs struct {
	fuseutil.NotImplementedFileSystem

	clock    timeutil.Clock
	filename string
	backend  readWriterAt
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
		size, err := f.backend.Size()
		if err != nil {
			return err
		}

		op.Attributes = fuseops.InodeAttributes{
			Nlink: 1,
			Mode:  os.ModePerm,
			Size:  uint64(size),
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
		size, err := f.backend.Size()
		if err != nil {
			return err
		}

		op.Entry.Child = fileInode
		op.Entry.Attributes = fuseops.InodeAttributes{
			Nlink: 1,
			Mode:  os.ModePerm,
			Size:  uint64(size),
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
		op.BytesRead, err = f.backend.ReadAt(op.Dst, op.Offset)

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

func (f *fs) WriteFile(
	ctx context.Context,
	op *fuseops.WriteFileOp) error {
	if op.Inode == fileInode {
		_, err := f.backend.WriteAt(op.Data, op.Offset)

		if err == io.EOF {
			return nil
		}

		return err
	}

	return fuse.ENOENT
}

func (f *fs) SetInodeAttributes(
	ctx context.Context,
	op *fuseops.SetInodeAttributesOp) error {
	return nil
}

func main() {
	mountpointFlag := flag.String("mountpoint", "/tmp/wbfuse-mountpoint", "Where to mount the FUSE to")
	filenameFlag := flag.String("filename", "file.entangled", "Name of the file to mount")
	backendFlag := flag.String("backend", "/tmp/wbfuse-backend.entangled", "Name of the file to use as the backend")
	verboseFlag := flag.Bool("verbose", false, "Whether to enable verbose logging")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	backend, err := os.OpenFile(*backendFlag, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(*mountpointFlag, os.ModePerm); err != nil {
		panic(err)
	}

	cfg := &fuse.MountConfig{}
	if *verboseFlag {
		cfg.DebugLogger = log.New(os.Stderr, "fuse: ", 0)
	}

	mount, err := fuse.Mount(
		*mountpointFlag,
		fuseutil.NewFileSystemServer(
			&fs{
				clock:    timeutil.RealClock(),
				filename: *filenameFlag,
				backend:  &fileWithSize{backend},
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

		if err := backend.Close(); err != nil {
			panic(err)
		}

		if err := fuse.Unmount(*mountpointFlag); err != nil {
			panic(err)
		}
	}()

	if err := mount.Join(ctx); err != nil {
		if err != context.Canceled {
			panic(err)
		}
	}
}
