// chapter7
package main

import "fmt"

func main() {
	slice1 := []int{1, 2, 3}
	fmt.Println(sum(slice1))
	fmt.Println(half(1.0))

	fmt.Println(half(2.0))

	fmt.Println(largest(slice1))
}

func sum(args ...[]int) int {
	total := 0
	for _, value := range args[0] {
		total += value
	}
	return total
}

func half(x int) (int, bool) {
	return x / 2, x%2 == 0
}

func largest(args ...[]int) int {
	output := 0
	for _, value := range args[0] {
		if value > output {
			output = value
		}
	}
	return output
}
