package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

// News source file (stored in user's config dir)
var newsConfigFile = filepath.Join(os.Getenv("HOME"), ".brightside_news.json")

// Default news sources (used if no config exists)
var defaultNewsSources = map[string][]string{
	"Tech": {
		"https://www.theverge.com/rss/index.xml",
		"https://www.wired.com/feed/rss",
		"https://www.techradar.com/rss",
		"https://rss.nytimes.com/services/xml/rss/nyt/Technology.xml",
	},
	"World": {
		"http://feeds.bbci.co.uk/news/world/rss.xml",
		"https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
		"https://www.aljazeera.com/xml/rss/all.xml",
	},
	"Hacker": {
		"https://news.ycombinator.com/rss",
	},
}

// Colors for terminal output
const (
	Blue   = "\033[1;34m"
	Cyan   = "\033[1;36m"
	Green  = "\033[1;32m"
	Gray   = "\033[1;30m"
	Reset  = "\033[0m"
)

// Load user-defined news sources or default ones
func loadNewsSources() map[string][]string {
	data, err := os.ReadFile(newsConfigFile)
	if err != nil {
		// If no custom sources exist, use default
		return defaultNewsSources
	}
	var sources map[string][]string
	if err := json.Unmarshal(data, &sources); err != nil {
		return defaultNewsSources
	}
	return sources
}

// Save user-defined news sources
func saveNewsSources(sources map[string][]string) {
	data, err := json.MarshalIndent(sources, "", "  ")
	if err != nil {
		fmt.Println("❌ Failed to save news sources!")
		return
	}
	_ = os.WriteFile(newsConfigFile, data, 0644)
}

// Fetch latest news from RSS feeds
func fetchNews(category string, limit int) {
	sources := loadNewsSources()

	feeds, exists := sources[category]
	if !exists {
		fmt.Println("❌ Invalid category. Available categories:")
		for k := range sources {
			fmt.Printf(Green+" - %s\n"+Reset, k)
		}
		return
	}

	// Fancy loading effect
	fmt.Printf("\n📡 Fetching %s news ", category)
	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("\n")

	fp := gofeed.NewParser()

	for _, feedURL := range feeds {
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			fmt.Printf("⚠️ Error fetching %s: %v\n", feedURL, err)
			continue
		}

		fmt.Printf(Blue+"📰 %s\n"+Reset, feed.Title)

		count := 0
		for _, item := range feed.Items {
			if count >= limit {
				break
			}
			fmt.Printf(Green+"  🔹 %s "+Cyan+"(%s)\n"+Reset, item.Title, item.Link)
			count++
		}
		fmt.Println(Gray + "-------------------------------------------------" + Reset)
	}
}

// Add a news source
func addNewsSource(category, url string) {
	sources := loadNewsSources()
	sources[category] = append(sources[category], url)
	saveNewsSources(sources)
	fmt.Println(Green + "✅ Source added successfully!" + Reset)
}

// Remove a news source
func removeNewsSource(category, url string) {
	sources := loadNewsSources()
	if _, exists := sources[category]; !exists {
		fmt.Println("❌ Category not found!")
		return
	}

	newFeeds := []string{}
	for _, feed := range sources[category] {
		if feed != url {
			newFeeds = append(newFeeds, feed)
		}
	}

	sources[category] = newFeeds
	saveNewsSources(sources)
	fmt.Println(Green + "✅ Source removed successfully!" + Reset)
}

// News command
var limit int

var newsCmd = &cobra.Command{
	Use:   "news [category]",
	Short: "Fetch latest news from RSS feeds",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fetchNews(args[0], limit)
	},
}


// Add source command
var newsAddCmd = &cobra.Command{
	Use:   "news-add [category] [url]",
	Short: "Add a news source to a category",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		addNewsSource(args[0], args[1])
	},
}

// Remove source command
var newsRemoveCmd = &cobra.Command{
	Use:   "news-remove [category] [url]",
	Short: "Remove a news source from a category",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		removeNewsSource(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(newsCmd)
	newsCmd.Flags().IntVarP(&limit, "limit", "l", 5, "Number of articles to fetch (default: 5)")
	rootCmd.AddCommand(newsAddCmd)
	rootCmd.AddCommand(newsRemoveCmd)
}