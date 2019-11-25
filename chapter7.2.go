// chapter7.2
package main

import "fmt"

func main() {
	fmt.Println(fib(20))
}

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}
