package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	netapi32DLL          = syscall.NewLazyDLL("netapi32.dll")
	netUserAddProc       = netapi32DLL.NewProc("NetUserAdd")
	netLocalGroupAddProc = netapi32DLL.NewProc("NetLocalGroupAddMembers")
	userPrivUser         = uint32(1)
	userFlagScript       = uint32(1)
	userFlagPasswordCantChange = uint32(0x40)
	userFlagDontExpirePasswd   = uint32(0x80)
)

type USER_INFO_1 struct {
	usri1_name      *uint16
	usri1_password  *uint16
	usri1_priv      uint32
	usri1_home_dir  *uint16
	usri1_comment   *uint16
	usri1_flags     uint32
	usri1_script_path *uint16
}

func NetUserAdd(servername *uint16, level uint32, buf *USER_INFO_1, parm_err *uint32) uint32 {
	ret, _, _ := netUserAddProc.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(parm_err)),
	)
	return uint32(ret)
}

func NetLocalGroupAddMembers(servername *uint16, groupname *uint16, level uint32, buf *USER_INFO_1, totalentries uint32) uint32 {
	ret, _, _ := netLocalGroupAddProc.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(unsafe.Pointer(groupname)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(totalentries),
	)
	return uint32(ret)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: CreateUser.exe username password")
		return
	}

	username := syscall.StringToUTF16Ptr(os.Args[1])
	password := syscall.StringToUTF16Ptr(os.Args[2])

	var ui USER_INFO_1
	ui.usri1_name = username
	ui.usri1_password = password
	ui.usri1_priv = userPrivUser
	ui.usri1_flags = userFlagScript | userFlagPasswordCantChange | userFlagDontExpirePasswd

	err := NetUserAdd(nil, 1, &ui, nil)
	if err != 0 {
		fmt.Printf("Error creating user: %d\n", err)
		return
	}

	groupname := syscall.StringToUTF16Ptr("Users")
	err = NetLocalGroupAddMembers(nil, groupname, 3, &ui, 1)
	if err != 0 {
		fmt.Printf("Error adding user to group: %d\n", err)
		return
	}

	fmt.Printf("User '%s' created and added to group 'Users'\n", os.Args[1])
}
