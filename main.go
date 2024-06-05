package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	timeout = flag.Int("t", 60, "the timeout until the the script is"+
		"re-ran. \"m\" is the default timeout strategy. use \"tstr\" to"+
		"specify a timeout strategy.")
	timeoutStrategy = flag.String("tstr", "ms", "the timeout strategy (unit)."+
		"accepted values are \"ms\", \"s\", and \"m\".")
	msg = flag.String("msg", "automated commit by autopush", "the commit"+
		"message to supply when an automatic commit is made")

	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func main() {
	flag.Parse()

	printAsciiArt()
	fmt.Println("autopush v1.0 created by github.com/villaleo")

	for {
		stage()
		if commit(*msg) {
			push()
			slog.Info(green("changes pushed"))
		} else {
			sleep()
		}
	}
}

// stage stages all the files in the working directory using git add -A.
func stage() {
	cmd := exec.Command("git", "add", "-A")
	if out, err := cmd.Output(); err != nil {
		slog.Error(red("failed to stage files"))
		if res := string(out); res != "" {
			slog.Error(red(res))
		}

		sleep()
	}
}

// commit commits the staging area in the working directory using git commit
// -m msg. Returns true if a commit was made.
func commit(msg string) bool {
	cmd := exec.Command("git", "commit", "-m", msg)
	out, err := cmd.Output()
	res := string(out)
	if strings.Contains(res, "up to date") || strings.Contains(res, "nothing to commit") {
		sleep()
		return false
	}

	if err != nil {
		slog.Error(red("failed to commit changes"))
		if msg != "" {
			slog.Error(red(res))
		}
		sleep()
	}

	return true
}

// push pushes the changes in the working directory using git push.
func push() {
	cmd := exec.Command("git", "push")
	out, err := cmd.Output()
	if err != nil {
		slog.Error(red("failed to push changes"))

		if res := string(out); res != "" {
			slog.Error(red(res))
		}
		sleep()

		return
	}
}

// sleep suspends the goroutine for timeout using timeoutStrategy units.
func sleep() {
	duration := time.Minute

	if timeoutStrategy == nil {
		time.Sleep(time.Duration(*timeout) * duration)
		return
	}

	switch *timeoutStrategy {
	case "ms":
		duration = time.Millisecond
	case "s":
		duration = time.Second
	case "m":
		duration = time.Minute
	}

	time.Sleep(time.Duration(*timeout) * duration)
}

// printAsciiArt prints ASCII art reading "autopush" to stdout.
func printAsciiArt() {
	fmt.Println("             _                        _")
	fmt.Println("            | |                      | |")
	fmt.Println("  __ _ _   _| |_ ___  _ __  _   _ ___| |__")
	fmt.Println(" / _` | | | | __/ _ \\| '_ \\| | | / __| '_ \\")
	fmt.Println("| (_| | |_| | || (_) | |_) | |_| \\__ \\ | | |")
	fmt.Println(" \\__,_|\\__,_|\\__\\___/| .__/ \\__,_|___/_| |_|")
	fmt.Println("                     | |")
	fmt.Println("                     |_|")
	fmt.Println()
}
