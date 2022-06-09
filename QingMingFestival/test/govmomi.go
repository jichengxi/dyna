package main

import (
	"fmt"
	"time"
)

//func vmConv(c *vim25.Client, mvm mo.VirtualMachine) {
//	vm := object.NewVirtualMachine(c, mvm.Reference())
//	fmt.Println(vm)
//}

func main() {
	a := time.Now().Format("2006-01-02-15-04-05")
	fmt.Println(a)
}
