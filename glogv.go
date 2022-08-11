package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/goccy/go-json"
)

const (
	bufSize       = 1 * 1024 * 1024 // read buffer size to use.
	maxKeys       = 100             // maximum size json key slice.
	errorExitCode = 4               // exit code if error occurs.
)

// ANSI color escape codes
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

// this struct will be used to marshall the json file into key/values.
type keyValues struct {
	Map map[string]any `json:"-"`
}

// this string slice will store keys and then sort them.
var keys = make([]string, 0, maxKeys)

func main() {
	// create a larger than the default buffer size of 64k.
	buf := make([]byte, 0, bufSize)

	// create a stdin scanner with the custom buffer.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(buf, bufSize)

	// loop until EOF.
	for scanner.Scan() {
		reformat(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("scanner error: %v\n", err)
		os.Exit(errorExitCode)
	}
}

// reformats the json log line into a prettier, more readable version.
func reformat(b []byte) {
	var tm time.Time
	var level, message, errorx string

	// first make sure the log line is json, if not return without processing.
	if string(b[:1]) != "{" {
		return
	}

	// marshall the current log entry into a key/value map.
	keyVals := &keyValues{}
	if err := json.Unmarshal(b, &keyVals.Map); err != nil {
		return
	}

	// first parse and format the standard logging fields.
	if val, ok := keyVals.Map["time"]; ok {
		tm, _ = time.Parse(time.RFC3339, val.(string))
	}
	if val, ok := keyVals.Map["level"]; ok {
		level = val.(string)
	}
	if val, ok := keyVals.Map["message"]; ok {
		message = val.(string)
	}
	if val, ok := keyVals.Map["error"]; ok {
		errorx = val.(string)
	}

	// reformat what we have parsed so far.
	tmStr := formatTime(tm)
	lvlStr := formatLevel(level)
	msgStr := formatMessage(message, level)
	errStr := formatError(errorx)

	// next delete the keys we just processed from the map.
	delete(keyVals.Map, "time")
	delete(keyVals.Map, "level")
	delete(keyVals.Map, "message")
	delete(keyVals.Map, "error")

	// now, parse through the remaining key/values in the map.
	valStr := formatMap(keyVals.Map)

	// finally, print the prettier log entry.
	fmt.Printf("%s%s%s%s%s\n", tmStr, lvlStr, msgStr, errStr, valStr)
}

// formats the 'time' portion of the json log line.
func formatTime(t time.Time) string {
	s := t.Format(time.Kitchen)

	if len(s) == 6 {
		return colorGray + "0" + s
	}

	return colorGray + s
}

// formats the 'level' portion of the json log line.
func formatLevel(s string) string {
	switch s {
	case "info":
		return " " + colorGreen + "INF"
	case "warn":
		return " " + colorYellow + "WRN"
	case "debug":
		return " " + colorCyan + "DBG"
	case "error":
		return " " + colorRed + "ERR"
	case "panic":
		return " " + colorPurple + "PNC"
	default:
		return " ???"
	}
}

// formats the 'message' portion of the json log line.
func formatMessage(s string, l string) string {
	if s == "" {
		return s
	}

	if l == "info" {
		return " " + colorReset + s
	}

	return " " + s
}

// formats the 'error' portion of the json log line.
func formatError(s string) string {
	if s == "" {
		return s
	}

	return " (error: " + s + ")"
}

// formats the remaining key/value pairs of the json log line.
func formatMap(m map[string]any) string {
	l := len(m)
	// if the map is empty then return nothing.
	if l == 0 {
		return ""
	}

	// if there is just one value left in the map, return it now.
	if l == 1 {
		for k, v := range m {
			return " " + colorGray + k + "=" + colorReset + v.(string)
		}
	}

	// there is more than 1 value in the map, so we will sort by
	// key to get a consistent order.
	keys = keys[:0]
	i := 0
	for k := range m {
		keys = append(keys, k)
		i++
		if i > maxKeys {
			break
		}
	}

	sort.Strings(keys)

	var s string
	for _, k := range keys {
		s += " " + colorGray + k + "=" + colorReset + m[k].(string)
	}

	return s
}
