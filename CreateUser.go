package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	modnetapi32 = syscall.NewLazyDLL("netapi32.dll")

	procNetUserAdd = modnetapi32.NewProc("NetUserAdd")
)

const (
	ERROR_SUCCESS = 0

	USER_PRIV_USER   = 1
	UF_SCRIPT        = 1
	UF_PASSWD_NOTREQD = 32
)

type USER_INFO_1 struct {
	Username   *uint16
	Password   *uint16
	PasswordAge uint32
	Priv       uint32
	HomeDir    *uint16
	Comment    *uint16
	Flags      uint32
	ScriptPath *uint16
}

func NetUserAdd(serverName *uint16, level uint32, buf *byte, parmErr *uint32) (netapiStatus uint32) {
	ret, _, _ := procNetUserAdd.Call(
		uintptr(unsafe.Pointer(serverName)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(parmErr)),
	)
	netapiStatus = uint32(ret)
	return
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: program.exe username password")
		return
	}

	username := syscall.StringToUTF16Ptr(os.Args[1])
	password := syscall.StringToUTF16Ptr(os.Args[2])

	ui1 := USER_INFO_1{
		Username: username,
		Password: password,
		Priv:     USER_PRIV_USER,
		Flags:    UF_SCRIPT | UF_PASSWD_NOTREQD,
	}

	ret := NetUserAdd(nil, 1, (*byte)(unsafe.Pointer(&ui1)), nil)
	if ret != ERROR_SUCCESS {
		fmt.Println("Error creating user:", ret)
		return
	}

	fmt.Println("User created successfully.")
}
