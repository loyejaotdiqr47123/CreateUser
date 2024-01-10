package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	USER_PRIV_GUEST uint32 = 0
	USER_PRIV_USER uint32 = 1
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procNetUserAdd = user32.NewProc("NetUserAdd")
	procNetLocalGroupAddMembers = user32.NewProc("NetLocalGroupAddMembers")
)

func createUser(username string, password string) error {
	userInfo := &USER_INFO_1{
		Name:       syscall.StringToUTF16Ptr(username),
		Password:   syscall.StringToUTF16Ptr(password),
		Priv:       USER_PRIV_USER,
		HomeDir:    nil,
		Comment:    nil,
		Flags:      0,
		ScriptPath: nil,
	}

	err := NetUserAdd(userInfo)
	if err != nil {
		fmt.Println("创建用户失败:", err)
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: CreateUser.exe username password")
		os.Exit(1)
	}

	username := os.Args[1]
	password := os.Args[2]

	err := createUser(username, password)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("创建用户成功.")
}

func NetUserAdd(userInfo *USER_INFO_1) error {
	r1, _, err := procNetUserAdd.Call(
		uintptr(0), // servername
		uintptr(1), // level
		uintptr(unsafe.Pointer(userInfo)),
		uintptr(0), // error code buffer
	)
	if r1 == 0 {
		return os.NewSyscallError("NetUserAdd", err)
	}
	return nil
}

type USER_INFO_1 struct {
	Name       *uint16
	Password   *uint16
	Priv       uint32
	HomeDir    *uint16
	Comment    *uint16
	Flags      uint32
	ScriptPath *uint16
}
