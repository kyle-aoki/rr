package rr

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

const debugFlag = "--debug"
const shortDebugFlag = "-d"

type fileFuncPair struct {
	filename string
	fn       string
}

func Check(err error) {
	defer customRecover(err)
	if err != nil {
		panic(err)
	}
}

func customRecover(err error) {
	if r := recover(); r != nil {
		if hasDebugFlag() {
			printCustomStacktrace(err)
		} else {
			fmt.Println(r)
		}
		os.Exit(0)
	}
}

func printCustomStacktrace(err error) {
	stacktrace := string(debug.Stack())
	ffps := fileFuncPairSliceFromStacktrace(stacktrace)
	maxFileNameLength := findMaxFileNameLength(ffps)
	dashedline(err)
	fmt.Println(err)
	dashedline(err)
	for i := 0; i < len(ffps); i++ {
		spaces := strings.Repeat(" ", maxFileNameLength-len(ffps[i].filename))
		fmt.Printf("%s%s :: %s", ffps[i].filename, spaces, ffps[i].fn)
	}
}

func hasDebugFlag() bool {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == debugFlag || os.Args[i] == shortDebugFlag {
			return true
		}
	}
	return false
}

func dashedline(err error) {
	fmt.Println(strings.Repeat("-", len(err.Error())))
}

func findMaxFileNameLength(ffs []fileFuncPair) int {
	var maxFileNameLength int
	for i := 0; i < len(ffs); i++ {
		maxFileNameLength = max(maxFileNameLength, len(ffs[i].filename))
	}
	return maxFileNameLength
}

func fileFuncPairSliceFromStacktrace(stacktrace string) []fileFuncPair {
	lines := strings.Split(stacktrace, "\n")
	var ffs []fileFuncPair
	// skip first line: "goroutine 1 [running]:"
	for i := 1; i < len(lines)-1; i += 2 {
		ffs = append(ffs, toFileFuncPair(lines[i], lines[i+1]))
	}
	ffs = ffs[4:] // omit stackframes from this file
	return ffs
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// firstline:  main.main()
// secondline: /Users/kyle/check/check.go:14 +0xd4
func toFileFuncPair(firstline, secondline string) fileFuncPair {
	fullFilePath := removeMemoryAddress(secondline)
	filename := getFileName(fullFilePath)
	return fileFuncPair{filename: filename, fn: firstline}
}

func removeMemoryAddress(line string) string {
	if !strings.Contains(line, "+0") {
		return line
	}
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return line
	}
	return parts[0]
}

func getFileName(fullFilePath string) string {
	parts := strings.Split(fullFilePath, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
