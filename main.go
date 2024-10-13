package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Config holds all configuration variables
type Config struct {
	Reader                   *bufio.Reader
	PackageManager           string
	PackageManagerCmd        *exec.Cmd
	InstallFrontendCmd       *exec.Cmd
	InstallDependenciesCmd   *exec.Cmd
	InstallTailwindCmd       *exec.Cmd
	InstallTailwindInitCmd   *exec.Cmd
	InstallDaisyUICmd        *exec.Cmd
	InstallReactRouterDomCmd *exec.Cmd
	InstallAxiosCmd          *exec.Cmd
	RunFrontendServerCmd     *exec.Cmd
}

// Global configuration instance
var AppConfig Config

// This special function runs before main().
func init() {
	// initialize the reader
	AppConfig.Reader = bufio.NewReader(os.Stdin)

	//decide which package manager to use
	fmt.Print("Do you want to use npm or bun? (npm/bun): ")
	choosePackageManager, _ := AppConfig.Reader.ReadString('\n')
	AppConfig.PackageManager = strings.TrimSpace(strings.ToLower(choosePackageManager))
}

// function for running command line commands
func runInteractiveCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var outputBuffer strings.Builder
	cmd.Stdout = io.MultiWriter(os.Stdout, &outputBuffer)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return outputBuffer.String(), err
}

// function for asking the user Yes or No
// it follows the following pattern:
// installTailwind := askYesNo(reader, "Do you want to install Tailwind CSS?")
// if (installTailwind {})
func askYesNo(question string) bool {
	fmt.Print(question + " (y/n): ")
	answer, _ := AppConfig.Reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(answer)) == "y"
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := AppConfig.Reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}
	return dir
}

func createFrontendProject() {
	fmt.Println("Okay, let's set up the frontend...")

	var projectName string

	// decide which command to run depending on the users package manager selection
	switch AppConfig.PackageManager {
	case "npm":
		AppConfig.InstallFrontendCmd = exec.Command("npm", "create", "vite@latest")
	case "bun":
		AppConfig.InstallFrontendCmd = exec.Command("bun", "create", "vite")
	default:
		fmt.Println("Invalid package manager. Please choose 'npm' or 'bun'.")
		os.Exit(1)
	}

	// Run the interactive command and capture output
	output, err := runInteractiveCommand(AppConfig.InstallFrontendCmd.Path, AppConfig.InstallFrontendCmd.Args[1:]...)
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	// Search for the "cd" instruction in the output
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "  cd ") {
			projectName = strings.TrimSpace(strings.TrimPrefix(line, "  cd "))
			break
		}
	}

	if projectName == "" {
		fmt.Println("Failed to capture project name.")
		os.Exit(1)
	}

	fmt.Printf("Project name captured: %s\n", projectName)

	err = os.Chdir(projectName)
	if err != nil {
		fmt.Printf("Error changing to project directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully changed to project directory: %s\n", projectName)
}

func installDependencies() {
	fmt.Println("Installing dependencies...")

	if AppConfig.PackageManager == "npm" {
		AppConfig.InstallDependenciesCmd = exec.Command("npm", "install")
	} else {
		AppConfig.InstallDependenciesCmd = exec.Command("bun", "install")
	}

	_, err := runInteractiveCommand(AppConfig.InstallDependenciesCmd.Path, AppConfig.InstallDependenciesCmd.Args[1:]...)
	if err != nil {
		fmt.Printf("Error installing dependencies: %v\n", err)
		os.Exit(1)
	}
}

func installTailwind() {
	installTailwind := askYesNo("Do you want to install Tailwind CSS?")

	if installTailwind {
		fmt.Println("Installing Tailwind CSS...")

		if AppConfig.PackageManager == "npm" {
			AppConfig.InstallTailwindCmd = exec.Command("npm", "install", "-D", "tailwindcss", "postcss", "autoprefixer")
			AppConfig.InstallTailwindInitCmd = exec.Command("npx", "tailwindcss", "init", "-p")
			// initCmd = exec.Command("npx", "tailwindcss", "init", "-p")
		} else {
			AppConfig.InstallTailwindCmd = exec.Command("bun", "add", "-D", "tailwindcss", "postcss", "autoprefixer")
			AppConfig.InstallTailwindInitCmd = exec.Command("bunx", "tailwindcss", "init", "-p")
			// initCmd = exec.Command("bunx", "tailwindcss", "init", "-p")
		}

		_, err := runInteractiveCommand(AppConfig.InstallTailwindCmd.Path, AppConfig.InstallTailwindCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing Tailwind CSS: %v\n", err)
			os.Exit(1)
		}

		_, err2 := runInteractiveCommand(AppConfig.InstallTailwindInitCmd.Path, AppConfig.InstallTailwindInitCmd.Args[1:]...)
		if err2 != nil {
			fmt.Printf("Error installing Tailwind CSS: %v\n", err2)
			os.Exit(1)
		}
	}
}

func installDaisyUI() {
	installDaisyUI := askYesNo("Do you want to install DaisyUI?")

	tailwindConfigWithDaisyUI := `/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
}`

	if installDaisyUI {
		fmt.Println("Installing DaisyUI...")

		if AppConfig.PackageManager == "npm" {
			AppConfig.InstallDaisyUICmd = exec.Command("npm", "install", "-D", "daisyui@latest")
		} else {
			AppConfig.InstallDaisyUICmd = exec.Command("bun", "add", "-D", "daisyui@latest")
		}

		// running installation command
		_, err1 := runInteractiveCommand(AppConfig.InstallDaisyUICmd.Path, AppConfig.InstallDaisyUICmd.Args[1:]...)
		if err1 != nil {
			fmt.Printf("Error installing DaisyUI: %v\n", err1)
			os.Exit(1)
		}

		// writing to tailwind config file with DaisyUI updates
		err2 := os.WriteFile("tailwind.config.js", []byte(tailwindConfigWithDaisyUI), 0644)
		if err2 != nil {
			fmt.Printf("Error updating tailwind.config.js: %v\n", err2)
			os.Exit(1)
		}

		indexCSS := `@tailwind base;
	@tailwind components;
	@tailwind utilities;`
		err3 := os.WriteFile("src/index.css", []byte(indexCSS), 0644)
		if err3 != nil {
			fmt.Printf("Error updating src/index.css: %v\n", err3)
			os.Exit(1)
		}

		fmt.Println("Tailwind CSS installed and configured successfully!")
	}
}

func installReactRouterDom() {
	installReactRouterDom := askYesNo("Do you want to install react-router-dom?")
	// Ask the user if they wan to install react-router-dom
	if installReactRouterDom {
		fmt.Println("Installing react-router-dom...")
		if AppConfig.PackageManager == "npm" {
			AppConfig.InstallReactRouterDomCmd = exec.Command("npm", "install", "react-router-dom")
		} else {
			AppConfig.InstallReactRouterDomCmd = exec.Command("bun", "add", "react-router-dom")
		}

		// install react-router-dom
		_, err := runInteractiveCommand(AppConfig.InstallReactRouterDomCmd.Path, AppConfig.InstallReactRouterDomCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing react-router-dom: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("react-router-dom installed successfully!")
	}
}

func installAxios() {
	installAxios := askYesNo("Do you want to install axios?")
	// Ask the user if they wan to install axios
	if installAxios {
		fmt.Println("Installing axios...")

		if AppConfig.PackageManager == "npm" {
			AppConfig.InstallAxiosCmd = exec.Command("npm", "install", "axios")
		} else {
			AppConfig.InstallAxiosCmd = exec.Command("bun", "add", "axios")
		}

		// install axios
		_, err := runInteractiveCommand(AppConfig.InstallAxiosCmd.Path, AppConfig.InstallAxiosCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing axios: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("axios installed successfully!")
	}
}

func RunFrontendServer() {
	// Ask the user if they want to run a development server now
	if askYesNo("\nDo you want to start the development server now?") {
		fmt.Println("Creating server...")
		currentDir := getCurrentDir()

		if AppConfig.PackageManager == "npm" {
			AppConfig.RunFrontendServerCmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell app "Terminal" to do script "cd '%s' && npm run dev"`, currentDir))
		} else {
			AppConfig.RunFrontendServerCmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell app "Terminal" to do script "cd '%s' && bun run dev"`, currentDir))
		}

		// Building the development server
		_, err := runInteractiveCommand(AppConfig.RunFrontendServerCmd.Path, AppConfig.RunFrontendServerCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error starting server in new window: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Server started in a new terminal window.")
	}
}

func deployFrontend() {
	createFrontendProject()
	installDependencies()
	installTailwind()
	installDaisyUI()
	installReactRouterDom()
	installAxios()
	fmt.Println("\nReact Vite project created successfully!")
	RunFrontendServer()
}

func installExpress() {

}

func deployBackend() {
	fmt.Print("Okay let's setup the backend...\n")

}

// EXPRESSJS
// install dev dependencies: express, @types/express, cors, @types/cors, dotenv
// AUTHENTICATION:
// clerk or firebase
// DATABASE
// generate docker yaml file
// install prisma (prisma, @prisma/client)
// prisma init
// generate prisma

// func deployFullstack() {
// 	deployFrontend()
// 	deployBackend()
// }

func main() {

	// Do you want to setup: backend, frontend or fullstack?
	// if frontend --> fn(frontend),
	// if backend --> fn(backend),
	// if fullstack --> fn(frontend) + fn(backend)
	deployFrontend()
}
