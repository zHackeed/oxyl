package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"google.golang.org/grpc/credentials"
	"zhacked.me/oxyl/agent/internal/models"
)

type activeToken struct {
	Token      string
	expiration time.Time
}

type AuthenticationService struct {
	// todo: remove for prod - only for testing as right now
	loginEndpoint    string
	refreshEndpoint  string
	shutdownEndpoint string

	identifier string

	enrollmentToken *string

	tokenMu sync.RWMutex
	token   *activeToken

	refreshTokenMu sync.RWMutex
	refreshToken   *activeToken

	client *http.Client
}

var _ credentials.PerRPCCredentials = (*AuthenticationService)(nil)

func NewAuthService(identifier, loginEndpoint, refreshEndpoint, shutdownEndpoint string) (*AuthenticationService, error) {
	if loginEndpoint == "" || refreshEndpoint == "" || shutdownEndpoint == "" {
		return nil, fmt.Errorf("the login or refresh url is invalid or malformed")
	}
	_, err := url.Parse(loginEndpoint)
	if err != nil {
		return nil, fmt.Errorf("the login url is invalid or malformed. %v", err)
	}

	_, err = url.Parse(refreshEndpoint)
	if err != nil {
		return nil, fmt.Errorf("the refresh url is invalid or malformed. %v", err)
	}

	_, err = url.Parse(shutdownEndpoint)
	if err != nil {
		return nil, fmt.Errorf("the shutdown url is invalid or malformed. %v", err)
	}

	//todo: add enrollment token for authentication if present

	srv := new(AuthenticationService{
		client:           &http.Client{},
		identifier:       identifier,
		loginEndpoint:    loginEndpoint,
		refreshEndpoint:  refreshEndpoint,
		shutdownEndpoint: shutdownEndpoint,
	})

	return srv, nil
}

func (s *AuthenticationService) StartTicking(ctx context.Context) {
	err := s.CreateAuthRequest()
	if err != nil {
		slog.Error("unable to create auth request", "error", err)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if time.Now().After(s.refreshToken.expiration.Add(-time.Minute)) {
					if err := s.CreateAuthRequest(); err != nil {
						slog.Error("unable to create auth request", "error", err)
					}
					continue
				}

				if time.Now().After(s.token.expiration.Add(-time.Minute)) {
					if err := s.RefreshToken(); err != nil {
						slog.Error("unable to refresh token", "error", err)
					}
				}
			}
		}
	}()
}

func (s *AuthenticationService) GetToken() string {
	s.tokenMu.RLock()
	defer s.tokenMu.RUnlock()

	return s.token.Token
}

func (s *AuthenticationService) IsTokenExpired() bool {
	s.tokenMu.RLock()
	defer s.tokenMu.RUnlock()

	return time.Now().After(s.token.expiration.Add(-time.Minute))
}

func (s *AuthenticationService) RefreshToken() error {
	request, err := http.NewRequest("POST", s.refreshEndpoint, bytes.NewBuffer(nil))
	if err != nil {
		return fmt.Errorf("could not create refresh token request: %v", err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.refreshToken.Token))

	resp, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("refresh token request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh token request failed: %v", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read data of refresh token request failed: %v", err)
	}

	var response models.AuthenticationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return fmt.Errorf("unmarshalling response body failed: %v", err)
	}

	s.tokenMu.Lock()
	s.refreshTokenMu.Lock()

	defer s.tokenMu.Unlock()
	defer s.refreshTokenMu.Unlock()

	s.token = &activeToken{
		Token:      response.AccessToken.Token,
		expiration: response.RefreshToken.ExpiresAt,
	}

	s.refreshToken = &activeToken{
		Token:      response.RefreshToken.Token,
		expiration: response.RefreshToken.ExpiresAt,
	}

	return nil
}

func (s *AuthenticationService) CreateAuthRequest() error {
	body := models.AuthenticationRequest{
		Identifier: s.identifier,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling request body failed: %v", err)
	}

	resp, err := s.client.Post(s.loginEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("auth request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth request failed: %v", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("reading response body failed: %v", err)
	}

	var response models.AuthenticationResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		return fmt.Errorf("unmarshalling response body failed: %v", err)
	}

	s.tokenMu.Lock()
	s.refreshTokenMu.Lock()
	defer func() {
		s.tokenMu.Unlock()
		s.refreshTokenMu.Unlock()
	}()

	s.token = &activeToken{
		Token:      response.AccessToken.Token,
		expiration: response.RefreshToken.ExpiresAt,
	}

	s.refreshToken = &activeToken{
		Token:      response.RefreshToken.Token,
		expiration: response.RefreshToken.ExpiresAt,
	}

	return nil
}

func (s *AuthenticationService) RequestShutdown() error {
	s.refreshTokenMu.RLock()
	defer s.refreshTokenMu.RUnlock()

	body := models.AuthenticationShutdownRequest{
		AgentId: s.identifier,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling request body failed: %v", err)
	}

	req, err := http.NewRequest("POST", s.shutdownEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("creating request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.refreshToken.Token))

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("auth request failed: %v", resp.Status)
	}

	return nil
}

// ----- gRPC authentication methods

func (s *AuthenticationService) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + s.token.Token,
		//"ag_enrollment": "",
	}, nil
}

func (s *AuthenticationService) RequireTransportSecurity() bool {
	return true // we require TLS (https)
}
