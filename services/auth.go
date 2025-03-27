package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, ctx context.Context) *http.Client {
	tokenFile := "token.json"
	token, err := tokenFromFile(tokenFile)
	// if err != nil {
	// 	token = getTokenFromWeb(config)
	// 	saveToken(tokenFile, token)
	// }

	tokenSource := config.TokenSource(ctx, token)
	refreshedToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("unable to refresh token: %v", err)
	}

	// 🔹 Save the new refreshed token back to file
	saveToken(tokenFile, refreshedToken)

	return config.Client(ctx, token)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	var authCode string
	fmt.Print("Enter authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	token, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return token
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	// f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache OAuth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Initialize Gmail Service
func getGmailService() (*gmail.Service, error) {
	ctx := context.Background()

	// Load credentials from the JSON file
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Configure the OAuth2 client
	config, err := google.ConfigFromJSON(b, gmail.GmailModifyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// client := getClient(config, ctx)

	// Load token from file
	tokenFile := "token.json"
	token, err := tokenFromFile(tokenFile)
	if err != nil {
		// return nil, fmt.Errorf("unable to load token: %v", err)
		token = getTokenFromWeb(config)
		saveToken(tokenFile, token)
	}

	log.Printf("Access Token: %s\n", token.AccessToken)
	log.Printf("Refresh Token: %s\n", token.RefreshToken)
	log.Printf("Token Expiry: %v\n", token.Expiry)

	// ❗️Check if refresh_token exists (it should never be empty)
	if token.RefreshToken == "" {
		log.Fatalf("❌ Refresh Token is missing! Re-authentication required.")
	}

	// 🔥 ✅ Use TokenSource to refresh tokens automatically
	tokenSource := config.TokenSource(ctx, token)

	// 🔄 If token is expired, refresh it with retries
	if token.Expiry.Before(time.Now()) {
		log.Println("🔄 Token expired, refreshing...")
		for i := 0; i < 3; i++ { // Retry up to 3 times
			newToken, err := tokenSource.Token()
			if err == nil {
				saveToken(tokenFile, newToken)
				token = newToken
				break
			}
			log.Printf("❌ Failed to refresh token (Attempt %d): %v", i+1, err)
			time.Sleep(time.Duration(i+1) * time.Second) // Exponential backoff
		}
	}

	// refreshedToken, err := tokenSource.Token()
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to refresh token: %v", err)
	// }

	// // 🔹 Save the new refreshed token back to file
	// saveToken(tokenFile, refreshedToken)

	// 🔥 Use `tokenSource` instead of old token
	client := oauth2.NewClient(ctx, tokenSource)

	// Create the Gmail service
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	log.Println("✅ Gmail service initialized successfully.")
	return srv, nil
}
