package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Config holds all configuration variables
type Config struct {
	Reader *bufio.Reader
	// package manager
	PackageManager    string
	PackageManagerCmd *exec.Cmd
	// frontend commands
	InstallFrontendCmd       *exec.Cmd
	InstallDependenciesCmd   *exec.Cmd
	InstallTailwindCmd       *exec.Cmd
	InstallTailwindInitCmd   *exec.Cmd
	InstallDaisyUICmd        *exec.Cmd
	InstallReactRouterDomCmd *exec.Cmd
	InstallAxiosCmd          *exec.Cmd
	RunFrontendServerCmd     *exec.Cmd
	// backend commands
	CreateBackendCmd   *exec.Cmd
	InstallExpressCmd  *exec.Cmd
	InstallClerkCmd    *exec.Cmd
	InstallFirebaseCmd *exec.Cmd
	InstallPrismaCmd   *exec.Cmd
	InitPrismaCmd      *exec.Cmd
	// general
	ChangeDirectoryIntoCmd  *exec.Cmd
	ChangeDirectoryOutOfCmd *exec.Cmd
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

func runInteractiveCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return "", err
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

func ChangeDirectoryInto(folder string) error {
	fmt.Printf("Changing directory to: %s\n", folder)
	err := os.Chdir(folder)
	if err != nil {
		return fmt.Errorf("failed to change directory: %v", err)
	}
	fmt.Printf("Successfully changed directory to: %s\n", folder)
	return nil
}

func ChangeDirectoryUp() error {
	fmt.Println("Moving up one directory...")
	err := os.Chdir("..")
	if err != nil {
		return fmt.Errorf("failed to move up one directory")
	}
	return nil
}

func MakeDirectory(folderName string) error {
	// Create the directory
	err := os.Mkdir(folderName, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func createFrontendProject() {
	fmt.Println("\nOkay, let's set up the frontend...")
	fmt.Println("What do you want to name this project?")
	projectName, _ := AppConfig.Reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	// decide which command to run depending on the users package manager selection
	switch AppConfig.PackageManager {
	case "npm":
		AppConfig.InstallFrontendCmd = exec.Command("npm", "create", "vite@latest", projectName)
	case "bun":
		AppConfig.InstallFrontendCmd = exec.Command("bun", "create", "vite", projectName)
	default:
		fmt.Println("Invalid package manager. Please choose 'npm' or 'bun'.")
		os.Exit(1)
	}

	// Run the interactive command and capture output
	_, err := runInteractiveCommand(AppConfig.InstallFrontendCmd.Path, AppConfig.InstallFrontendCmd.Args[1:]...)
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	ChangeDirectoryInto(projectName)

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

		// launch a new terminal window and run the server
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

func createBackendProject() {
	fmt.Println("Okay, let's set up the backend...")

	if askYesNo("Do you want to create a new folder for the backend? (no will deploy in current folder)") {
		fmt.Print("Name of the folder: ")
		folderName, _ := AppConfig.Reader.ReadString('\n')
		folderName = strings.TrimSpace(folderName)

		MakeDirectory(folderName)
		ChangeDirectoryInto(folderName)
	}

	if AppConfig.PackageManager == "npm" {
		AppConfig.CreateBackendCmd = exec.Command("npm", "init", "-y")
	} else {
		AppConfig.CreateBackendCmd = exec.Command("bun", "init", "-y")
	}

	_, err := runInteractiveCommand(AppConfig.CreateBackendCmd.Path, AppConfig.CreateBackendCmd.Args[1:]...)
	if err != nil {
		fmt.Printf("Error creating backend project: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Backend project created successfully!\n")
}

func installExpress() {
	fmt.Println("Installing Express and related packages...")
	if AppConfig.PackageManager == "npm" {
		AppConfig.InstallExpressCmd = exec.Command("npm", "install", "-D", "express", "@types/express", "cors", "@types/cors", "dotenv")
	} else {
		AppConfig.InstallExpressCmd = exec.Command("bun", "add", "-D", "express", "@types/express", "cors", "@types/cors", "dotenv")
	}

	_, err := runInteractiveCommand(AppConfig.InstallExpressCmd.Path, AppConfig.InstallExpressCmd.Args[1:]...)
	if err != nil {
		fmt.Printf("Error installing Express and related packages: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Express and related packages installed successfully!")
}

func createDotEnvFile() {
	fmt.Println("Creating .env file...")
	envContent := `PORT=3000
DATABASE_URL=postgresql://postgres:postgres@localhost:10017
# Add other environment variables as needed
`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		fmt.Printf("Error creating .env file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(".env file created successfully!")
}

func installClerk() {
	if askYesNo("Do you want to install Clerk for authentication?") {
		fmt.Println("Installing Clerk...")
		if AppConfig.PackageManager == "npm" {
			AppConfig.InstallClerkCmd = exec.Command("npm", "install", "@clerk/clerk-sdk-node")
		} else {
			AppConfig.InstallClerkCmd = exec.Command("bun", "add", "@clerk/clerk-sdk-node")
		}

		_, err := runInteractiveCommand(AppConfig.InstallClerkCmd.Path, AppConfig.InstallClerkCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing Clerk: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Clerk installed successfully!")
	}
}

func installFirebase() {
	if askYesNo("Do you want to install Firebase?") {
		fmt.Println("Installing Firebase...")
		if AppConfig.PackageManager == "npm" {
			AppConfig.InstallFirebaseCmd = exec.Command("npm", "install", "firebase-admin")
		} else {
			AppConfig.InstallFirebaseCmd = exec.Command("bun", "add", "firebase-admin")
		}

		_, err := runInteractiveCommand(AppConfig.InstallFirebaseCmd.Path, AppConfig.InstallFirebaseCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing Firebase: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Firebase installed successfully!")
	}
}

// func setupAuth() {

// }

func createDatabase() {
	if askYesNo("Do you want to set up a PostgreSQL database with Docker?") {
		fmt.Println("Setting up PostgreSQL database...")

		gistURL := "https://gist.githubusercontent.com/kmankan/e2d9414a15a669af21840cce146e7201/raw/7a906c7d9a6d82a61e645486c30de8c7738791a4/docker-compose.yml"
		curlCmd := exec.Command("curl", "-s", gistURL)

		dockerComposeContent, err_get := curlCmd.Output()
		if err_get != nil {
			fmt.Printf("Error fetching docker-compose content: %v\n", err_get)
			os.Exit(1)
		}

		err_write := os.WriteFile("docker-compose.yml", []byte(dockerComposeContent), 0644)
		if err_write != nil {
			fmt.Printf("Error creating docker-compose.yml file: %v\n", err_write)
			os.Exit(1)
		}

		fmt.Println("docker-compose.yml file created. You can start the database by running 'docker compose up -d' in this directory.")
	}
}

func installPrisma() {
	if askYesNo("Do you want to install Prisma?") {
		fmt.Println("Installing Prisma...")
		// install prisma
		if AppConfig.PackageManager == "npm" {
			AppConfig.InstallPrismaCmd = exec.Command("npm", "install", "-D", "prisma", "@prisma/client")
		} else {
			AppConfig.InstallPrismaCmd = exec.Command("bun", "add", "-D", "prisma", "@prisma/client")
		}

		_, err := runInteractiveCommand(AppConfig.InstallPrismaCmd.Path, AppConfig.InstallPrismaCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error installing Prisma: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Initializing Prisma...")
		// init prisma
		if AppConfig.PackageManager == "npm" {
			AppConfig.InitPrismaCmd = exec.Command("npx", "prisma", "init")
		} else {
			AppConfig.InitPrismaCmd = exec.Command("bunx", "prisma", "init")
		}

		_, err = runInteractiveCommand(AppConfig.InitPrismaCmd.Path, AppConfig.InitPrismaCmd.Args[1:]...)
		if err != nil {
			fmt.Printf("Error initializing Prisma: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Prisma installed and initialized successfully!")
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

func deployBackend() {
	createBackendProject()
	installExpress()
	createDotEnvFile()
	installClerk()
	installFirebase()
	createDatabase()
	installPrisma()
	fmt.Println("\nBackend setup completed successfully!")

}

func deployFullstack() {
	deployBackend()
	ChangeDirectoryUp()
	deployFrontend()

}

func main() {
	fmt.Println("What do you want to set up?")
	fmt.Println("1. Frontend")
	fmt.Println("2. Backend")
	fmt.Println("3. Fullstack")
	choice := getUserInput("Enter your choice (1/2/3): ")

	switch choice {
	case "1":
		deployFrontend()
	case "2":
		deployBackend()
	case "3":
		deployFullstack()
	default:
		fmt.Println("Invalid choice. Exiting.")
		os.Exit(1)
	}
}
