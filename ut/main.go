package main

import "fmt"

func FizzBuzz(numbers []int) []string {
	var result []string
	for _, number := range numbers {
		switch {
		case number > 0 && number%15 == 0:
			result = append(result, "FizzBuzz")
		case number > 0 && number%3 == 0:
			result = append(result, "Fizz")
		case number > 0 && number%5 == 0:
			result = append(result, "Buzz")
		default:
			result = append(result, fmt.Sprintf("%d", number))
		}
	}
	return result
}
