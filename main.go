package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runInteractiveCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Okay let's setup the frontend. What do you want to name this folder? ")
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	fmt.Print("Do you want to use npm or bun? (npm/bun): ")
	packageManager, _ := reader.ReadString('\n')
	packageManager = strings.TrimSpace(strings.ToLower(packageManager))

	var cmd *exec.Cmd
	switch packageManager {
	case "npm":
		cmd = exec.Command("npm", "create", "vite@latest", projectName)
	case "bun":
		cmd = exec.Command("bun", "create", "vite", projectName)
	default:
		fmt.Println("Invalid package manager. Please choose 'npm' or 'bun'.")
		os.Exit(1)
	}

	err := runInteractiveCommand(cmd.Path, cmd.Args[1:]...)
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	err = os.Chdir(projectName)
	if err != nil {
		fmt.Printf("Error changing to project directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Installing dependencies...")
	var installCmd *exec.Cmd
	if packageManager == "npm" {
		installCmd = exec.Command("npm", "install")
	} else {
		installCmd = exec.Command("bun", "install")
	}

	err = runInteractiveCommand(installCmd.Path, installCmd.Args[1:]...)
	if err != nil {
		fmt.Printf("Error installing dependencies: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("React Vite project created successfully!")
	fmt.Println("To start your development server, run:")
	if packageManager == "npm" {
		fmt.Printf("cd %s && npm run dev\n", projectName)
	} else {
		fmt.Printf("cd %s && bun run dev\n", projectName)
	}
}
