package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// grabCmd represents the grab command
var grabCmd = &cobra.Command{
	Use:   "grab [URL]",
	Short: "Download videos, images, or files from the internet",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		downloadFile(url)
	},
}

// Detect file type and download appropriately
func downloadFile(url string) {
	fmt.Println("🔍 Detecting file type...")

	// Check if URL is a YouTube link (or TikTok)
	if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") || strings.Contains(url, "tiktok.com") {
		fmt.Println("🎥 Detected Video Platform! Using yt-dlp...")
		downloadWithYTDLP(url)
		return
	}

	// Check if it's a direct file download (like .jpg, .mp3, .pdf, etc.)
	ext := filepath.Ext(url)
	if ext != "" {
		fmt.Printf("📂 Detected File Download! File Type: %s\n", ext)
		downloadWithWget(url)
		return
	}

	// If it's a webpage, save it as an offline page
	if isWebPage(url) {
		fmt.Println("🌍 Detected Webpage! Saving for offline use...")
		downloadWithWget(url)
		return
	}

	// If no match, just try to download it
	fmt.Println("📡 Unknown type. Attempting to download...")
	downloadWithCurl(url)
}

// Uses yt-dlp for video downloads
func downloadWithYTDLP(url string) {
	cmd := exec.Command("yt-dlp", "-f", "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]", "-o", "~/Downloads/%(title)s.%(ext)s", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Failed to download video:", err)
	} else {
		fmt.Println("✅ Download complete! Saved in ~/Downloads")
	}
}

// Uses wget for direct file downloads
func downloadWithWget(url string) {
	cmd := exec.Command("wget", "-c", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Failed to download file:", err)
	} else {
		fmt.Println("✅ Download complete!")
	}
}

// Uses curl if wget is unavailable
func downloadWithCurl(url string) {
	cmd := exec.Command("curl", "-O", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("❌ Failed to download:", err)
	} else {
		fmt.Println("✅ Download complete!")
	}
}

// Simple check if it's a webpage (not a direct file)
func isWebPage(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	contentType := resp.Header.Get("Content-Type")
	return strings.Contains(contentType, "text/html")
}

func init() {
	rootCmd.AddCommand(grabCmd)
}
