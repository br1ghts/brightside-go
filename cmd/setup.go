package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Setup flags
var resetFlag bool
var silentFlag bool

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Installs dependencies and configures Brightside-Go",
	Run: func(cmd *cobra.Command, args []string) {
		runSetup()
	},
}

func runSetup() {
	fmt.Println("🚀 Running Brightside-Go Setup...\n")

	if resetFlag {
		resetInstallation()
	}

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

	packages := []string{"git", "yt-dlp", "ffmpeg", "wget", "zsh", "font-hack-nerd-font"}
	for _, pkg := range packages {
		if !commandExists(pkg) {
			fmt.Printf("🔹 Installing %s...\n", pkg)
			exec.Command("brew", "install", pkg).Run()
		}
	}

	installP10K()
}

// 🛠 Install Dependencies (Linux)
func installLinuxDependencies() {
	fmt.Println("🔹 Checking for APT...")
	if !commandExists("apt") {
		fmt.Println("❌ APT package manager not found! Make sure you're on a Debian-based system.")
		os.Exit(1)
	}

	packages := []string{"git", "yt-dlp", "ffmpeg", "wget", "zsh", "fonts-powerline"}
	for _, pkg := range packages {
		if !commandExists(pkg) {
			fmt.Printf("🔹 Installing %s...\n", pkg)
			exec.Command("sudo", "apt", "install", "-y", pkg).Run()
		}
	}

	installP10K()
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

// 🔧 Configure Shell (Import `.zshrc`, `.p10k.zsh`, and install plugins)
// 🔧 Configure Shell (Import `.zshrc` and plugins)
func configureShell() {
	shellConfig := getShellConfig()
	if shellConfig == "" {
		fmt.Println("❌ Could not detect shell configuration file.")
		return
	}

	// 🛠 Ensure `.zshrc` is restored
	restoreZshConfig()

	// ✅ Now it's safe to modify `.zshrc`
	fmt.Println("🔧 Adding Brightside-Go to PATH in", shellConfig)
	exportCmd := "export PATH=\"/usr/local/bin:$PATH\""
	appendToFile(shellConfig, exportCmd)

	// 🛠 Install Powerlevel10k & Plugins
	installP10K()

	// 🛠 Double-check Oh My Zsh installation before sourcing `.zshrc`
	if _, err := os.Stat(os.Getenv("HOME") + "/.oh-my-zsh/oh-my-zsh.sh"); os.IsNotExist(err) {
		fmt.Println("⚠️ Warning: Oh My Zsh did not install correctly. Please run 'brightside setup' again.")
		return
	}

	// 🔄 Source `.zshrc` properly AFTER ensuring installation is done
	fmt.Println("🔄 Sourcing .zshrc to apply changes...")
	err := exec.Command("zsh", "-c", "sleep 2 && source ~/.zshrc").Run()
	if err != nil {
		fmt.Println("⚠️ Warning: Could not source .zshrc automatically. Try running 'source ~/.zshrc' manually.")
	} else {
		fmt.Println("✅ Shell configuration complete! Run 'source ~/.zshrc' if needed.")
	}
}

func restoreZshConfig() {
	zshrcPath := os.Getenv("HOME") + "/.zshrc"
	configPath := ("config/.zshrc")

	fmt.Println("🛠 Ensuring .zshrc is correctly set...")

	// Force overwrite with Brightside's .zshrc
	err := copyFile(configPath, zshrcPath)
	if err != nil {
		fmt.Println("❌ Failed to overwrite .zshrc:", err)
	} else {
		fmt.Println("✅ .zshrc successfully replaced with Brightside config.")
	}
}

// 📂 Copy a File from One Location to Another (Ensures Parent Directories Exist)
func copyFile(src, dest string) error {
	fmt.Printf("📂 Copying %s → %s\n", src, dest)

	// Ensure source exists
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("❌ ERROR: Could not open source file %s: %v", src, err)
	}
	defer srcFile.Close()

	// Ensure destination directory exists
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("❌ ERROR: Could not create parent directory for %s: %v", dest, err)
	}

	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("❌ ERROR: Could not create destination file %s: %v", dest, err)
	}
	defer destFile.Close()

	// Copy file contents
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("❌ ERROR: Failed to copy file %s to %s: %v", src, dest, err)
	}

	// Set correct permissions
	srcInfo, err := os.Stat(src)
	if err == nil {
		os.Chmod(dest, srcInfo.Mode())
	}

	fmt.Printf("✅ Successfully copied %s → %s\n", src, dest)
	return nil
}

// 📌 Append a Line to a File (Used for PATH updates)
func appendToFile(filePath, line string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("❌ Failed to modify", filePath, ":", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("\n" + line + "\n")
	if err != nil {
		fmt.Println("❌ Failed to write to", filePath, ":", err)
	} else {
		fmt.Println("✅ Updated", filePath)
	}
}

// 🔥 Install Oh My Zsh, Powerlevel10k, and Zsh Plugins
func installP10K() {
	fmt.Println("🎨 Checking Oh My Zsh & Powerlevel10k installation...")

	// ✅ Ensure Zsh is installed first
	if !commandExists("zsh") {
		fmt.Println("⚡ Installing Zsh...")
		exec.Command("sudo", "apt", "install", "-y", "zsh").Run()
	}

	// ✅ Ensure Oh My Zsh is installed
	ohMyZshPath := os.Getenv("HOME") + "/.oh-my-zsh"
	if _, err := os.Stat(ohMyZshPath); os.IsNotExist(err) {
		fmt.Println("⚡ Installing Oh My Zsh...")
		cmd := exec.Command("/bin/bash", "-c", "curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh | bash")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("❌ Failed to install Oh My Zsh:", err)
			return
		}
	} else {
		fmt.Println("✅ Oh My Zsh is already installed.")
	}

	// ✅ Ensure Powerlevel10k is installed
	zshCustom := os.Getenv("HOME") + "/.oh-my-zsh/custom"
	os.Setenv("ZSH_CUSTOM", zshCustom)

	p10kPath := zshCustom + "/themes/powerlevel10k"
	if _, err := os.Stat(p10kPath); os.IsNotExist(err) {
		fmt.Println("🎨 Installing Powerlevel10k...")
		err := exec.Command("git", "clone", "--depth=1", "https://github.com/romkatv/powerlevel10k.git", p10kPath).Run()
		if err != nil {
			fmt.Println("❌ Failed to install Powerlevel10k:", err)
		}
	} else {
		fmt.Println("✅ Powerlevel10k is already installed.")
	}

	// ✅ Ensure Zsh Plugins are Installed
	pluginDir := zshCustom + "/plugins"
	os.MkdirAll(pluginDir, os.ModePerm)

	// 🔹 Autosuggestions
	autosuggestionsPath := pluginDir + "/zsh-autosuggestions"
	if _, err := os.Stat(autosuggestionsPath); os.IsNotExist(err) {
		fmt.Println("💡 Installing zsh-autosuggestions...")
		err := exec.Command("git", "clone", "https://github.com/zsh-users/zsh-autosuggestions", autosuggestionsPath).Run()
		if err != nil {
			fmt.Println("❌ Failed to install zsh-autosuggestions:", err)
		}
	}

	// 🔹 Syntax Highlighting
	syntaxHighlightingPath := pluginDir + "/zsh-syntax-highlighting"
	if _, err := os.Stat(syntaxHighlightingPath); os.IsNotExist(err) {
		fmt.Println("💡 Installing zsh-syntax-highlighting...")
		err := exec.Command("git", "clone", "https://github.com/zsh-users/zsh-syntax-highlighting", syntaxHighlightingPath).Run()
		if err != nil {
			fmt.Println("❌ Failed to install zsh-syntax-highlighting:", err)
		}
	}

	fmt.Println("✅ Powerlevel10k and plugins installed successfully!")
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

// 📂 Reset Installation
func resetInstallation() {
	fmt.Println("🗑 Resetting Brightside-Go installation...")

	os.Remove("/usr/local/bin/brightside")
	os.RemoveAll(os.Getenv("HOME") + "/.oh-my-zsh")
	os.RemoveAll(os.Getenv("HOME") + "/.p10k.zsh")
	os.RemoveAll(os.Getenv("HOME") + "/.zshrc")

	fmt.Println("✅ Brightside-Go has been reset! Run 'brightside setup' again.")
}

// 📂 Check if a command exists
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func init() {
	setupCmd.Flags().BoolVar(&resetFlag, "reset", false, "Reset Brightside-Go installation")
	setupCmd.Flags().BoolVar(&silentFlag, "silent", false, "Run Brightside-Go setup without prompts")
	rootCmd.AddCommand(setupCmd)
}
