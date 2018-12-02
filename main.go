package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	if len(os.Args) == 2 {
		switch door := os.Args[1]; door {
		case "1_1":
			door1_1(scan_digits())
		case "1_2":
			door1_2(scan_digits())
		case "2_1":
			door2_1(scan_strings())
		case "2_2":
			door2_2(scan_strings())
		default:
			fmt.Println("Invalid door!")
		}
	}

}

func scan_digits() []int {

	digits := make([]int, 0)
	var digit int

	for _, err := fmt.Fscan(os.Stdin, &digit); err != io.EOF; _, err = fmt.Fscan(os.Stdin, &digit) {
		if err == nil {
			digits = append(digits, digit)
		}
	}

	return digits
}

func scan_strings() []string {
	scanner := bufio.NewScanner(os.Stdin)
	strings := make([]string, 0)

	for scanner.Scan() {
		strings = append(strings, scanner.Text())
	}

	return strings
}

func door1_1(digits []int) {
	freq := 0
	for i := 0; i < len(digits); i++ {
		freq += digits[i]
	}

	fmt.Println("Frequency:", freq)
}

func door1_2(digits []int) {
	set := make(map[int]struct{})
	exists := struct{}{}

	freq := 0
	digitsLen := len(digits)
	for i := 0; true; i++ {
		freq += digits[i%digitsLen]
		if _, ok := set[freq]; !ok {
			set[freq] = exists
		} else {
			break
		}
	}

	fmt.Println("Repeated frequency:", freq)
}

func door2_1(strings []string) {
	twos := 0
	threes := 0

	for _, line := range strings {
		two, three := checkLine(line)

		if two {
			twos++
		}

		if three {
			threes++
		}

	}
	fmt.Println("Twos:", twos, ", threes:", threes, " twos*threes:", twos*threes)
	/**/
}

func checkLine(line string) (hasTwo bool, hasThree bool) {

	charCounts := make(map[rune]int)
	for _, c := range line {
		charCounts[c] += 1
	}

	for _, v := range charCounts {
		hasTwo = hasTwo || v == 2
		hasThree = hasThree || v == 3
	}

	return
}

func door2_2(text []string) {
	for i, line := range text {
		for _, laterLine := range text[i+1:] {
			if laterLine == "" {
				continue
			}
			dist, char := hammingDistance(line, laterLine)
			if dist == 1 {
				fmt.Printf("Char: %s\n", string(char))
				fmt.Println(line)
				fmt.Println(laterLine)
				fmt.Println(strings.Replace(line, string(char), " ", -1))
			}
		}
	}
}

func hammingDistance(a string, b string) (distance int, lastDiff rune) {
	for i, char := range a {
		if char != rune(b[i]) {
			distance++
			lastDiff = char
		}
	}
	return
}
