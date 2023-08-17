package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	netapi32DLL      = syscall.NewLazyDLL("netapi32.dll")
	netUserAdd       = netapi32DLL.NewProc("NetUserAdd")
	usrerrInfo       = netapi32DLL.NewProc("NetUserGetInfo")
)

const (
	NERR_Success          = 0
	ERROR_ACCESS_DENIED   = 5
	ERROR_ALREADY_EXISTS  = 183
	ERROR_INVALID_PARAMETER = 87
)

type USER_INFO_1 struct {
	usri1_name       *uint16
	usri1_password   *uint16
	usri1_password_age uint32
	usri1_priv       uint32
	usri1_home_dir   *uint16
	usri1_comment    *uint16
	usri1_flags      uint32
	usri1_script_path *uint16
}

func NetUserAdd(servername *uint16, level uint32, buf *byte, parm_err *uint32) uint32 {
	ret, _, _ := netUserAdd.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(parm_err)),
	)
	return uint32(ret)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: CreateUser.exe <username> <password>")
		os.Exit(1)
	}

	username := syscall.StringToUTF16Ptr(os.Args[1])
	password := syscall.StringToUTF16Ptr(os.Args[2])

	var ui1 USER_INFO_1
	ui1.usri1_name = username
	ui1.usri1_password = password
	ui1.usri1_priv = 1
	ui1.usri1_flags = 0x200
	ui1.usri1_home_dir = nil
	ui1.usri1_comment = nil

	ret := NetUserAdd(nil, 1, (*byte)(unsafe.Pointer(&ui1)), nil)
	if ret == NERR_Success {
		fmt.Println("User created successfully.")
	} else if ret == ERROR_ACCESS_DENIED {
		fmt.Println("Access denied. Run the program with administrative privileges.")
	} else if ret == ERROR_ALREADY_EXISTS {
		fmt.Println("User already exists.")
	} else if ret == ERROR_INVALID_PARAMETER {
		fmt.Println("Invalid parameter.")
	} else {
		fmt.Printf("Error code: %d\n", ret)
	}
}
