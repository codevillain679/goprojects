// chapter7.1
package main

import "fmt"

func makeOddGenerator() func() int {
	fmt.Println("the factory")
	i := -1
	return func() int {
		fmt.Println("the closure")
		fmt.Println("i are: ", i)
		i += 2
		return i
	}
}

func main() {
	nextOdd := makeOddGenerator()
	fmt.Println(nextOdd()) // 1
	fmt.Println(nextOdd()) // 3
	fmt.Println(nextOdd()) // 5
	fmt.Println(nextOdd()) // 7
	fmt.Println(nextOdd()) // 9
}
