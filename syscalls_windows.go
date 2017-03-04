package main

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

func GetDiskSpace() (free uint64, total uint64) {
	h := syscall.MustLoadDLL("kernel32.dll")
	c := h.MustFindProc("GetDiskFreeSpaceExW")

	cwd, _ := os.Getwd()

	log.Println(cwd)

	_, _, err := c.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(cwd))),
		uintptr(unsafe.Pointer(&free)),
		uintptr(unsafe.Pointer(&total)),
		uintptr(0))

	if err != nil {
		log.Println(err)
	}

	return
}
