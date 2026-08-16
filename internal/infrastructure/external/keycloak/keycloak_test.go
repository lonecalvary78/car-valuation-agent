package keycloak

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
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

var errTokenRequestFailed = errors.New("keycloak: token request failed")

func TestVerifyToken(t *testing.T) {
	ctx := context.Background()
	keycloakContainer, err := createKeycloakContainer(ctx)
	require.NoError(t, err)
	defer func() { require.NoError(t, keycloakContainer.Terminate(ctx)) }()

	baseUrl, err := keycloakBaseUrl(ctx, keycloakContainer)
	require.NoError(t, err)

	client, err := NewClient(ctx, baseUrl, adminRealm, adminClientId)
	require.NoError(t, err)

	idToken, err := fetchAdminIdToken(ctx, baseUrl)
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
	defer func() { require.NoError(t, keycloakContainer.Terminate(ctx)) }()

	baseUrl, err := keycloakBaseUrl(ctx, keycloakContainer)
	require.NoError(t, err)

	client, err := NewClient(ctx, baseUrl, adminRealm, adminClientId)
	require.NoError(t, err)

	_, err = client.VerifyToken(ctx, "not-a-valid-jwt")
	require.Error(t, err)
}

func createKeycloakContainer(ctx context.Context) (*testcontainers.DockerContainer, error) {
	container, err := testcontainers.Run(ctx, "quay.io/keycloak/keycloak:26.7",
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
	if err != nil {
		return nil, fmt.Errorf("keycloak_test: failed to start container: %w", err)
	}

	return container, nil
}

func keycloakBaseUrl(ctx context.Context, keycloakContainer *testcontainers.DockerContainer) (string, error) {
	host, err := keycloakContainer.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("keycloak_test: failed to get container host: %w", err)
	}

	mappedPort, err := keycloakContainer.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return "", fmt.Errorf("keycloak_test: failed to get mapped port: %w", err)
	}

	return "http://" + net.JoinHostPort(host, mappedPort.Port()), nil
}

// fetchAdminIdToken authenticates against the realm's built-in admin-cli client
// (public, direct-access-grants-enabled by default) so the test doesn't need to
// provision its own realm/client/user via the admin REST API. It requests the
// "openid" scope and returns the id_token, since that is what carries the "sub"
// and "preferred_username" claims that Client.VerifyToken (an ID token verifier)
// expects to parse.
func fetchAdminIdToken(ctx context.Context, baseUrl string) (string, error) {
	form := url.Values{
		"grant_type": {"password"},
		"client_id":  {adminClientId},
		"username":   {adminUsername},
		"password":   {adminPassword},
		"scope":      {"openid"},
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost,
		baseUrl+"/realms/"+adminRealm+"/protocol/openid-connect/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("keycloak_test: failed to build token request: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("keycloak_test: failed to send token request: %w", err)
	}
	defer func() {
		closeErr := response.Body.Close()
		if closeErr != nil {
			log.Printf("keycloak_test: failed to close response body: %v", closeErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: status %d", errTokenRequestFailed, response.StatusCode)
	}

	var tokenResponse struct {
		IdToken string `json:"id_token"`
	}
	err = json.NewDecoder(response.Body).Decode(&tokenResponse)
	if err != nil {
		return "", fmt.Errorf("keycloak_test: failed to decode token response: %w", err)
	}

	return tokenResponse.IdToken, nil
}
