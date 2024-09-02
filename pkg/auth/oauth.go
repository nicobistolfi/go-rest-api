package auth

import (
	"context"
	"errors"

	"go-boilerplate/internal/config"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
	oidcProvider *oidc.Provider
	oidcConfig   oauth2.Config
)

func InitOAuth(cfg *config.Config) error {
	var err error
	oidcProvider, err = oidc.NewProvider(context.Background(), cfg.OIDCIssuer)
	if err != nil {
		return err
	}

	oidcConfig = oauth2.Config{
		ClientID:     cfg.OAuthClientID,
		ClientSecret: cfg.OAuthClientSecret,
		RedirectURL:  cfg.OAuthRedirectURL,
		Endpoint:     oidcProvider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return nil
}

func ValidateOAuthToken(ctx context.Context, token string, cfg *config.Config) (*oidc.UserInfo, error) {
	if oidcProvider == nil {
		return nil, errors.New("OIDC provider not initialized")
	}

	oauth2Token := &oauth2.Token{AccessToken: token}
	userInfo, err := oidcProvider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}
