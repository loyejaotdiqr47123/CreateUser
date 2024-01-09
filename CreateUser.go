package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/unixpickle/pe"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: gobuild <input.exe> <output.exe>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	exe, err := pe.NewExe(f)
	if err != nil {
		log.Fatal(err)
	}

	var sysCalls []string
	for _, imp := range exe.Imports {
		if imp.Name == "syscall" {
			sysCalls = append(sysCalls, imp.Hint)
		}
	}

	if len(sysCalls) == 0 {
		fmt.Println("No syscall functions found.")
		os.Exit(0)
	}

	fmt.Println("Found syscall functions:", sysCalls)

	// Replace syscall functions with obfuscated names
	for i, call := range sysCalls {
		exe.Imports[i].Hint = obfuscate(call)
	}

	f, err = os.Create(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = exe.Write(f)
	if err != nil {
		log.Fatal(err)
	}
}

func obfuscate(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		b.WriteByte(s[i] ^ 0x55)
	}
	return b.String()
}
