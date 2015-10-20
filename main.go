package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		errorln("Usage: timeconv <timestamp> <timezone>")
		return
	}

	timestamp, err := time.Parse(time.RFC3339, flag.Arg(0))
	if err != nil {
		errorln("Invalid timestamp: %s", err.Error())
		return
	}

	location, ok := parseOffset(flag.Arg(1))
	if !ok {
		location, err = time.LoadLocation(flag.Arg(1))
		if err != nil {
			errorln("Unknown timezone %s", flag.Arg(1))
			return
		}
	}

	fmt.Println(timestamp.In(location).Format(time.RFC3339))
}

func parseOffset(input string) (*time.Location, bool) {
	var sign rune
	var hh, mm int
	if _, err := fmt.Sscanf(input, `%c%2d:%2d`, &sign, &hh, &mm); err != nil {
		return nil, false
	}

	if mm < 0 {
		return nil, false
	}

	offset := (3600 * hh) + (60 * mm)
	if sign == '-' {
		offset *= -1
	} else if sign != '+' {
		return nil, false
	}

	return time.FixedZone(input, offset), true
}

func errorln(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, color.RedString(format, args...))
}
