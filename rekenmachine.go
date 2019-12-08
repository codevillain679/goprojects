// rekenmachine
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	for {
		var sign string = scanSign()
		a, b := scanDigits()
		var result = calculateResult(sign, a, b)
		fmt.Println(a, sign, b, " = ", result)
	}
}

func scanSign() string {
	fmt.Printf("Optellen(+), aftrekken(-), vermenigvuldigen(*) of delen(/)?")
	var sign string
	fmt.Scan(&sign)
	if sign == "+" || sign == "-" || sign == "*" || sign == "/" {
		return sign
	}
	printRetry()
	return scanSign()
}

func scanDigits() (int, int) {
	fmt.Printf("Voer twee hele getallen in gescheiden met een komma ")
	var input string
	fmt.Scan(&input)
	var s = strings.Split(input, ",")
	if len(s) < 2 {
		printRetry()
		return scanDigits()
	}
	var dig1, err = strconv.Atoi(s[0])
	var dig2, err2 = strconv.Atoi(s[1])
	if err != nil || err2 != nil {

		printRetry()
		return scanDigits()
	}
	return dig1, dig2
}

func calculateResult(sign string, dig1 int, dig2 int) float64 {
	if sign == "+" {
		return float64(dig1 + dig2)
	}
	if sign == "-" {
		return float64(dig1 - dig2)
	}
	if sign == "*" {
		return float64(dig1 * dig2)
	}
	if sign == "/" {
		return float64(dig1) / float64(dig2)
	}
	return 0
}

func printRetry() {
	fmt.Println("Dat is geen geldige invoer, probeer het opnieuw!")
}
