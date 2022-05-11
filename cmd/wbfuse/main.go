package main

import (
	"context"
	"flag"
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
	fileNode
)

type fs struct {
	fuseutil.NotImplementedFileSystem
	clock timeutil.Clock
}

func (f *fs) GetInodeAttributes(
	ctx context.Context,
	op *fuseops.GetInodeAttributesOp) error {
	now := f.clock.Now()

	if op.Inode == rootInode {
		op.Attributes = fuseops.InodeAttributes{
			Nlink:  1,
			Mode:   os.ModePerm | os.ModeDir,
			Atime:  now,
			Mtime:  now,
			Crtime: now,
		}

		return nil
	}

	if op.Inode == fileNode {
		op.Attributes = fuseops.InodeAttributes{
			Nlink:  1,
			Mode:   os.ModePerm,
			Size:   0,
			Atime:  now,
			Mtime:  now,
			Crtime: now,
		}

		return nil
	}

	return fuse.ENOENT
}

func main() {
	mountpoint := flag.String("mountpoint", "/tmp/wbfuse", "Where to mount the FUSE to")
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
				clock: timeutil.RealClock(),
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
