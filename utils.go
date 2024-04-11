package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func writeToFile(data string) {
	err := os.WriteFile(timerFile, []byte(data), FileMode)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}
}

func saveTimerData(projectName string, sessionsArray [][2]int64) {
	data := fmt.Sprintf("%s\n%d\n", projectName, len(sessionsArray))
	for _, session := range sessionsArray {
		data += fmt.Sprintf("%d\n%d\n", session[0], session[1])
	}
	writeToFile(data)
}

func printSessions(sessions [][2]int64) {
	totalTime := time.Duration(0)
	for idx, session := range sessions {
		startTime := time.Unix(session[0], 0)
		stopTime0 := session[1]
		stopTime := time.Unix(stopTime0, 0)
		if stopTime0 == 0 {
			stopTime = time.Unix(time.Now().Unix(), 0)
		}
		elapsed := stopTime.Sub(startTime)
		elapsedMinutes := int(elapsed.Minutes())
		elapsedSeconds := int(elapsed.Seconds()) % 60 // remaining seconds after minutes
		details := fmt.Sprintf("Sesssion %d: %d minutes %d seconds", idx+1, elapsedMinutes, elapsedSeconds)
		println(details)
		totalTime += elapsed
	}
	totalMinutes := int(totalTime.Minutes())
	totalSeconds := int(totalTime.Seconds()) % 60
	details := fmt.Sprintf("Total: %d minutes %d seconds", totalMinutes, totalSeconds)
	println(details)
}

func readTimer() (*string, [][2]int64, bool) {
	data, err := os.ReadFile(timerFile)
	if err != nil {
		fmt.Println("No timer is currently running.")
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	projectName := lines[0]
	sessionsCount, err := strconv.ParseInt(lines[1], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing timer info: %v\nStop current timer to remove the error\n", err)
		return nil, make([][2]int64, 0), true
	}

	// sessionsArray is an array with each element has start and stop time
	// each element has 0-index as start time and 1-index as stop time
	sessionsArray := make([][2]int64, sessionsCount)
	for i := 0; i < int(sessionsCount); i++ {
		sessionsArray[i][0], err = strconv.ParseInt(lines[2*i+2], 10, 64)
		if err != nil {
			fmt.Printf("Error parsing timer info: %v\nStop current timer to remove the error\n", err)
			return nil, make([][2]int64, 0), true
		}
		sessionsArray[i][1], err = strconv.ParseInt(lines[2*i+3], 10, 64)
		if err != nil {
			fmt.Printf("Error parsing timer info: %v\nStop current timer to remove the error\n", err)
			return nil, make([][2]int64, 0), true
		}
	}

	stopTimeOfLastSession := sessionsArray[len(sessionsArray)-1][1]
	isStopped := stopTimeOfLastSession != 0

	return &projectName, sessionsArray, isStopped
}
