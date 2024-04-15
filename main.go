package main

import (
	"fmt"
	"os"
	"strconv"
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

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a timer for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sessions := make([][2]int64, 1)
		sessions[0] = [2]int64{time.Now().Unix(), 0}

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

			saveTimerData(args[0], sessions)
			println("\nstarted project: ", args[0])
		} else {
			projectName, _, isStopped := readTimer()
			if !isStopped {
				fmt.Printf("A timer is already running for project: %s\n", *projectName)
				return
			} else {
				saveTimerData(args[0], sessions)
				println("started project: ", args[0])
			}
		}
	},
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List currently running timer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, sessions, isStopped := readTimer()

		if projectName == nil {
			println("No timer is currently running.")
			return
		}

		status := "Running"
		if isStopped {
			status = "Stopped"
		}
		title := fmt.Sprintf("\nProject: %s; Status: %s\n", *projectName, status)
		println(title)
		printSessions(sessions)
	},
}

var editCmd = &cobra.Command{
	Use:                "edit",
	Short:              "Edit time for the project",
	Args:               cobra.ExactArgs(1),
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		errMsg1 := "Please provide number of minutes to add(e.g. 31) or remove(e.g. -20) from last session."
		if len(args) < 1 {
			println(errMsg1)
			return
		}
		editInfo := args[0]
		if len(editInfo) < 2 {
			println(errMsg1)
			return
		}
		minutes, err := strconv.Atoi(editInfo)
		if err != nil {
			println(errMsg1)
			return
		}
		projectName, sessions, _ := readTimer()
		if projectName == nil {
			println("No timer is currently running.")
			return
		}
		if len(sessions) == 0 {
			println("No sessions found.")
			return
		}

		lastSession := sessions[len(sessions)-1]
		if minutes < 0 {
			// add negative minutes from endTime to remove time
			lastSession[1] += int64(minutes * 60)
		} else {
			// remove positive minutes from startTime to add time
			lastSession[0] -= int64(minutes * 60)
		}
		sessions[len(sessions)-1] = lastSession
		saveTimerData(*projectName, sessions)
		println("\nedited project:", *projectName, "\n")
		printSessions(sessions)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop currently running timer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, sessions, isStopped := readTimer()
		if projectName == nil || isStopped {
			println("No timer is currently running.")
			return
		}
		// stop timer by adding stop time to last session
		sessions[len(sessions)-1][1] = time.Now().Unix()
		saveTimerData(*projectName, sessions)
		println("\nstopped project:", *projectName, "\n")
		printSessions(sessions)
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a timer for a project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, sessions, isStopped := readTimer()
		if projectName == nil {
			println("No timer is currently running.")
			return
		}
		if !isStopped {
			println("Timer is already running for project:", *projectName)
			return
		}

		// restart timer by adding new session
		sessions = append(sessions, [2]int64{time.Now().Unix(), 0})
		saveTimerData(*projectName, sessions)
		println("\nrestarted project:", *projectName, "\n")
		printSessions(sessions)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
}

func main() {
	rootCmd.Execute()
}
