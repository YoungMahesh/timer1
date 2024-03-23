package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "app is a simple CLI application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World!")
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.Join(args, " ")
		f, err := os.OpenFile("projects.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		if _, err := f.WriteString(projectName + "\n"); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Added project: ", projectName)
	},
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile("projects.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(data))
	},
}

var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a project name to remove")

		}
		projectName := args[0]
		data, err := os.ReadFile("projects.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		projects := strings.Split(string(data), "\n")
		for i, project := range projects {
			if project == projectName {
				projects = append(projects[:i], projects[i+1:]...)
				break
			}
		}

		output := strings.Join(projects, "\n")
		err = os.WriteFile("projects.txt", []byte(output), 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Removed project: ", projectName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(removeCmd)
}

func main() {
	rootCmd.Execute()
}
