package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ANSI color codes
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
)

type Chirp struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	printBanner()

	apiURL := "http://localhost:8080"
	if len(os.Args) > 1 {
		apiURL = os.Args[1]
	}

	fmt.Printf("%s🔍 Fetching chirps from: %s%s\n\n", Cyan, apiURL, Reset)

	chirps, err := fetchChirps(apiURL)
	if err != nil {
		fmt.Printf("%s❌ Error fetching chirps: %v%s\n", Red, err, Reset)
		fmt.Printf("%s💡 Make sure the Chirpy API is running on %s%s\n", Yellow, apiURL, Reset)
		return
	}

	if len(chirps) == 0 {
		fmt.Printf("%s📭 No chirps found. Create some chirps first!%s\n", Yellow, Reset)
		return
	}

	displayChirps(chirps)
	displayStats(chirps)
}

func printBanner() {
	banner := `
   _____ _     _                __      ___
  / ____| |   (_)              \ \    / (_)
 | |    | |__  _ _ __ _ __  _   \ \  / / _  _____      _____ _ __
 | |    | '_ \| | '__| '_ \| | | \ \/ / | |/ _ \ \ /\ / / _ \ '__|
 | |____| | | | | |  | |_) | |_| |\  /  | |  __/\ V  V /  __/ |
  \_____|_| |_|_|_|  | .__/ \__, | \/   |_|\___| \_/\_/ \___|_|
                     | |     __/ |
                     |_|    |___/
`
	colors := []string{Cyan, Blue, Magenta, Green, Yellow}
	lines := strings.Split(banner, "\n")

	for i, line := range lines {
		color := colors[i%len(colors)]
		fmt.Printf("%s%s%s\n", color, line, Reset)
	}
	fmt.Println()
}

func fetchChirps(baseURL string) ([]Chirp, error) {
	resp, err := http.Get(baseURL + "/api/chirps")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var chirps []Chirp
	if err := json.NewDecoder(resp.Body).Decode(&chirps); err != nil {
		return nil, err
	}

	return chirps, nil
}

func displayChirps(chirps []Chirp) {
	fmt.Printf("%s%s╔════════════════════════════════════════════════════════════╗%s\n", Bold, Blue, Reset)
	fmt.Printf("%s%s║                     🐦 CHIRP FEED 🐦                      ║%s\n", Bold, Blue, Reset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════╝%s\n\n", Bold, Blue, Reset)

	for i, chirp := range chirps {
		// Alternate colors for each chirp
		borderColor := Cyan
		if i%2 == 0 {
			borderColor = Magenta
		}

		fmt.Printf("%s┌─────────────────────────────────────────────────────────┐%s\n", borderColor, Reset)

		// Display chirp body with word wrapping
		words := strings.Fields(chirp.Body)
		line := "│ "
		lineLength := 0
		maxLineLength := 55

		for _, word := range words {
			if lineLength+len(word)+1 > maxLineLength {
				// Pad the rest of the line
				padding := maxLineLength - lineLength
				line += strings.Repeat(" ", padding) + "│"
				fmt.Printf("%s%s%s%s\n", borderColor, Bold, line, Reset)
				line = "│ "
				lineLength = 0
			}
			if lineLength > 0 {
				line += " "
				lineLength++
			}
			line += word
			lineLength += len(word)
		}

		// Print remaining line
		if lineLength > 0 {
			padding := maxLineLength - lineLength
			line += strings.Repeat(" ", padding) + "│"
			fmt.Printf("%s%s%s%s\n", borderColor, Bold, line, Reset)
		}

		// Footer with metadata
		fmt.Printf("%s├─────────────────────────────────────────────────────────┤%s\n", borderColor, Reset)

		timeStr := chirp.CreatedAt.Format("Jan 02, 2006 15:04")
		authorStr := chirp.UserID[:8] + "..."

		metadata := fmt.Sprintf("│ %s👤 %s  %s🕐 %s", Dim, authorStr, Dim, timeStr)
		padding := 57 - len(authorStr) - len(timeStr) - 6
		metadata += strings.Repeat(" ", padding) + Reset + borderColor + "│" + Reset

		fmt.Printf("%s%s\n", borderColor, metadata)
		fmt.Printf("%s└─────────────────────────────────────────────────────────┘%s\n\n", borderColor, Reset)
	}
}

func displayStats(chirps []Chirp) {
	fmt.Printf("%s%s╔════════════════════════════════════════════════════════════╗%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s║                     📊 STATISTICS 📊                      ║%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s╚════════════════════════════════════════════════════════════╝%s\n\n", Bold, Green, Reset)

	// Count chirps per user
	userCounts := make(map[string]int)
	totalChars := 0
	longestChirp := 0
	shortestChirp := 140

	for _, chirp := range chirps {
		userCounts[chirp.UserID]++
		length := len(chirp.Body)
		totalChars += length
		if length > longestChirp {
			longestChirp = length
		}
		if length < shortestChirp {
			shortestChirp = length
		}
	}

	avgLength := 0
	if len(chirps) > 0 {
		avgLength = totalChars / len(chirps)
	}

	fmt.Printf("%s  📝 Total Chirps:      %s%d%s\n", Yellow, Bold, len(chirps), Reset)
	fmt.Printf("%s  👥 Unique Users:      %s%d%s\n", Yellow, Bold, len(userCounts), Reset)
	fmt.Printf("%s  📏 Avg Length:        %s%d chars%s\n", Yellow, Bold, avgLength, Reset)
	fmt.Printf("%s  📐 Longest Chirp:     %s%d chars%s\n", Yellow, Bold, longestChirp, Reset)
	fmt.Printf("%s  📌 Shortest Chirp:    %s%d chars%s\n", Yellow, Bold, shortestChirp, Reset)

	// Find most active user
	maxChirps := 0
	mostActiveUser := ""
	for user, count := range userCounts {
		if count > maxChirps {
			maxChirps = count
			mostActiveUser = user
		}
	}

	if mostActiveUser != "" {
		fmt.Printf("%s  🏆 Most Active:       %s%s (%d chirps)%s\n",
			Yellow, Bold, mostActiveUser[:8]+"...", maxChirps, Reset)
	}

	// Fun ASCII art bar chart
	fmt.Printf("\n%s  📊 Activity Bar:%s\n", Cyan, Reset)
	maxBar := 40
	for user, count := range userCounts {
		barLength := (count * maxBar) / len(chirps)
		if barLength == 0 {
			barLength = 1
		}
		bar := strings.Repeat("█", barLength)
		fmt.Printf("     %s%-12s %s%s%s %d\n",
			Dim, user[:8]+"...", Green, bar, Reset, count)
	}

	fmt.Println()
}
