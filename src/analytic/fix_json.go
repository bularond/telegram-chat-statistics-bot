package analytic

import (
	"strings"
)

func RemoveBrokenLinesFromJsonFile(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	brokenLines := make([]int, 0)

	for i, line := range lines {
		if strings.Contains(line, "custom_emoji") {
			brokenLines = append(brokenLines, i)
		}
	}
	if len(brokenLines) == 0 {
		return data
	}

	leftPointer := 0
	rightPointer := 0
	brokenPointer := 0
	for {
		if rightPointer == len(lines) {
			break
		}
		for brokenPointer != len(brokenLines) && rightPointer == brokenLines[brokenPointer]-1 {
			rightPointer += 5
			brokenPointer += 1
		}

		if leftPointer != rightPointer {
			lines[leftPointer] = lines[rightPointer]
		}

		leftPointer++
		rightPointer++
	}

	lines = lines[:leftPointer]
	return []byte(strings.Join(lines, ""))
}
