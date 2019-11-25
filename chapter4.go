// chapter4
package main

import "fmt"

func main() {
	fmt.Print("Enter feet: ")
	var Ft float64
	fmt.Scanf("%f", &Ft)

	M := Ft * 0.3048
	fmt.Println(Ft, "\nMeters: ", M)
}
