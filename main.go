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

func askYesNo(reader *bufio.Reader, question string) bool {
	fmt.Print(question + " (y/n): ")
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(answer)) == "y"
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Okay let's setup the frontend...\nWhat do you want to name this folder? ")
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

	installTailwind := askYesNo(reader, "Do you want to install Tailwind CSS?")

	if installTailwind {
		fmt.Println("Installing Tailwind CSS...")
		var tailwindCmd *exec.Cmd
		var initCmd *exec.Cmd
		if packageManager == "npm" {
			tailwindCmd = exec.Command("npm", "install", "-D", "tailwindcss", "postcss", "autoprefixer")
			initCmd = exec.Command("npx", "tailwindcss", "init", "-p")
		} else {
			tailwindCmd = exec.Command("bun", "add", "-D", "tailwindcss", "postcss", "autoprefixer")
			initCmd = exec.Command("bunx", "tailwindcss", "init", "-p")
		}

		err = runInteractiveCommand(tailwindCmd.Path, tailwindCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing Tailwind CSS: %v\n", err)
			os.Exit(1)
		}

		err = runInteractiveCommand(initCmd.Path, initCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error initializing Tailwind CSS: %v\n", err)
			os.Exit(1)
		}

		installDaisyUI := askYesNo(reader, "Do you want to install DaisyUI?")

		tailwindConfig := `/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [`

		if installDaisyUI {
			fmt.Println("Installing DaisyUI...")
			var daisyUICmd *exec.Cmd
			if packageManager == "npm" {
				daisyUICmd = exec.Command("npm", "install", "-D", "daisyui@latest")
			} else {
				daisyUICmd = exec.Command("bun", "add", "-D", "daisyui@latest")
			}

			err = runInteractiveCommand(daisyUICmd.Path, daisyUICmd.Args[1:]...)
			if err != nil {
				fmt.Printf("Error installing DaisyUI: %v\n", err)
				os.Exit(1)
			}

			tailwindConfig += `
    require('daisyui'),`
		}

		tailwindConfig += `
  ],
}`

		err = os.WriteFile("tailwind.config.js", []byte(tailwindConfig), 0644)
		if err != nil {
			fmt.Printf("Error updating tailwind.config.js: %v\n", err)
			os.Exit(1)
		}

		indexCSS := `@tailwind base;
@tailwind components;
@tailwind utilities;`
		err = os.WriteFile("src/index.css", []byte(indexCSS), 0644)
		if err != nil {
			fmt.Printf("Error updating src/index.css: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Tailwind CSS installed and configured successfully!")
	}

	// Ask the user if they wan to install react-router-dom
	if askYesNo(reader, "Do you want to install react-router-dom?") {
		fmt.Println("Installing react-router-dom...")
		var routerCmd *exec.Cmd
		if packageManager == "npm" {
			routerCmd = exec.Command("npm", "install", "react-router-dom")
		} else {
			routerCmd = exec.Command("bun", "add", "react-router-dom")
		}

		// install react-router-dom
		err = runInteractiveCommand(routerCmd.Path, routerCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing react-router-dom: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("react-router-dom installed successfully!")
	}

	// Ask the user if they wan to install axios
	if askYesNo(reader, "Do you want to install axios?") {
		fmt.Println("Installing axios...")
		var axiosCmd *exec.Cmd
		if packageManager == "npm" {
			axiosCmd = exec.Command("npm", "install", "axios")
		} else {
			axiosCmd = exec.Command("bun", "add", "axios")
		}

		// install axios
		err = runInteractiveCommand(axiosCmd.Path, axiosCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing axios: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("axios installed successfully!")
	}

	fmt.Println("\nReact Vite project created successfully!")

	// Ask the user if they want to run a development server now
	if askYesNo(reader, "\nDo you want to start the development server now?") {
		fmt.Println("Creating server...")
		var runServerCmd *exec.Cmd
		if packageManager == "npm" {
			runServerCmd = exec.Command("npm", "run", "dev")
		} else {
			runServerCmd = exec.Command("bun", "run", "dev")
		}

		// Building the development server
		err = runInteractiveCommand(runServerCmd.Path, runServerCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing axios: %v\n", err)
			os.Exit(1)
		}
		//fmt.Println("axios installed successfully!")
	}

}
