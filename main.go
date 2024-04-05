package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	TimerFileDir = "/.timer1/"
	TimerFile    = "timer.txt"
	FileMode     = 0755
)

var timerFileDir = os.Getenv("HOME") + TimerFileDir
var timerFile = timerFileDir + TimerFile

var rootCmd = &cobra.Command{
	Use:   "Timer1",
	Short: "Timer1 is a simple CLI application to track time spent on a certain project",
	Run: func(cmd *cobra.Command, args []string) {
		println("timer1 start <project-name>")
		println("timer1 ls")
		println("timer1 stop")
	},
}

func writeToFile(data string) {
	err := os.WriteFile(timerFile, []byte(data), FileMode)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}
}

func startTimer(projectName string) {
	data := fmt.Sprintf("%s\n%d\n%d", projectName, time.Now().Unix(), 0)
	writeToFile(data)
	fmt.Printf("Started timer for project: %s\n", projectName)
}

func timerDetails() (*string, int, int, bool, int64, int64) {
	data, err := os.ReadFile(timerFile)
	if err != nil {
		fmt.Println("No timer is currently running.")
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	projectName := lines[0]
	startTime, err := strconv.ParseInt(lines[1], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing timer info: %v\nStop current timer to remove the error\n", err)
		return nil, 0, 0, false, 0, 0
	}
	stopTime, err := strconv.ParseInt(lines[2], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing timer info: %v\nStop current timer to remove the error\n", err)
	}
	isStopped := stopTime != 0

	elapsed := time.Since(time.Unix(startTime, 0))
	if isStopped {
		elapsed = time.Unix(stopTime, 0).Sub(time.Unix(startTime, 0))
	}
	elapsedMinutes := int(elapsed.Minutes())
	elapsedSeconds := int(elapsed.Seconds()) % 60 // remaining seconds after minutes

	return &projectName, elapsedMinutes, elapsedSeconds, isStopped, startTime, stopTime
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a timer for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// make timerFile if not already exists
		if _, err := os.Stat(timerFileDir); os.IsNotExist(err) {
			err := os.MkdirAll(timerFileDir, FileMode)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = os.Create(timerFile)
			if err != nil {
				fmt.Println(err)
				return
			}

			startTimer(args[0])
		} else {
			projectName, _, _, isStopped, _, _ := timerDetails()
			if !isStopped {
				fmt.Printf("A timer is already running for project: %s\n", *projectName)
				return
			} else {
				startTimer(args[0])
			}
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List currently running timer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, elapsedMinutes, elapsedSeconds, isStopped, _, _ := timerDetails()
		if projectName != nil {
			if isStopped {
				println("Timer is stopped")
			}
			details := fmt.Sprintf("Project: %s, elapsed time: %d minutes %d seconds \n", *projectName, elapsedMinutes, elapsedSeconds)
			println(details)
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop currently running timer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, elapsedMinutes, elapsedSeconds, isStopped, startTime, _ := timerDetails()
		if projectName != nil {
			if isStopped {
				println("No Timer is running currently.")
				return
			}
			details := fmt.Sprintf("Project: %s, elapsed time: %d minutes %d seconds \n", *projectName, elapsedMinutes, elapsedSeconds)
			println(details)
			println("Succesfully stopped project:", *projectName)
		}
		// stop the timer by adding stop time at 3rd line of the file
		data := fmt.Sprintf("%s\n%d\n%d", *projectName, startTime, time.Now().Unix())
		writeToFile(data)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(stopCmd)
}

func main() {
	rootCmd.Execute()
}
