package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	modnetapi32 = syscall.NewLazyDLL("netapi32.dll")

	procNetUserAdd             = modnetapi32.NewProc("NetUserAdd")
	procNetLocalGroupAddMembers = modnetapi32.NewProc("NetLocalGroupAddMembers")
)

const (
	ERROR_SUCCESS      = 0
	NERR_GroupNotFound = 2220

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

type LOCALGROUP_MEMBERS_INFO_3 struct {
	lgrmi3_domainandname *uint16
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

func NetLocalGroupAddMembers(serverName *uint16, groupName *uint16, level uint32, buf *byte, totalEntries uint32) (netapiStatus uint32) {
	ret, _, _ := procNetLocalGroupAddMembers.Call(
		uintptr(unsafe.Pointer(serverName)),
		uintptr(unsafe.Pointer(groupName)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(totalEntries),
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

	groupName := syscall.StringToUTF16Ptr("Users")
	memberInfo := LOCALGROUP_MEMBERS_INFO_3{lgrmi3_domainandname: username}
	buf := &memberInfo

	ret = NetLocalGroupAddMembers(nil, groupName, 3, (*byte)(unsafe.Pointer(buf)), 1)
	if ret != ERROR_SUCCESS && ret != NERR_GroupNotFound {
		fmt.Println("Error adding user to group:", ret)
		return
	}

	fmt.Println("User created and added to 'Users' group successfully.")
}
