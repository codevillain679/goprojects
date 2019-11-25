// chapter5
package main

import (
	"fmt"
)

func main() {
	for i := 1; i <= 100; i++ {
		if i%3 == 0 {
			fmt.Print("Fizz")
			if i%5 == 0 {
				fmt.Println("Buzz")
			}
			if i%5 != 0 {
				fmt.Println()
			}
			continue
		}
		if i%5 == 0 {
			fmt.Println("Buzz")
			continue
		}
		fmt.Println(i)
	}
}
