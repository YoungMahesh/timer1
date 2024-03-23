package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const timerFile = "timer.txt"

var rootCmd = &cobra.Command{
	Use:   "Timer1",
	Short: "Timer1 is a simple CLI application to track time spent on a certain project",
	Run: func(cmd *cobra.Command, args []string) {
		println("timer1 start <project-name>")
		println("timer1 ls")
		println("timer1 stop")
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a timer for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(timerFile); err == nil {
			fmt.Println("A timer is already running; Please stop the current timer before starting a new one.")
			return
		}
		projectName := args[0]
		startTime := time.Now().Unix()
		data := fmt.Sprintf("%s\n%d", projectName, startTime)
		err := os.WriteFile(timerFile, []byte(data), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Started timer for project: %s\n", projectName)
	},
}

func timerDetails() (*string, int, int) {
	data, err := os.ReadFile(timerFile)
	if err != nil {
		fmt.Println("No timer is currently running.")
		return nil, 0, 0
	}
	lines := strings.Split(string(data), "\n")
	projectName := lines[0]
	startTime, err := strconv.ParseInt(lines[1], 10, 64)
	if err != nil {
		fmt.Println("Error parssing timer info: ", err)
		fmt.Println("Stop current timer to remove the error")
		return nil, 0, 0
	}
	elapsed := time.Since(time.Unix(startTime, 0))
	elapsedMinutes := int(elapsed.Minutes())
	elapsedSeconds := int(elapsed.Seconds()) % 60 // remaining seconds after minutes

	return &projectName, elapsedMinutes, elapsedSeconds
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List currently running timer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, elapsedMinutes, elapsedSeconds := timerDetails()
		if projectName != nil {
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
		projectName, elapsedMinutes, elapsedSeconds := timerDetails()
		if projectName != nil {
			details := fmt.Sprintf("Project: %s, elapsed time: %d minutes %d seconds \n", *projectName, elapsedMinutes, elapsedSeconds)
			println(details)
		}
		os.Remove(timerFile)
		println("Succesfully stopped project:", *projectName)
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
