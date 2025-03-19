package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Installs dependencies and configures Brightside-Go",
	Run: func(cmd *cobra.Command, args []string) {
		runSetup()
	},
}

// 🚀 Run Full Setup
func runSetup() {
	fmt.Println("🚀 Running Brightside-Go Setup...\n")

	switch runtime.GOOS {
	case "darwin":
		fmt.Println("🍏 macOS detected!")
		installMacDependencies()
	case "linux":
		fmt.Println("🐧 Linux detected!")
		installLinuxDependencies()
	default:
		fmt.Println("❌ Unsupported OS.")
		os.Exit(1)
	}

	// Move binary & update environment
	moveBinary()
	configureShell()

	fmt.Println("\n✅ Setup Complete! Run 'brightside --help' to get started.")
}

// 🛠 Install Dependencies (macOS)
func installMacDependencies() {
	fmt.Println("🔹 Checking for Homebrew...")
	if !commandExists("brew") {
		fmt.Println("🍺 Homebrew not found! Installing...")
		exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)").Run()
	}

	packages := []string{"git", "yt-dlp", "ffmpeg", "wget", "zsh"}
	for _, pkg := range packages {
		fmt.Printf("🔹 Installing %s...\n", pkg)
		exec.Command("brew", "install", pkg).Run()
	}

	installP10K() // Install Powerlevel10k & plugins AFTER installing Zsh
}

// 🛠 Install Dependencies (Linux)
// 🛠 Install Dependencies (Linux)
func installLinuxDependencies() {
	fmt.Println("🔹 Checking for APT...")
	if !commandExists("apt") {
		fmt.Println("❌ APT package manager not found! Make sure you're on a Debian-based system.")
		os.Exit(1)
	}

	packages := []string{"git", "yt-dlp", "ffmpeg", "wget", "zsh"}
	for _, pkg := range packages {
		fmt.Printf("🔹 Installing %s...\n", pkg)
		exec.Command("sudo", "apt", "install", "-y", pkg).Run()
	}

	installP10K() // Install Powerlevel10k & plugins AFTER installing Zsh
}

// 🚚 Move Binary to `/usr/local/bin`
func moveBinary() {
	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Println("❌ Error detecting binary location:", err)
		return
	}

	targetPath := "/usr/local/bin/brightside"
	fmt.Println("🚚 Moving Brightside-Go to", targetPath)

	err = exec.Command("sudo", "mv", binaryPath, targetPath).Run()
	if err != nil {
		fmt.Println("❌ Failed to move binary:", err)
	} else {
		fmt.Println("✅ Brightside-Go installed globally!")
	}
}

// 🔧 Configure Shell (Import `.zshrc` and `p10k.zsh`)
// 🔧 Configure Shell (Import `.zshrc`, `.p10k.zsh`, and install plugins)
func configureShell() {
	shellConfig := getShellConfig()
	if shellConfig == "" {
		fmt.Println("❌ Could not detect shell configuration file.")
		return
	}

	// Ensure Brightside is in PATH
	fmt.Println("🔧 Adding Brightside-Go to PATH in", shellConfig)
	exportCmd := "export PATH=\"/usr/local/bin:$PATH\""
	appendToFile(shellConfig, exportCmd)

	// 🛠 Install Powerlevel10k & Plugins
	installP10K()

	// Import Custom `.zshrc` and `p10k.zsh`
	fmt.Println("🛠 Importing custom .zshrc and Powerlevel10k config...")
	configDir := "/usr/local/brightside-go/config"
	os.MkdirAll(configDir, os.ModePerm)

	copyFile("config/.zshrc", os.Getenv("HOME")+"/.zshrc")
	copyFile("config/.p10k.zsh", os.Getenv("HOME")+"/.p10k.zsh")

	// Ensure Zsh is the default shell
	if strings.Contains(os.Getenv("SHELL"), "zsh") {
		fmt.Println("✅ Zsh is already the default shell.")
	} else {
		fmt.Println("🛠 Setting Zsh as the default shell...")
		exec.Command("chsh", "-s", "/bin/zsh").Run()
	}

	fmt.Println("✅ Shell configuration complete! Run 'source ~/.zshrc'")
}

// 🔥 Install Oh My Zsh, Powerlevel10k, and Plugins
func installP10K() {
	fmt.Println("🎨 Checking Oh My Zsh & Powerlevel10k installation...")

	// 1️⃣ Ensure Oh My Zsh is Installed
	if _, err := os.Stat(os.Getenv("HOME") + "/.oh-my-zsh"); os.IsNotExist(err) {
		fmt.Println("⚡ Installing Oh My Zsh...")
		exec.Command("sh", "-c", "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)").Run()
	}

	// 2️⃣ Ensure `$ZSH_CUSTOM` is Set
	zshCustom := os.Getenv("HOME") + "/.oh-my-zsh/custom"
	os.Setenv("ZSH_CUSTOM", zshCustom)

	// 3️⃣ Install Powerlevel10k
	p10kPath := zshCustom + "/themes/powerlevel10k"
	if _, err := os.Stat(p10kPath); os.IsNotExist(err) {
		fmt.Println("🎨 Installing Powerlevel10k...")
		exec.Command("git", "clone", "--depth=1", "https://github.com/romkatv/powerlevel10k.git", p10kPath).Run()
	} else {
		fmt.Println("✅ Powerlevel10k is already installed.")
	}

	// 4️⃣ Install Zsh Plugins
	pluginDir := zshCustom + "/plugins"
	os.MkdirAll(pluginDir, os.ModePerm)

	// Autosuggestions
	autosuggestionsPath := pluginDir + "/zsh-autosuggestions"
	if _, err := os.Stat(autosuggestionsPath); os.IsNotExist(err) {
		fmt.Println("💡 Installing zsh-autosuggestions...")
		exec.Command("git", "clone", "https://github.com/zsh-users/zsh-autosuggestions", autosuggestionsPath).Run()
	}

	// Syntax Highlighting
	syntaxHighlightingPath := pluginDir + "/zsh-syntax-highlighting"
	if _, err := os.Stat(syntaxHighlightingPath); os.IsNotExist(err) {
		fmt.Println("💡 Installing zsh-syntax-highlighting...")
		exec.Command("git", "clone", "https://github.com/zsh-users/zsh-syntax-highlighting", syntaxHighlightingPath).Run()
	}
}

// 🔍 Get the correct shell configuration file (`.zshrc` or `.bashrc`)
func getShellConfig() string {
	shell := os.Getenv("SHELL")

	if strings.Contains(shell, "zsh") {
		return os.Getenv("HOME") + "/.zshrc"
	} else if strings.Contains(shell, "bash") {
		return os.Getenv("HOME") + "/.bashrc"
	} else {
		return ""
	}
}

// 📂 Copy File (Used for .zshrc and .p10k.zsh)
func copyFile(src, dest string) {
	input, err := os.ReadFile(src)
	if err != nil {
		fmt.Println("❌ Error reading", src, ":", err)
		return
	}

	err = os.WriteFile(dest, input, 0644)
	if err != nil {
		fmt.Println("❌ Error writing", dest, ":", err)
	} else {
		fmt.Println("✅ Successfully imported", dest)
	}
}

// 📌 Append a Line to a File (Used for PATH updates)
func appendToFile(filePath, line string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("❌ Failed to modify", filePath, ":", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString("\n" + line + "\n"); err != nil {
		fmt.Println("❌ Failed to write to", filePath, ":", err)
		return
	}
}

// 🛠 Check if a command exists (Used for `brew` and `apt` detection)
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
