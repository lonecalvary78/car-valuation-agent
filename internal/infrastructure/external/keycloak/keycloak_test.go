package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	adminUsername = "admin"
	adminPassword = "admin"
	adminRealm    = "master"
	adminClientId = "admin-cli"
)

func TestVerifyToken(t *testing.T) {
	ctx := context.Background()
	keycloakContainer, err := createKeycloakContainer(ctx)
	require.NoError(t, err)
	defer keycloakContainer.Terminate(ctx)

	baseUrl, err := keycloakBaseUrl(ctx, keycloakContainer)
	require.NoError(t, err)

	client, err := NewClient(ctx, baseUrl, adminRealm, adminClientId)
	require.NoError(t, err)

	idToken, err := fetchAdminIdToken(baseUrl)
	require.NoError(t, err)

	claims, err := client.VerifyToken(ctx, idToken)
	require.NoError(t, err)
	require.NotEmpty(t, claims.Subject)
	require.Equal(t, adminUsername, claims.PreferredUsername)
}

func TestVerifyToken_InvalidToken(t *testing.T) {
	ctx := context.Background()
	keycloakContainer, err := createKeycloakContainer(ctx)
	require.NoError(t, err)
	defer keycloakContainer.Terminate(ctx)

	baseUrl, err := keycloakBaseUrl(ctx, keycloakContainer)
	require.NoError(t, err)

	client, err := NewClient(ctx, baseUrl, adminRealm, adminClientId)
	require.NoError(t, err)

	_, err = client.VerifyToken(ctx, "not-a-valid-jwt")
	require.Error(t, err)
}

func createKeycloakContainer(ctx context.Context) (*testcontainers.DockerContainer, error) {
	return testcontainers.Run(ctx, "quay.io/keycloak/keycloak:26.7",
		testcontainers.WithEnv(map[string]string{
			"KEYCLOAK_ADMIN":          adminUsername,
			"KEYCLOAK_ADMIN_PASSWORD": adminPassword,
		}),
		testcontainers.WithCmd("start-dev"),
		testcontainers.WithExposedPorts("8080/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForHTTP("/realms/"+adminRealm+"/.well-known/openid-configuration").
				WithPort("8080/tcp").
				WithStartupTimeout(2*time.Minute),
		),
	)
}

func keycloakBaseUrl(ctx context.Context, keycloakContainer *testcontainers.DockerContainer) (string, error) {
	host, err := keycloakContainer.Host(ctx)
	if err != nil {
		return "", err
	}

	mappedPort, err := keycloakContainer.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%s", host, mappedPort.Port()), nil
}

// fetchAdminIdToken authenticates against the realm's built-in admin-cli client
// (public, direct-access-grants-enabled by default) so the test doesn't need to
// provision its own realm/client/user via the admin REST API. It requests the
// "openid" scope and returns the id_token, since that is what carries the "sub"
// and "preferred_username" claims that Client.VerifyToken (an ID token verifier)
// expects to parse.
func fetchAdminIdToken(baseUrl string) (string, error) {
	form := url.Values{
		"grant_type": {"password"},
		"client_id":  {adminClientId},
		"username":   {adminUsername},
		"password":   {adminPassword},
		"scope":      {"openid"},
	}

	response, err := http.PostForm(baseUrl+"/realms/"+adminRealm+"/protocol/openid-connect/token", form)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("keycloak: token request failed with status %d", response.StatusCode)
	}

	var tokenResponse struct {
		IdToken string `json:"id_token"`
	}
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.IdToken, nil
}
