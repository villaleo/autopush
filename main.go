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
		"message to supply when an automatic commit is made")

	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func main() {
	flag.Parse()

	printAsciiArt()
	for {
		stageWorkingDir()
		if willCommitChanges() {
			promptForCustomCommitMsg()
			pushCommittedChanges()
		} else {
			sleep()
		}
	}
}

// stageWorkingDir stages all the files in the working directory using git add -A.
func stageWorkingDir() {
	cmd := exec.Command("git", "add", "-A")
	if out, err := cmd.Output(); err != nil {
		slog.Error(red("failed to stage files"))
		if res := string(out); res != "" {
			slog.Error(red(res))
		}

		sleep()
	}
}

// promptForCustomCommitMsg prompts the user to enter a custom commit message.
// If no message is entered after a minute, the default commit message is used
// instead.
func promptForCustomCommitMsg() {
	slog.Info(green("changes detected"))

	customMsg := make(chan string, 1)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	go func() {
		fmt.Println("enter a commit message: ")

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			slog.Error(red(err))
			return
		}

		select {
		case customMsg <- line:
		case <-ctx.Done():
		}
	}()

	select {
	case msg := <-customMsg:
		*commitMsg = msg
	case <-ctx.Done():
		slog.Info("no message provided.")
		slog.Info(fmt.Sprintf("falling back to default: %q", *commitMsg))
	}
}

// willCommitChanges commits the staging area in the working directory using git commit
// -m msg. Returns true if a commit was made.
func willCommitChanges() bool {
	cmd := exec.Command("git", "commit", "-m", *commitMsg)
	out, err := cmd.Output()
	res := string(out)
	if strings.Contains(res, "up to date") || strings.Contains(res, "nothing to commit") {
		sleep()
		return false
	}

	if err != nil {
		slog.Error(red("failed to commit changes"))
		if res != "" {
			slog.Error(red(res))
		}
		sleep()
	}

	return true
}

// pushCommittedChanges pushes the changes in the working directory using git push.
func pushCommittedChanges() {
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

	slog.Info(green("changes pushed"))
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
