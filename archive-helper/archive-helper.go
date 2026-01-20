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

	// B. Cookie Logic (Preserved from previous step)
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

	// C. Execute Docker
	if err := executeDockerCommand(finalURL, finalCookiePath); err != nil {
		fmt.Printf("Execution failed: %v\n", err)
		os.Exit(1)
	}
}

// --- SYSTEM OPERATIONS ---

func checkDocker() error {
	// 1. Check if binary exists
	_, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker binary not found in PATH")
	}

	// 2. Check if daemon is responsive (and if we need sudo)
	// We run 'docker info'. If it works, great.
	// If it fails with "permission denied", we need sudo.
	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		// Simple heuristic: If it failed, assume it's a permission issue on Linux
		if runtime.GOOS == "linux" {
			fmt.Println(">> Docker requires privileges. Enabling sudo mode...")
			needsSudo = true
		} else {
			return fmt.Errorf("docker daemon is not running or accessible")
		}
	}
	return nil
}

func executeDockerCommand(url string, cookiePath string) error {
	// 1. Gather System Info for Variables
	currentUser, _ := user.Current()
	currentDir, _ := os.Getwd()

	// Default to standard "pwd" mount
	volumeMount := fmt.Sprintf("%s:/out", currentDir)

	// Default to current user's ID (Linux/Mac only)
	uid := currentUser.Uid
	gid := currentUser.Gid

	// 2. Build Arguments List
	// Start with the base command parts
	var cmdArgs []string

	if needsSudo {
		cmdArgs = append(cmdArgs, "docker")
	} else {
		// If no sudo, the executable is "docker", arguments start after
	}

	cmdArgs = append(cmdArgs, "run", "-it", "--rm")

	// 3. Conditional Logic: Cookies vs No Cookies
	if cookiePath != "None" && cookiePath != "" {
		// --- COOKIE MODE ---
		// User requested specific logic: Mount specific folder if cookies are present.
		// NOTE: In a real scenario, you likely want to mount the cookies FILE itself.
		// For this example, we follow your request to change the volume mount logic.

		// If the user provided a relative path, make it absolute
		absCookiePath, _ := filepath.Abs(cookiePath)
		cookieDir := filepath.Dir(absCookiePath)

		// Overwrite defaults as per your prompt example
		// "If cookies... use sensible defaults (hardcoded or based on input)"
		volumeMount = fmt.Sprintf("%s:/out", cookieDir)
		uid = "1000"
		gid = "1000"

		// We also technically need to pass the cookies to the container.
		// Assuming the container takes a flag, we add a volume for the file:
		cmdArgs = append(cmdArgs, "-v", fmt.Sprintf("%s:/app/cookies.txt", absCookiePath))
	}

	// 4. Add Common Flags
	cmdArgs = append(cmdArgs,
		"-v", volumeMount,
		"-e", fmt.Sprintf("MY_UID=%s", uid),
		"-e", fmt.Sprintf("MY_GID=%s", gid),
		"zeppelinsforever/livestream_dl_containerized:latest",
		"--log-level", "DEBUG",
		"--wait-for-video", "60",
		"--live-chat",
		"--resolution", "best",
	)

	// 5. Add Cookie Flag to the IMAGE (if needed)
	if cookiePath != "None" && cookiePath != "" {
		cmdArgs = append(cmdArgs, "--cookies", "/app/cookies.txt")
	}

	// 6. Add URL (Final Argument)
	cmdArgs = append(cmdArgs, url)

	// 7. Execution
	var finalCmd *exec.Cmd
	if needsSudo {
		// Run: sudo docker run ...
		finalCmd = exec.Command("sudo", cmdArgs...)
	} else {
		// Run: docker run ...
		finalCmd = exec.Command("docker", cmdArgs...)
	}

	// Connect input/output so the user sees the docker progress bar
	finalCmd.Stdout = os.Stdout
	finalCmd.Stderr = os.Stderr
	finalCmd.Stdin = os.Stdin

	fmt.Println("\n>> Launching Docker Container...")
	fmt.Printf("Command: %s\n\n", finalCmd.String())

	return finalCmd.Run()
}

// --- HELPERS (Same as before) ---
func getUserInput(r *bufio.Reader) string {
	input, _ := r.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func promptAndSaveNewPath(r *bufio.Reader) string {
	fmt.Print("Enter path to cookies.txt: ")
	input, _ := r.ReadString('\n')
	newPath := strings.TrimSpace(input)
	fmt.Print("Save this path as default for future runs? (y/n): ")
	if getUserInput(r) == "y" {
		saveEnvVar("LIVESTREAM_DL_CONTAINERIZED_COOKIES", newPath)
	}
	return newPath
}

func saveEnvVar(key, value string) {
	// (Use the "Search and Replace" version from the previous step here)
	// Abbreviated for space in this view, but include the full logic.
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
