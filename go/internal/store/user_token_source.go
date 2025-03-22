package store

import (
	"context"

	"golang.org/x/oauth2"
)

// UserTokenSource allows using a stored access token with refresh support
type UserTokenSource struct {
	User *User
}

func (u *UserTokenSource) Token() (*oauth2.Token, error) {
	tok := &oauth2.Token{
		AccessToken:  u.User.AccessToken,
		RefreshToken: u.User.RefreshToken,
		Expiry:       u.User.TokenExpiry,
		TokenType:    "Bearer",
	}

	if tok.Valid() {
		return tok, nil
	}

	ts := (&oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}).TokenSource(context.Background(), tok)

	newToken, err := ts.Token()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}
