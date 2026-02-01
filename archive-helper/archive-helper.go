package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time" // [ADDED] Needed for generating the date string on Windows

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

// CONFIGURATION
type Config struct {
	CookiesPath string `envconfig:"COOKIES"`
}

var (
	cookiePathFlag string
	silentFlag     bool
	appConfig      Config
	// We determine these at runtime
	needsSudo bool
)

func main() {
	// 0. SYSTEM CHECKS (Run immediately on startup)
	if err := checkDocker(); err != nil {
		fmt.Println("CRITICAL ERROR: Docker is not working.")
		fmt.Println(err)
		fmt.Println("Please ensure Docker is installed and running.")
		os.Exit(1)
	}

	// 1. SETUP PERSISTENCE
	_ = godotenv.Load(".env")
	envconfig.Process("LIVESTREAM_DL_CONTAINERIZED", &appConfig)

	// 2. DEFINE COMMAND
	var rootCmd = &cobra.Command{
		Use:   "archive-helper [URL]",
		Short: "Docker wrapper for livestream_dl",
		Args:  cobra.MaximumNArgs(1),
		Run:   runApplicationLogic,
	}

	rootCmd.Flags().StringVar(&cookiePathFlag, "cookies", "", "Path to cookies.txt")
	rootCmd.Flags().BoolVar(&silentFlag, "silent", false, "Run silently in background")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// --- CORE LOGIC ---

func runApplicationLogic(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)
	var finalURL string
	var finalCookiePath string

	// A. URL Handling
	if len(args) > 0 {
		finalURL = args[0]
	} else {
		fmt.Print("Please enter the URL: ")
		input, _ := reader.ReadString('\n')
		finalURL = strings.TrimSpace(input)
	}

	// B. Cookie Logic
	if cmd.Flags().Changed("cookies") {
		finalCookiePath = cookiePathFlag
		if appConfig.CookiesPath != "" && finalCookiePath != appConfig.CookiesPath {
			fmt.Printf("\nNotice: You used flag '%s', but saved default is '%s'.\n", finalCookiePath, appConfig.CookiesPath)
			fmt.Print("Update saved default to match this flag? (y/n): ")
			if getUserInput(reader) == "y" {
				saveEnvVar("LIVESTREAM_DL_CONTAINERIZED_COOKIES", finalCookiePath)
			}
		}
	} else {
		if appConfig.CookiesPath != "" {
			fmt.Printf("\nSaved cookies path detected: %s\n", appConfig.CookiesPath)
			fmt.Print("Do you want to use this path? (y/n/x): ")
			choice := getUserInput(reader)
			if choice == "y" {
				finalCookiePath = appConfig.CookiesPath
			} else if choice == "n" {
				finalCookiePath = promptAndSaveNewPath(reader)
			} else {
				finalCookiePath = "None"
			}
		} else {
			fmt.Print("\nDo you want to provide a cookies file? (y/n): ")
			if getUserInput(reader) == "y" {
				finalCookiePath = promptAndSaveNewPath(reader)
			} else {
				finalCookiePath = "None"
			}
		}
	}
	// Ensure the image is up to date
	if err := pullDockerImage(); err != nil {
		fmt.Printf("Warning: Could not pull latest image: %v\n", err)
		fmt.Println("Attempting to run with local version...")
	}
	// C. Execute Docker
	if err := executeDockerCommand(finalURL, finalCookiePath, silentFlag); err != nil {
		fmt.Printf("Execution failed: %v\n", err)
		os.Exit(1)
	}
}

// --- SYSTEM OPERATIONS ---

func checkDocker() error {
	_, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker binary not found in PATH")
	}

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		if runtime.GOOS == "linux" {
			fmt.Println(">> Docker requires privileges. Enabling sudo mode...")
			needsSudo = true
		} else {
			return fmt.Errorf("docker daemon is not running or accessible")
		}
	}
	return nil
}

func executeDockerCommand(url string, cookiePath string, silent bool) error {
	// 1. Gather System Info
	currentUser, _ := user.Current()
	currentDir, _ := os.Getwd()

	// BASELINE: Always mount the CURRENT directory to /out
	volumeMount := fmt.Sprintf("%s:/out", currentDir)

	// Defaults for "No Cookie" mode
	uid := currentUser.Uid
	gid := currentUser.Gid

	// 2. Build Arguments
	var cmdArgs []string

	if needsSudo {
		cmdArgs = append(cmdArgs, "docker")
	}

	cmdArgs = append(cmdArgs, "run", "--rm")

	// 3. Conditional Logic
	if cookiePath != "None" && cookiePath != "" {
		uid = "1000"
		gid = "1000"

		absCookiePath, err := filepath.Abs(cookiePath)
		if err != nil {
			return fmt.Errorf("failed to resolve cookie path: %v", err)
		}
		cmdArgs = append(cmdArgs, "-v", fmt.Sprintf("%s:/app/cookies.txt", absCookiePath))
	}

	// 4. Add Common Flags
	cmdArgs = append(cmdArgs,
		"-v", volumeMount,
	)

	// [CHANGE] Only append UID/GID on non-Windows systems
	if runtime.GOOS != "windows" {
		cmdArgs = append(cmdArgs,
			"-e", fmt.Sprintf("MY_UID=%s", uid),
			"-e", fmt.Sprintf("MY_GID=%s", gid),
		)
	}

	cmdArgs = append(cmdArgs,
		"zeppelinsforever/livestream_dl_containerized:latest",
		"--log-level", "INFO",
		"--threads", "4",
		"--dash", 
		"--m3u8",
		"--wait-for-video", "60:600",
		"--write-thumbnail",
		"--embed-thumbnail",
		"--clean-info-json", 
		"--remove-ip-from-json",
		"--live-chat",
		"--resolution", "best",
	)

	// 5. Add Cookie Flag to the IMAGE
	if cookiePath != "None" && cookiePath != "" {
		cmdArgs = append(cmdArgs, "--cookies", "/app/cookies.txt")
	}

	// 6. Add URL
	cmdArgs = append(cmdArgs, url)

	// 7. Execution
	var finalCmd *exec.Cmd

	if silent {
		// [CHANGE] Branching Logic for Windows vs Linux Silent Mode
		if runtime.GOOS == "windows" {
			// --- WINDOWS SILENT MODE ---
			// Windows does not have "nohup". We must manually create the file and direct output.

			// 1. Create the log file name: nohup.YYYY-MM-DD.out
			dateStr := time.Now().Format("2006-01-02")
			fileName := fmt.Sprintf("nohup.%s.out", dateStr)

			// 2. Open/Create the file
			outfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return fmt.Errorf("failed to create log file for silent mode: %v", err)
			}
			// Note: We deliberately do not defer outfile.Close() here because we want
			// it to stay open for the duration of the child process if possible,
			// though Go's exec.Start handles file descriptor inheritance.

			// 3. Setup the command
			finalCmd = exec.Command("docker", cmdArgs...)
			finalCmd.Stdout = outfile
			finalCmd.Stderr = outfile

			fmt.Println("\n>> Launching Docker Container in Background (Windows Mode)...")
			fmt.Printf("Logs redirected to: %s\n", fileName)

			// 4. Start asynchronously (detach)
			return finalCmd.Start()
		} else {
			// --- LINUX/MAC SILENT MODE (Original Logic) ---
			binary := "docker"
			if needsSudo {
				binary = "sudo"
			}

			quotedArgs := make([]string, len(cmdArgs))
			for i, arg := range cmdArgs {
				quotedArgs[i] = fmt.Sprintf("%q", arg)
			}

			fullCmdStr := fmt.Sprintf("%s %s", binary, strings.Join(quotedArgs, " "))
			shellCmd := fmt.Sprintf("nohup %s > nohup.$(date +\"%%F\").out 2>&1 &", fullCmdStr)

			fmt.Println("\n>> Launching Docker Container in Background (Silent Mode)...")
			fmt.Printf("Command: %s\n", shellCmd)

			finalCmd = exec.Command("sh", "-c", shellCmd)
			return finalCmd.Run()
		}
	} else {
		// -- INTERACTIVE MODE --
		if needsSudo {
			finalCmd = exec.Command("sudo", cmdArgs...)
		} else {
			finalCmd = exec.Command("docker", cmdArgs...)
		}

		finalCmd.Stdout = os.Stdout
		finalCmd.Stderr = os.Stderr
		finalCmd.Stdin = os.Stdin

		fmt.Println("\n>> Launching Docker Container...")
		displayCmd := finalCmd.String()
		fmt.Printf("Command: %s\n\n", displayCmd)

		return finalCmd.Run()
	}
}

// --- HELPERS ---

func getUserInput(r *bufio.Reader) string {
	input, _ := r.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func promptAndSaveNewPath(r *bufio.Reader) string {
	fmt.Print("Enter path to cookies.txt (no perentheses): ")
	input, _ := r.ReadString('\n')
	newPath := strings.TrimSpace(input)
	fmt.Print("Save this path as default for future runs? (y/n): ")
	if getUserInput(r) == "y" {
		saveEnvVar("LIVESTREAM_DL_CONTAINERIZED_COOKIES", newPath)
	}
	return newPath
}

func pullDockerImage() error {
	fmt.Println(">> Checking for updates for zeppelinsforever/livestream_dl_containerized...")

	imageName := "zeppelinsforever/livestream_dl_containerized:latest"
	var pullCmd *exec.Cmd

	if needsSudo {
		pullCmd = exec.Command("sudo", "docker", "pull", imageName)
	} else {
		pullCmd = exec.Command("docker", "pull", imageName)
	}

	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr

	return pullCmd.Run()
}

func saveEnvVar(key, value string) {
	envFile := ".env"
	newEntry := fmt.Sprintf("%s=%s", key, value)
	input, _ := os.ReadFile(envFile)
	lines := strings.Split(string(input), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = newEntry
			found = true
			break
		}
	}
	if !found {
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines[len(lines)-1] = newEntry
		} else {
			lines = append(lines, newEntry)
		}
	}
	output := strings.Join(lines, "\n")
	_ = os.WriteFile(envFile, []byte(output), 0644)
}
