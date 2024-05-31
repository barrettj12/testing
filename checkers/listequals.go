// Copyright 2024 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package checkers

import (
	"fmt"
	gc "gopkg.in/check.v1"
)

type listEqualsChecker struct {
	*gc.CheckerInfo
}

// The ListEquals checker verifies if two lists are equal. If they are not,
// it will essentially run a "diff" algorithm to provide the developer with
// an easily understandable summary of the difference between the two lists.
var ListEquals gc.Checker = &listEqualsChecker{
	&gc.CheckerInfo{Name: "ListEquals", Params: []string{"obtained", "expected"}},
}

func (l listEqualsChecker) Check(params []interface{}, names []string) (result bool, error string) {
	obtained := params[0]
	expected := params[1]

	// TODO: Simple pre-checks:
	// - both obtained and expected are indeed slices
	// - of the same element type
	// - and this element type is comparable.
	_, _ = obtained, expected

	// The approach here is to find a longest-common subsequence using dynamic
	// programming, and use this to generate the diff. This algorithm runs in
	// O(n^2). However, naive list equality is only O(n). Hence, to be more
	// efficient, we should first check if the lists are equal, and if they are
	// not, we do the more complicated work to find out exactly *how* they are
	// different.

	// Check length is equal
	// Iterate through and check every element

	// If we're here, the lists are not equal, so run the DP algorithm to
	// compute the diff.
	return false, "" //generateDiff(obtained, expected)
}

func generateDiff(obtained, expected []rune) string {
	// lenLCS[m][n] stores the length of the longest common subsequence of
	// obtained[:m] and expected[:n]
	lenLCS := make([][]int, len(obtained)+1)
	for i := 0; i <= len(obtained); i++ {
		lenLCS[i] = make([]int, len(expected)+1)
	}

	// lenLCS[i][0] and lenLCS[0][j] are already correctly initialised to 0

	for i := 1; i <= len(obtained); i++ {
		for j := 1; j <= len(expected); j++ {
			if obtained[i-1] == expected[j-1] {
				// We can extend the longest subsequence of obtained[:i-1] and expected[:j-1]
				lenLCS[i][j] = lenLCS[i-1][j-1] + 1
			} else {
				// We can't extend a previous subsequence
				lenLCS[i][j] = max(lenLCS[i-1][j], lenLCS[i][j-1])
			}
		}
	}

	// print table
	fmt.Print("      ")
	for j := 0; j < len(expected); j++ {
		fmt.Printf("%c  ", expected[j])
	}
	fmt.Println()
	for i := 0; i <= len(obtained); i++ {
		fmt.Printf("%c  ", append([]rune{' '}, obtained...)[i])
		for j := 0; j <= len(expected); j++ {
			fmt.Printf("%d  ", lenLCS[i][j])
		}
		fmt.Println()
	}

	// "Traceback" to calculate the diff
	var diffs []diff
	i := len(obtained)
	j := len(expected)

	for i > 0 && j > 0 {
		if lenLCS[i][j] == lenLCS[i-1][j-1] {
			// Element changed at this index
			diffs = append(diffs, elementChanged{j, expected[j-1], obtained[i-1]})
			i -= 1
			j -= 1

		} else if lenLCS[i][j] == lenLCS[i-1][j] {
			// Additional/unexpected element at this index
			diffs = append(diffs, elementAdded{j, obtained[i-1]})
			i -= 1

		} else if lenLCS[i][j] == lenLCS[i][j-1] {
			// Element missing at this index
			diffs = append(diffs, elementRemoved{j, expected[j-1]})
			j -= 1

		} else {
			// Elements are the same at this index - no diff
			i -= 1
			j -= 1
		}
	}

	// Convert diffs array into human-readable error
	description := ""
	for k := len(diffs) - 1; k >= 0; k-- {
		description += diffs[k].String() + "\n"
	}
	fmt.Println(description)
	return description
}

var GenerateDiff = generateDiff

// diff represents a single difference between the two slices.
type diff interface {
	// Prints a user friendly description of this difference.
	String() string
}

type elementAdded struct {
	index   int
	element any
}

func (d elementAdded) String() string {
	return fmt.Sprintf("at index %d: unexpected element %v", d.index, d.element)
}

type elementChanged struct {
	index int

	original, changed any
}

func (d elementChanged) String() string {
	return fmt.Sprintf("at index %d: obtained element %v, expected %v", d.index, d.original, d.changed)
}

type elementRemoved struct {
	index   int
	element any
}

func (d elementRemoved) String() string {
	return fmt.Sprintf("at index %d: missing element %v", d.index, d.element)
}
