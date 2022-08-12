package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/klauspost/compress/gzip"
)

const (
	maxKeys       = 100 // maximum size of json key slice.
	errorExitCode = 4   // exit code if an error occurs.
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

// default colors by level.
var color = map[string]string{
	"info":  colorGreen,
	"warn":  colorYellow,
	"debug": colorCyan,
	"error": colorRed,
	"panic": colorPurple,
	"fatal": colorPurple,
	"trace": colorCyan,
}

// other default colors.
var (
	timeColor = colorGray
	tagColor  = colorGray
	infoColor = colorWhite
)

// this struct will be used to marshall the json file into key/values.
type keyValues struct {
	Map map[string]any `json:"-"`
}

// this string slice will store keys and then sort them.
var keys = make([]string, 0, maxKeys)

func main() {
	// parse flags
	tailFile := flag.Bool("tail", false, "tail the file provided")
	flag.Parse()
	files := flag.Args()

	// make sure there is a file provided if the -tail option is set
	if *tailFile && len(files) == 0 {
		fmt.Printf("-tail option used without a file being provided\n")
		os.Exit(errorExitCode)
	}

	// check for tail mode if flag set.
	if *tailFile {
		if err := tail(files); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(errorExitCode)
		}
		return
	}

	// check for cat mode if not tail mode and file provided.
	if len(files) > 0 {
		if err := cat(files); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(errorExitCode)
		}
		return
	}

	// otherwise, scan stdin.
	if err := scan(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(errorExitCode)
	}
}

// scan continues to scan stdin until EOF.
func scan() error {
	scanner := bufio.NewScanner(os.Stdin)

	// loop until EOF.
	for scanner.Scan() {
		reformat(scanner.Bytes())
	}

	return scanner.Err()
}

// tail will run the linux tail command and log the output
func tail(files []string) error {
	// check if file(s) exists first
	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			return err
		}
	}

	args := []string{"--follow=name"}
	args = append(args, files...)

	cmd := exec.CommandContext(context.Background(), "tail", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			reformat(scanner.Bytes())
		}
		wg.Done()
	}()

	if err = cmd.Start(); err != nil {
		return err
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	wg.Wait()

	return cmd.Wait()
}

// cat will read the given file(s) and reformat it
func cat(files []string) error {
	fn := func(file string) error {
		read, err := os.Open(file)
		if err != nil {
			return err
		}
		defer read.Close()

		// pick a reader based on if the file is compressed or not.
		var scanner *bufio.Scanner
		if filepath.Ext(file) == ".gz" {
			gz, err := gzip.NewReader(read)
			if err != nil {
				return err
			}
			defer gz.Close()
			scanner = bufio.NewScanner(gz)
		} else {
			scanner = bufio.NewScanner(read)
		}

		// loop until EOF.
		for scanner.Scan() {
			reformat(scanner.Bytes())
		}

		return scanner.Err()
	}

	for _, file := range files {
		if err := fn(file); err != nil {
			return err
		}
	}

	return nil
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

	// if level is unknown, set it to default
	if _, ok := color[level]; !ok {
		level = "info"
	}

	// now, parse through the remaining key/values in the map.
	valStr := formatMap(keyVals.Map, level)

	// finally, print the prettier log entry.
	fmt.Printf("%s%s%s%s%s\n", tmStr, lvlStr, msgStr, errStr, valStr)
}

// formats the 'time' portion of the json log line.
func formatTime(t time.Time) string {
	s := t.Format(time.Kitchen)

	if len(s) == 6 {
		return timeColor + "0" + s
	}

	return timeColor + s
}

// formats the 'level' portion of the json log line.
func formatLevel(s string) string {
	switch s {
	case "info":
		return " " + color[s] + "INF"
	case "warn":
		return " " + color[s] + "WRN"
	case "debug":
		return " " + color[s] + "DBG"
	case "error":
		return " " + color[s] + "ERR"
	case "panic":
		return " " + color[s] + "PNC"
	case "fatal":
		return " " + color[s] + "PNC"
	case "trace":
		return " " + color[s] + "PNC"
	default:
		return color["info"] + " ???"
	}
}

// formats the 'message' portion of the json log line.
func formatMessage(s string, l string) string {
	if s == "" {
		return s
	}

	if l == "info" {
		return " " + infoColor + s
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
func formatMap(m map[string]any, l string) string {
	length := len(m)
	// if the map is empty then return nothing.
	if length == 0 {
		return ""
	}

	// compute value color
	var clr string
	if l == "info" {
		clr = infoColor
	} else {
		clr = color[l]
	}

	// if there is just one value left in the map, return it now.
	if length == 1 {
		for k, v := range m {
			return " " + tagColor + k + "=" + clr + v.(string)
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
		s += " " + tagColor + k + "=" + clr + m[k].(string)
	}

	return s
}
