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
	timeout = flag.Int("timeout", 60, "the time that the script will be re-ran, in minutes")

	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func main() {
	flag.Parse()

	fmt.Println("             _                        _")
	fmt.Println("            | |                      | |")
	fmt.Println("  __ _ _   _| |_ ___  _ __  _   _ ___| |__")
	fmt.Println(" / _` | | | | __/ _ \\| '_ \\| | | / __| '_ \\")
	fmt.Println("| (_| | |_| | || (_) | |_) | |_| \\__ \\ | | |")
	fmt.Println(" \\__,_|\\__,_|\\__\\___/| .__/ \\__,_|___/_| |_|")
	fmt.Println("                     | |")
	fmt.Println("                     |_|")
	fmt.Println()

	fmt.Println("autopush v1.0 created by github.com/villaleo")
	for {
		stage()
		if commit("automated commit by autopush") {
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

// sleep suspends the goroutine for timeout minutes.
func sleep() {
	time.Sleep(time.Duration(*timeout) * time.Minute)
}
