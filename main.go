package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
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
		case "4":
			door4(scan_strings())
		case "5":
			door5(scan_strings())
		case "6":
			door6(scan_strings())
		case "7":
			door7(scan_strings())
		case "8":
			door8()
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

const guard = `\[.+\] Guard #(\d+) .+`
const sleep = `\[.+:0?(\d|\d\d)\] falls.+`
const awakes = `\[.+:0?(\d|\d\d)\] wakes.+`

type span struct {
	start int
	end   int
}

func door4(lines []string) {
	sort.Strings(lines)

	var guardRe = regexp.MustCompile(guard)
	var sleepRe = regexp.MustCompile(sleep)
	var awakesRe = regexp.MustCompile(awakes)
	var curGuard int
	var curSleepStart int
	var curSpans []span
	var curSpan span
	guardSleep := make(map[int]int)
	guardSleepSpans := make(map[int][][]span)
	for _, line := range lines {

		if matches := guardRe.FindStringSubmatch(line); len(matches) > 0 {
			if curSpans != nil {
				guardSleepSpans[curGuard] = append(guardSleepSpans[curGuard], curSpans)
			}
			curGuard, _ = strconv.Atoi(matches[1])
			curSpans = make([]span, 0)
		} else if matches := sleepRe.FindStringSubmatch(line); len(matches) > 0 {
			curSleepStart, _ = strconv.Atoi(matches[1])
			curSpan = span{curSleepStart, -1}
		} else if matches := awakesRe.FindStringSubmatch(line); len(matches) > 0 {
			awokeAt, _ := strconv.Atoi(matches[1])
			curSpan.end = awokeAt
			curSpans = append(curSpans, curSpan)
			time := awokeAt - curSleepStart

			guardSleep[curGuard] += time
		}
	}

	guardSleepSpans[curGuard] = append(guardSleepSpans[curGuard], curSpans)
	maxGuard := -1
	maxTime := -1
	for g, time := range guardSleep {
		if time > maxTime {
			maxTime = time
			maxGuard = g
		}
	}

	days := guardSleepSpans[maxGuard]
	bestMin, bestOverlaps := findBestTimeIn(maxGuard, days)
	fmt.Println("bestMin * guard", maxGuard*bestMin)

	consistentTime := bestMin
	consistentGuard := maxGuard
	maxOverlaps := bestOverlaps
	for guard, guardDays := range guardSleepSpans {

		minute, overlaps := findBestTimeIn(guard, guardDays)
		if overlaps >= maxOverlaps {
			consistentGuard = guard
			consistentTime = minute
			maxOverlaps = overlaps
		}
	}

	fmt.Println("consistentGuard * consistentTime", consistentGuard*(consistentTime), "guard", consistentGuard, "time", consistentTime)

}

func (s span) in(min int) bool {
	return (s.start <= min) && (min < s.end)
}

func findBestTimeIn(g int, days [][]span) (int, int) {
	bestMin := 0
	maxOverlaps := 0
	for i := 0; i < 60; i++ {
		overlaps := 0
		for _, d := range days {
			for _, s := range d {
				if s.in(i) {
					overlaps++
				}
			}
		}
		if overlaps >= maxOverlaps {
			maxOverlaps = overlaps
			bestMin = i
		}
	}

	return bestMin, maxOverlaps
}

func door5(lines []string) {

	//var s scanner.Scanner
	line := lines[0]
	remaining := collapse(line)
	fmt.Println("remaining:", remaining) //, string(remaining))

	minSize := remaining
	for _, c := range "abcdefghijklmnopqrstuvxyz" {
		excludingChar := strings.Replace(strings.Replace(line, string(unicode.ToUpper(c)), "", -1), string(c), "", -1)
		size := collapse(excludingChar)
		if size < minSize {
			minSize = size
		}
	}

	fmt.Println("Minimized:", minSize)
}

func collapse(line string) int {

	remaining := make([]rune, 0)
	reader := strings.NewReader(line)
	var destroys rune
	for r, _, err := reader.ReadRune(); err != io.EOF; r, _, err = reader.ReadRune() {
		if r == destroys && len(remaining) > 1 {
			remaining = remaining[:len(remaining)-1]
			r = remaining[len(remaining)-1]
		} else if r == destroys {
			remaining = remaining[:len(remaining)-1]
			destroys = 0
			continue
		} else {
			remaining = append(remaining, r)
		}

		if unicode.IsLower(r) {
			destroys = unicode.ToUpper(r)
		} else {
			destroys = unicode.ToLower(r)
		}
	}

	return len(remaining)
}

type point struct {
	x           int
	y           int
	id          int
	touchesEdge bool
}

func door6(lines []string) {
	points := make([]*point, 0)
	pointRe := regexp.MustCompile(`(\d+), (\d+)`)

	maxX := 0
	maxY := 0
	minX := math.MaxInt64
	minY := math.MaxInt64

	for i, line := range lines {
		matches := pointRe.FindStringSubmatch(line)
		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])

		points = append(points, &point{x, y, i, false})
		if maxX < x {
			maxX = x
		}
		if minX > x {
			minX = x
		}

		if maxY < y {
			maxY = y
		}
		if minY > y {
			minY = y
		}

	}
	counts := make([]int, len(points))
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			minDist := math.MaxInt64
			var minP *point
			disputed := false
			for _, p := range points {
				dist := manhattanDistance(x, y, p)
				if dist < minDist {
					minDist = dist
					minP = p
					disputed = false
				} else {
					disputed = disputed || dist == minDist
				}
			}

			if !disputed {
				counts[minP.id] += 1
			}
			if x == maxX || x == minX || y == maxY || y == minY {
				minP.touchesEdge = true
			}

		}
	}
	fmt.Println()
	maxCount := 0
	for _, p := range points {

		count := counts[p.id]
		if !p.touchesEdge {
			if maxCount < count {
				maxCount = count
			}
		}

	}

	fmt.Println("largest:", maxCount)
	minX = 0
	minY = 0
	below10kCount := 0
	for x := minX; x < maxX; x++ {
		for y := minY; y < maxY; y++ {
			totDist := 0
			for _, p := range points {
				dist := manhattanDistance(x, y, p)
				totDist += dist
			}

			if totDist <= 10000 {
				below10kCount++
			}

		}
	}

	fmt.Println("below10kCount", below10kCount)
}

func manhattanDistance(x, y int, p *point) int {
	return Abs(x-p.x) + Abs(y-p.y)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func maxInt(v []int) int {
	m := math.MinInt64
	for i, e := range v {
		if e > m {
			fmt.Println("i is bigest", i)
			m = e
		}
	}
	return m
}

func door7(lines []string) {
	instructionRe := regexp.MustCompile(`Step (.+) must be finished before step (.+) can begin.`)
	preToDeps := make(map[string][]string)
	depToPrereqs := make(map[string][]string)

	allSteps := make(map[string]bool)
	isPrereq := make(map[string]bool)
	hasPrereq := make(map[string]bool)

	for _, line := range lines {
		matches := instructionRe.FindStringSubmatch(line)

		prereq := matches[1]
		step := matches[2]
		allSteps[step] = true
		hasPrereq[step] = true

		allSteps[prereq] = true
		isPrereq[prereq] = true

		prereqsForStep := preToDeps[prereq]
		prereqsForStep = append(prereqsForStep, step)
		preToDeps[prereq] = prereqsForStep

		deps := depToPrereqs[step]
		deps = append(deps, prereq)
		depToPrereqs[step] = deps
	}

	var toExploreOrig string
	for s, _ := range allSteps {
		if !hasPrereq[s] {
			toExploreOrig += s
		}
	}

	toExplore := toExploreOrig
	path := ""
	/*	for len(toExplore) > 0 {
			toExplore = SortString(toExplore)

			nextSteps := calculateNextSteps(toExplore, path, depToPrereqs)
			nextStep := strings.Split(nextSteps, "")[0]
			path += nextStep
			toExplore = expandStep(nextStep, path, toExplore, preToDeps)

		}
		fmt.Println(path)
	/**/
	toExplore = toExploreOrig
	path = ""
	times := make([]byte, 5)
	work := make([]string, 5)
	totalTime := 0
	for len(toExplore) > 0 {
		toExplore = SortString(toExplore)
		nextSteps := calculateNextSteps(toExplore, path, depToPrereqs)
		fmt.Println(nextSteps)
		nextStepsSlice := strings.Split(nextSteps, "")
		assigned := ""
		for i := 0; i < len(nextStepsSlice); i++ {
			s := nextStepsSlice[i]
			workerIndex := bytes.IndexByte(times, 0)
			if workerIndex < 0 {
				break
			}
			assigned += s
			stepWorkTime := workTimeFor(s)
			work[workerIndex] = s
			times[workerIndex] = stepWorkTime
			toExplore = expandStep(s, path, toExplore, preToDeps)
		}

		// completeWork
		completesIn := minNonZeroByte(times)
		nextToFinish := bytes.IndexByte(times, completesIn)
		finished := ""
		for nextToFinish > -1 {
			finished += work[nextToFinish]

			times[nextToFinish] = 0
			nextToFinish = bytes.IndexByte(times, completesIn)

		}

		path += SortString(finished)

		for i := range times {
			t := times[i]
			if t == 0 {
				work[i] = ""
			} else {
				times[i] -= completesIn
			}
		}
		totalTime += int(completesIn)

		//assigned, work, times = assignWork(nextSteps, work, times)
		//assigned := SortString(assigned)
		//for _, a := range assigned {

		//}
	}
	fmt.Println(path, totalTime)
	/**/
}

func SortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

func calculateNextSteps(toExplore string, path string, depToPrereqs map[string][]string) string {

	var nextSteps string
	for _, s := range toExplore {
		step := string(s)

		prereqsMet := arePrereqsMet(step, path, depToPrereqs)
		if prereqsMet {
			nextSteps += step
		}
	}

	return nextSteps
}

func arePrereqsMet(step string, path string, depToPrereqs map[string][]string) bool {
	stepDependsOn := depToPrereqs[step]
	prereqsMet := true
	for _, d := range stepDependsOn {
		if !strings.Contains(path, d) {
			prereqsMet = false
			break
		}
	}

	return prereqsMet
}

func expandStep(nextStep, path, toExplore string, preToDeps map[string][]string) string {
	toExplore = strings.Replace(toExplore, nextStep, "", -1)
	deps := preToDeps[nextStep]
	for _, d := range deps {
		//fmt.Println(d, path, toExplore)
		if !(strings.Contains(path, d) || strings.Contains(toExplore, d)) {

			toExplore += d
		}
	}

	return toExplore

}

func workTimeFor(step string) byte {
	t := step[0]
	t = t - 'A' + 61
	return t
}

func minNonZeroByte(v []byte) byte {
	m := byte(255)
	for _, e := range v {
		if e != 0 && e < m {
			m = e
		}
	}
	return m
}

/*
func assignWork(nextSteps string, work [5]string, times [5]int) (string, [5]string, [5]int) {
	return nextSteps, work, times
}/**/

type Node struct {
	meta     []int
	children []*Node
}

func door8() {

	root := ReadNode()

	sum := sumMeta(root)

	fmt.Println(sum)

	rootValue := valueNode(root)

	fmt.Println(rootValue)
}

func sumMeta(node *Node) int {
	sum := 0
	for _, m := range node.meta {
		sum += m
	}

	for _, c := range node.children {
		sum += sumMeta(c)
	}

	return sum
}

func valueNode(node *Node) int {
	if len(node.children) == 0 {
		return sumMeta(node)
	}

	value := 0
	for _, m := range node.meta {
		index := m - 1
		if index >= 0 && index < len(node.children) {
			value += valueNode(node.children[index])
		}
	}
	return value
}

func ReadNode() *Node {
	var childCount int
	var metaCount int
	fmt.Fscan(os.Stdin, &childCount)
	fmt.Fscan(os.Stdin, &metaCount)

	children := make([]*Node, childCount)
	meta := make([]int, metaCount)
	for i := 0; i < childCount; i++ {
		child := ReadNode()
		children[i] = child
	}

	for i := 0; i < metaCount; i++ {
		var m int
		fmt.Fscan(os.Stdin, &m)
		meta[i] = m
	}

	node := Node{meta, children}
	return &node
}
