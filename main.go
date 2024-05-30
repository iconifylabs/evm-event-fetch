package main

import "fmt"

func main() {
	fmt.Println("Call Message Events")
	detailed_logs := false

	call_message(46048251, detailed_logs)

	fmt.Println()

	fmt.Println("Call Executed Events")
	call_executed(46048251, detailed_logs)

}
