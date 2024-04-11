package main

import (
	"fmt"
	"os"
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
			println("started project: ", args[0])
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

		if projectName != nil {
			status := "Running"
			if isStopped {
				status = "Stopped"
			}
			title := fmt.Sprintf("Project: %s; Status: %s\n", *projectName, status)
			println(title)
			printSessions(sessions)

		} else {
			println("No timer is currently running.")
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop currently running timer",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, sessions, isStopped := readTimer()
		if projectName != nil {
			if isStopped {
				println("No Timer is running currently.")
				return
			}
			printSessions(sessions)
			println("\nStopped project:", *projectName)
		}

		// stop timer by adding stop time to last session
		sessions[len(sessions)-1][1] = time.Now().Unix()
		saveTimerData(*projectName, sessions)
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a timer for a project",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName, sessions, isStopped := readTimer()
		if projectName != nil {
			printSessions(sessions)
			if !isStopped {
				println("Timer is already running for project:", *projectName)
				return
			}
		}

		// restart timer by adding new session
		sessions = append(sessions, [2]int64{time.Now().Unix(), 0})
		saveTimerData(*projectName, sessions)
		println("\nRestarted project:", *projectName)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
}

func main() {
	rootCmd.Execute()
}
