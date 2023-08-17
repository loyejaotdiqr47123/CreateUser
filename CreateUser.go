package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")
	modnetapi32 = syscall.NewLazyDLL("netapi32.dll")

	procNetUserAdd                  = modnetapi32.NewProc("NetUserAdd")
	procNetLocalGroupAddMembers    = modnetapi32.NewProc("NetLocalGroupAddMembers")
	procNetLocalGroupGetMembers    = modnetapi32.NewProc("NetLocalGroupGetMembers")
	procNetApiBufferFree           = modnetapi32.NewProc("NetApiBufferFree")
)

const (
	ERROR_SUCCESS = 0

	USER_PRIV_USER = 1
	UF_SCRIPT      = 1
	UF_PASSWD_NOTREQD = 32
)

type USER_INFO_1 struct {
	Username  *uint16
	Password  *uint16
	PasswordAge uint32
	Priv      uint32
	HomeDir   *uint16
	Comment   *uint16
	Flags     uint32
	ScriptPath *uint16
}

type LOCALGROUP_MEMBERS_INFO_3 struct {
	lgrmi3_domainandname uintptr
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

func NetLocalGroupGetMembers(serverName *uint16, groupName *uint16, level uint32, buf **byte, prefMaxLen uint32, entriesRead, totalEntries *uint32, resumeHandle *uintptr) (netapiStatus uint32) {
	ret, _, _ := procNetLocalGroupGetMembers.Call(
		uintptr(unsafe.Pointer(serverName)),
		uintptr(unsafe.Pointer(groupName)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(prefMaxLen),
		uintptr(unsafe.Pointer(entriesRead)),
		uintptr(unsafe.Pointer(totalEntries)),
		uintptr(unsafe.Pointer(resumeHandle)),
	)
	netapiStatus = uint32(ret)
	return
}

func NetApiBufferFree(buf *byte) {
	procNetApiBufferFree.Call(uintptr(unsafe.Pointer(buf)))
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

	buf, err := syscall.Marshal(&ui1)
	if err != nil {
		fmt.Println("Error marshaling user info:", err)
		return
	}

	ret := NetUserAdd(nil, 1, &buf[0], nil)
	if ret != ERROR_SUCCESS {
		fmt.Println("Error creating user:", ret)
		return
	}

	groupName := syscall.StringToUTF16Ptr("Users")
	memberInfo := LOCALGROUP_MEMBERS_INFO_3{lgrmi3_domainandname: uintptr(unsafe.Pointer(username))}
	buf, err = syscall.Marshal(&memberInfo)
	if err != nil {
		fmt.Println("Error marshaling group member info:", err)
		return
	}

	ret = NetLocalGroupAddMembers(nil, groupName, 3, &buf[0], 1)
	if ret != ERROR_SUCCESS {
		fmt.Println("Error adding user to group:", ret)
		return
	}

	fmt.Println("User created and added to group successfully.")
}
