package commands

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var reader *bufio.Reader

func inputDate() time.Time {
	var input string

	fmt.Print("Year: ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	year, _ := strconv.Atoi(input)

	fmt.Print("Month (1-12): ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	month, _ := strconv.Atoi(input)

	fmt.Print("Day in month (1-31): ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	day, _ := strconv.Atoi(input)

	fmt.Print("Hour (0-23): ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	hour, _ := strconv.Atoi(input)

	fmt.Print("Minute (0-59): ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimRight(input, "\n")
	minute, _ := strconv.Atoi(input)

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC).UTC()
}
