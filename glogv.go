package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-json"
)

var (
	colorReset  = "\033[0m"
	colorGray   = "\033[90m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type ljson struct {
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Error   string    `json:"error"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := reformat(scanner.Bytes()); err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(4)
	}
}

func reformat(b []byte) error {
	jsn := &ljson{}
	if err := json.Unmarshal(b, jsn); err != nil {
		return err
	}

	fmt.Printf(
		"%s %s %s%s\n",
		formatTime(jsn.Time),
		formatLevel(jsn.Level),
		formatMessage(jsn.Message, jsn.Level),
		formatError(jsn.Error))

	return nil
}

func formatTime(t time.Time) string {
	return colorGray + t.Format(time.Kitchen)
}

func formatLevel(s string) string {
	switch s {
	case "info":
		return colorGreen + "INF"
	case "warn":
		return colorYellow + "WRN"
	case "debug":
		return colorCyan + "DBG"
	case "error":
		return colorRed + "ERR"
	case "panic":
		return colorPurple + "PNC"
	default:
		return "???"
	}
}

func formatMessage(s string, l string) string {
	if l == "info" {
		return colorReset + s
	}
	return s
}

func formatError(s string) string {
	if s == "" {
		return s
	}

	return " (error: " + s + ")"
}
