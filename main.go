package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	timeout = flag.Int("t", 60, "the timeout until the the script is "+
		"re-ran. \"m\" is the default timeout strategy. use \"tstr\" to "+
		"specify a timeout strategy.")
	timeoutStrategy = flag.String("tstr", "ms", "the timeout strategy (unit). "+
		"accepted values are \"ms\", \"s\", and \"m\".")
	commitMsg = flag.String("msg", "automated commit by autopush", "the commit "+
		"message to supply when an automatic commit is made.")

	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func main() {
	flag.Parse()

	printAsciiArt()
	for {
		if !isWorkingDirClean() {
			if ok := stageWorkingDir(); !ok {
				sleep()
				continue
			}
			if ok := promptForCustomCommitMsg(); !ok {
				sleep()
				continue
			}
			if ok := commitStagedChanges(); !ok {
				sleep()
				continue
			}
			if ok := pushCommittedChanges(); !ok {
				sleep()
				continue
			}
		}
		sleep()
	}
}

// isWorkingDirClean checks if the working directory is clean using git status.
// If an error occurrs when checking git status, the directory is assumed to be
// dirty.
func isWorkingDirClean() bool {
	cmd := exec.Command("git", "status")

	out, err := cmd.Output()
	res := string(out)
	if err != nil {
		fmt.Println(red(err))
		fmt.Println(red(res))
		return false
	}

	fmt.Println(res)
	return strings.Contains(res, "nothing to commit, working tree clean")
}

// stageWorkingDir stages all the files in the working directory using git add -A.
// Returns false if an error occurred. Errors are printed to Stdout.
func stageWorkingDir() (ok bool) {
	slog.Info("changes detected.")
	cmd := exec.Command("git", "add", "-A")

	out, err := cmd.Output()
	res := string(out)
	if err != nil {
		fmt.Println(red(err))
		fmt.Println(red(res))
		return false
	}

	slog.Info(green("all files staged."))
	return true
}

// promptForCustomCommitMsg prompts the user to enter a custom commit message.
// If no message is entered after a minute, the default commit message is used
// instead.
// Returns false if an error occurred. Errors are printed to Stdout.
func promptForCustomCommitMsg() (ok bool) {
	customMsg := make(chan string, 1)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	go func() {
		fmt.Println("enter a commit message: ")

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(red(err))
			return
		}

		select {
		case customMsg <- line:
		case <-ctx.Done():
		}

		ok = true
	}()

	select {
	case msg := <-customMsg:
		*commitMsg = msg
	case <-ctx.Done():
		slog.Info("no message provided.")
		slog.Info(fmt.Sprintf("falling back to default: %q", *commitMsg))
	}

	return ok
}

// commitStagedChanges commits the staging area in the working directory using git commit
// -m commitMsg.
// Returns false if an error occurred. Errors are printed to Stdout.
func commitStagedChanges() (ok bool) {
	cmd := exec.Command("git", "commit", "-m", *commitMsg)

	out, err := cmd.Output()
	res := string(out)
	if err != nil {
		fmt.Println(red(err))
		fmt.Println(red(res))
		return false
	}

	slog.Info(green("changes committed."))
	return true
}

// pushCommittedChanges pushes the changes in the working directory using git push.
// Returns false if an error occurred. Errors are printed to Stdout.
func pushCommittedChanges() (ok bool) {
	cmd := exec.Command("git", "push")

	out, err := cmd.Output()
	res := string(out)
	if err != nil {
		fmt.Println(red(err))
		fmt.Println(red(res))
		return false
	}

	slog.Info(green("changes pushed."))
	return true
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
	fmt.Println("autopush v1.0 created by github.com/villaleo")
}
