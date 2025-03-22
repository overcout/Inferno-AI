package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/overcout/Inferno-AI/internal/store"
)

func RegisterHandlers(db *store.Store) {
	http.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
		}

		_, err := db.GetValidAuthLink(token)
		if err != nil {
			log.Println("[OAUTH] Invalid token:", err)
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}

		params := url.Values{}
		params.Add("client_id", googleOAuthConfig.ClientID)
		params.Add("redirect_uri", googleOAuthConfig.RedirectURL)
		params.Add("response_type", "code")
		params.Add("scope", "https://www.googleapis.com/auth/calendar https://www.googleapis.com/auth/gmail.send")
		params.Add("access_type", "offline")
		params.Add("prompt", "consent")
		params.Add("state", token)

		authURL := fmt.Sprintf("%s?%s", googleOAuthConfig.Endpoint.AuthURL, params.Encode())
		http.Redirect(w, r, authURL, http.StatusFound)
	})

	http.HandleFunc("/oauth/callback", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("state")
		code := r.URL.Query().Get("code")
		if token == "" || code == "" {
			http.Error(w, "Missing token or code", http.StatusBadRequest)
			return
		}

		authLink, err := db.GetValidAuthLink(token)
		if err != nil {
			log.Println("[OAUTH] Invalid or expired token:", err)
			http.Error(w, "Invalid or expired token", http.StatusBadRequest)
			return
		}

		oauthToken, err := googleOAuthConfig.Exchange(context.Background(), code)
		if err != nil {
			log.Println("[OAUTH] Exchange failed:", err)
			http.Error(w, "Token exchange failed", http.StatusInternalServerError)
			return
		}

		user, err := db.GetOrCreateUser(authLink.TelegramID)
		if err != nil {
			log.Println("[OAUTH] Failed to get user:", err)
			http.Error(w, "User lookup failed", http.StatusInternalServerError)
			return
		}

		user.AccessToken = oauthToken.AccessToken
		user.RefreshToken = oauthToken.RefreshToken
		user.TokenExpiry = oauthToken.Expiry
		_ = db.MarkAuthLinkUsed(token)

		if err := db.DB.Save(user).Error; err != nil {
			log.Println("[OAUTH] Failed to save tokens:", err)
			http.Error(w, "DB save failed", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "✅ Google account linked successfully!")
	})
}
