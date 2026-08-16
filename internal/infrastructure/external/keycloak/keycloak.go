package keycloak

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

// Client verifies JWTs issued by a Keycloak realm against its published JWKS.
type Client struct {
	verifier *oidc.IDTokenVerifier
}

// Claims holds the subset of Keycloak's JWT claims relevant for authorization.
type Claims struct {
	Subject           string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

func (claims Claims) Roles() []string {
	return claims.RealmAccess.Roles
}

// NewClient discovers the given realm's OIDC configuration and builds a token verifier for it.
func NewClient(ctx context.Context, baseUrl string, realm string, clientId string) (*Client, error) {
	issuer := fmt.Sprintf("%s/realms/%s", baseUrl, realm)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("keycloak: failed to discover provider at %s: %w", issuer, err)
	}

	// Keycloak access tokens are typically issued with aud=account rather than the
	// requesting client's ID, so the audience check is skipped here.
	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})

	return &Client{verifier: verifier}, nil
}

// VerifyToken validates the raw JWT's signature, issuer and expiry, and returns its claims.
func (client *Client) VerifyToken(ctx context.Context, rawToken string) (Claims, error) {
	idToken, err := client.verifier.Verify(ctx, rawToken)
	if err != nil {
		return Claims{}, fmt.Errorf("keycloak: invalid token: %w", err)
	}

	var claims Claims
	err = idToken.Claims(&claims)
	if err != nil {
		return Claims{}, fmt.Errorf("keycloak: failed to parse claims: %w", err)
	}

	return claims, nil
}
