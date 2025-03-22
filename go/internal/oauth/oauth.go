package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"github.com/overcout/Inferno-AI/internal/config"
	"github.com/overcout/Inferno-AI/internal/store"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOAuthConfig *oauth2.Config

func InitOAuth(cfg *config.Config) {
	googleOAuthConfig = &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.OAuthRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/calendar",
			"https://www.googleapis.com/auth/gmail.send",
		},
		Endpoint: google.Endpoint,
	}

	log.Println("[OAUTH] Loaded ClientID:", googleOAuthConfig.ClientID)
	log.Println("[OAUTH] Loaded RedirectURL:", googleOAuthConfig.RedirectURL)
}

func StartServer(db *store.Store) {
	http.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "Missing user_id", http.StatusBadRequest)
			return
		}
		log.Println("[OAUTH] Redirecting user", userID)

		params := url.Values{}
		params.Add("client_id", googleOAuthConfig.ClientID)
		params.Add("redirect_uri", googleOAuthConfig.RedirectURL)
		params.Add("response_type", "code")
		params.Add("scope", "https://www.googleapis.com/auth/calendar https://www.googleapis.com/auth/gmail.send")
		params.Add("access_type", "offline")
		params.Add("prompt", "consent")
		params.Add("state", userID)

		authURL := fmt.Sprintf("%s?%s", googleOAuthConfig.Endpoint.AuthURL, params.Encode())
		http.Redirect(w, r, authURL, http.StatusFound)
	})

	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[OAUTH] Callback hit")

		code := r.URL.Query().Get("code")
		userID := r.URL.Query().Get("state")
		if code == "" || userID == "" {
			http.Error(w, "Missing code or user_id", http.StatusBadRequest)
			return
		}

		log.Println("[OAUTH] Exchanging code for token...")
		token, err := googleOAuthConfig.Exchange(context.Background(), code)
		if err != nil {
			log.Println("[OAUTH] Token exchange failed:", err)
			http.Error(w, "Token exchange failed", http.StatusInternalServerError)
			return
		}

		id, err := parseUserID(userID)
		if err != nil {
			log.Println("[OAUTH] Invalid user_id:", err)
			http.Error(w, "Invalid user_id", http.StatusBadRequest)
			return
		}

		log.Println("[OAUTH] Fetching user:", id)
		user, err := db.GetOrCreateUser(id)
		if err != nil {
			log.Println("[OAUTH] User lookup failed:", err)
			http.Error(w, "User lookup failed", http.StatusInternalServerError)
			return
		}

		log.Println("[OAUTH] Saving tokens to database...")
		user.AccessToken = token.AccessToken
		user.RefreshToken = token.RefreshToken
		user.TokenExpiry = token.Expiry

		if err := db.DB.Save(user).Error; err != nil {
			log.Println("[OAUTH] Failed to save user:", err)
			http.Error(w, "DB update failed", http.StatusInternalServerError)
			return
		}

		log.Println("[OAUTH] Success! Account linked for user", user.ID)
		fmt.Fprintln(w, "✅ Google account linked successfully!")
	})

	log.Println("[OAUTH] Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func parseUserID(s string) (int64, error) {
	var id int64
	err := json.Unmarshal([]byte(s), &id)
	if err != nil {
		_, err = fmt.Sscanf(s, "%d", &id)
	}
	return id, err
}
