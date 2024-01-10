package main

import (
   "fmt"
   "io/ioutil"
   "log"
   "os"
   "regexp"
)

func main() {
   if len(os.Args) != 2 {
   	fmt.Println("Usage: go run replace.go <path-to-exe>")
   	os.Exit(1)
   }

   exePath := os.Args[1]
   data, err := ioutil.ReadFile(exePath)
   if err != nil {
   	log.Fatalf("Error reading file: %v", err)
   }

   // Replace syscall with syscawsl in the binary code
   re := regexp.MustCompile(`(?m)(\x5b\x69\x66\x73\x63\x61\x6c\x63\x68\x20\x28\x24\x73\x79\x73\x74\x65\x6d\x29\x20\x3b)`)
   newData := re.ReplaceAll(data, []byte("(\x5b\x69\x66\x73\x63\x61\x77\x73\x6c\x20\x28\x24\x73\x79\x73\x74\x65\x6d\x29\x20\x3b)"))

   err = ioutil.WriteFile(exePath, newData, 0644)
   if err != nil {
   	log.Fatalf("Error writing file: %v", err)
   }

   fmt.Printf("Replaced syscall with syscawsl in %s\n", exePath)
}
