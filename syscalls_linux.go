package main

import (
	"log"
	"os"
	"syscall"
)

func GetDiskSpace() (free uint64, total uint64) {
	cwd, _ := os.Getwd()

	var fs syscall.Statfs_t
	err := syscall.Statfs(cwd, &fs)

	if err != nil {
		log.Println(err)
	} else {
		free = fs.Bfree * uint64(fs.Bsize)
		total = fs.Blocks * uint64(fs.Bsize)
	}

	return
}
