package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
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
		case "3_1":
			door3_1()
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

type claim struct {
	id     int
	x      int
	y      int
	width  int
	height int
}

func door3_1() {
	const size = 1000
	scanner := bufio.NewScanner(os.Stdin)
	claims := [size][size][]claim{}
	var claims_re = regexp.MustCompile(`#(\d+) @ (\d+),(\d+): (\d+)x(\d+)`)

	overlaps := 0
	for scanner.Scan() {

		line := scanner.Text()
		matches := claims_re.FindStringSubmatch(line)
		id, _ := strconv.Atoi(matches[1])
		x, _ := strconv.Atoi(matches[2])
		y, _ := strconv.Atoi(matches[3])
		w, _ := strconv.Atoi(matches[4])
		h, _ := strconv.Atoi(matches[5])
		c := claim{id, x, y, w, h}
		for i := x; i < x+w; i++ {
			for j := y; j < y+h; j++ {
				claims[i][j] = append(claims[i][j], c)
				if len(claims[i][j]) == 2 {
					overlaps++
				}
			}
		}
	}
	fmt.Println("Overlapping size: ", overlaps)

	// Nasty code o'hoy!
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			claims_in_bucket := claims[i][j]
			// check if claim is no overlap candidate
			if len(claims_in_bucket) == 1 && claims_in_bucket[0].x == i && claims_in_bucket[0].y == j {
				no_overlap_found := true
				c2 := claims_in_bucket[0]

				// exit loop early if no_overlap_found is false
				for a := c2.x; no_overlap_found && a < c2.x+c2.width; a++ {
					for b := c2.y; no_overlap_found && b < c2.y+c2.height; b++ {
						// check if there is overlap at position a,b
						if len(claims[a][b]) != 1 {
							no_overlap_found = false
						}
					}
				}
				// Found the non overlapping claim!
				if no_overlap_found {
					fmt.Println("Non overlapping claim: ", c2.id)
				}
			}
		}
	}

}
