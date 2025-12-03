package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
	Bold   = "\033[1m"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateChirpRequest struct {
	Body string `json:"body"`
}

var funChirps = []string{
	"Just deployed my first Go API! 🚀",
	"Why do programmers prefer dark mode? Because light attracts bugs! 🐛",
	"Coffee: because adulting is hard ☕",
	"There are only 10 types of people: those who understand binary and those who don't",
	"Debugging is like being a detective in a crime movie where you're also the murderer 🔍",
	"My code works and I don't know why 🤷",
	"Git happens! Just committed 500 lines of awesome code 💪",
	"REST APIs are RESTful because they REST on HTTP 🛌",
	"Today's mood: sudo make me a sandwich 🥪",
	"Keep calm and clear cache 🧹",
	"Hello World! Actually, hello universe! 🌌",
	"Tabs vs Spaces? I choose violence 😈",
	"404: Motivation not found. Please try again later",
	"Programming is 10% writing code and 90% figuring out why it doesn't work",
	"I don't always test my code, but when I do, I do it in production 😎",
	"Roses are red, violets are blue, unexpected { on line 32",
	"Life is too short for bad code. Write poetry instead 📝",
	"Cloud computing is just someone else's computer ☁️",
	"May the fork be with you! 🍴",
	"I speak fluent sarcasm and broken code",
	"Coding in Go feels like driving a race car! 🏎️",
	"The best thing about boolean jokes is they're either funny or they're not",
	"I changed my password to 'incorrect'. Now when I forget it, it tells me my password is incorrect",
	"Artificial intelligence is no match for natural stupidity",
	"I would tell you a UDP joke, but you might not get it",
	"A SQL query walks into a bar, walks up to two tables and asks... Can I join you? 🍻",
}

func main() {
	apiURL := "http://localhost:8080"
	if len(os.Args) > 1 {
		apiURL = os.Args[1]
	}

	fmt.Printf("%s%s╔═══════════════════════════════════════════════╗%s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s║     🎨 Chirpy Demo Data Generator 🎨        ║%s\n", Bold, Cyan, Reset)
	fmt.Printf("%s%s╚═══════════════════════════════════════════════╝%s\n\n", Bold, Cyan, Reset)

	rand.Seed(time.Now().UnixNano())

	fmt.Printf("%s🌐 API URL: %s%s\n\n", Yellow, apiURL, Reset)

	// Create 3 demo users
	users := []struct {
		email    string
		password string
		token    string
	}{
		{"alice@chirpy.com", "password123", ""},
		{"bob@chirpy.com", "password123", ""},
		{"charlie@chirpy.com", "password123", ""},
	}

	fmt.Printf("%s👥 Creating demo users...%s\n", Cyan, Reset)
	for i := range users {
		// Try to create user
		createReq := CreateUserRequest{
			Email:    users[i].email,
			Password: users[i].password,
		}

		body, _ := json.Marshal(createReq)
		resp, err := http.Post(apiURL+"/api/users", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("  %s❌ Error creating user %s: %v%s\n", Yellow, users[i].email, err, Reset)
			continue
		}
		resp.Body.Close()

		// Login to get token
		loginReq := LoginRequest{
			Email:    users[i].email,
			Password: users[i].password,
		}

		body, _ = json.Marshal(loginReq)
		resp, err = http.Post(apiURL+"/api/login", "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Printf("  %s❌ Error logging in as %s: %v%s\n", Yellow, users[i].email, err, Reset)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var loginResp LoginResponse
			json.NewDecoder(resp.Body).Decode(&loginResp)
			users[i].token = loginResp.Token
			fmt.Printf("  %s✅ %s (token obtained)%s\n", Green, users[i].email, Reset)
		} else {
			fmt.Printf("  %s⚠️  %s (already exists, trying login...)%s\n", Yellow, users[i].email, Reset)
		}
		resp.Body.Close()
	}

	fmt.Printf("\n%s🐦 Creating fun chirps...%s\n", Cyan, Reset)

	// Shuffle chirps for randomness
	shuffled := make([]string, len(funChirps))
	copy(shuffled, funChirps)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	chirpCount := 0
	maxChirps := 15 // Create 15 chirps

	for i := 0; i < maxChirps && i < len(shuffled); i++ {
		// Randomly select a user
		userIdx := rand.Intn(len(users))
		if users[userIdx].token == "" {
			continue
		}

		chirpReq := CreateChirpRequest{
			Body: shuffled[i],
		}

		body, _ := json.Marshal(chirpReq)
		req, _ := http.NewRequest("POST", apiURL+"/api/chirps", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+users[userIdx].token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("  %s❌ Error creating chirp: %v%s\n", Yellow, err, Reset)
			continue
		}

		if resp.StatusCode == http.StatusCreated {
			chirpCount++
			preview := shuffled[i]
			if len(preview) > 50 {
				preview = preview[:50] + "..."
			}
			fmt.Printf("  %s✅ Chirp #%d: \"%s\"%s\n", Green, chirpCount, preview, Reset)
		}
		resp.Body.Close()

		// Small delay to make it look cool
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("\n%s%s🎉 Demo data created successfully! 🎉%s\n", Bold, Green, Reset)
	fmt.Printf("%s%s   Created %d chirps from %d users%s\n\n", Bold, Green, chirpCount, len(users), Reset)
	fmt.Printf("%sRun the chirpy-viewer to see them in action!%s\n", Cyan, Reset)
	fmt.Printf("%s  cd cmd/chirpy-viewer && go run main.go%s\n\n", Yellow, Reset)
}
