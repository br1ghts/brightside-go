package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var format string

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert audio/video files (MP3 ↔ WAV, MP4 → MP3)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		convertFile(filePath, format)
	},
}

// Convert the file using ffmpeg
func convertFile(inputPath string, targetFormat string) {
	if !fileExists(inputPath) {
		fmt.Println("❌ Error: File not found!")
		return
	}

	ext := strings.ToLower(filepath.Ext(inputPath))
	fileName := strings.TrimSuffix(inputPath, ext)

	var outputPath string
	if targetFormat == "" {
		fmt.Println("⚠️ No format specified. Use --format mp3 or --format wav")
		return
	}

	switch targetFormat {
	case "mp3":
		if ext == ".mp3" {
			fmt.Println("✅ Already in MP3 format!")
			return
		}
		outputPath = fileName + ".mp3"
	case "wav":
		if ext == ".wav" {
			fmt.Println("✅ Already in WAV format!")
			return
		}
		outputPath = fileName + ".wav"
	default:
		fmt.Println("❌ Unsupported format! Use mp3 or wav.")
		return
	}

	fmt.Printf("🎵 Converting %s → %s...\n", inputPath, outputPath)

	cmd := exec.Command("ffmpeg", "-i", inputPath, outputPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Println("❌ Conversion failed:", err)
	} else {
		fmt.Println("✅ Conversion successful! File saved as", outputPath)
	}
}

// Check if the file exists
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&format, "format", "f", "", "Target format (mp3, wav)")
}
