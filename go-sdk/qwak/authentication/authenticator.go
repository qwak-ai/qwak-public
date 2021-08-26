package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"qwak.ai/inference-sdk/http"
	"sync"
	"time"
)

const (
	TOKEN_EXPIRATION_BUFFER = 30 * time.Minute
)

type Authenticator struct {
	parentCtx         context.Context
	ctx               context.Context
	cancelContext     context.CancelFunc
	runingStateLocker sync.Locker
	apiKey            string
	ruunning          bool
	tokenWrapper      tokenWrapper
	runtimeAuthError  error
	httpClient        http.HttpClient
}


type AuthenticatorOptions struct {
	Ctx        context.Context
	ApiKey     string
	HttpClient http.HttpClient
}

type authResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiredAt   int64  `json:"expiredAt"`
}

type tokenWrapper struct {
	accessToken string
	expiredAt   time.Time
}

func NewAuthenticator(options *AuthenticatorOptions) (*Authenticator) {
	ctx, cancel := context.WithCancel(options.Ctx)
	authenticator := &Authenticator{
		ruunning:          false,
		runingStateLocker: &sync.Mutex{},
		parentCtx:         options.Ctx,
		ctx: ctx,
		cancelContext: cancel,
		runtimeAuthError:  nil,
		httpClient:        options.HttpClient,
		apiKey: options.ApiKey,
	}

	return authenticator
}

func (a *Authenticator) GetToken() (string, error) {
	if a.getDurationForNextAuth() > 0 {
		err := a.renewToken()

		if err != nil {
			return "", err
		}
	}

	return a.tokenWrapper.accessToken, nil
}

func (a *Authenticator) getDurationForNextAuth() time.Duration {
	now := time.Now()
	nextAuthIn := a.tokenWrapper.expiredAt.Sub(now) - (TOKEN_EXPIRATION_BUFFER)

	if nextAuthIn < 0 {
		return 0
	}

	return nextAuthIn
}

func (a *Authenticator) renewToken() error {

	tokenResponse, err := a.makeTokenRequest(a.apiKey)

	if err != nil {
		a.runtimeAuthError = err
		return err
	}

	a.tokenWrapper = tokenWrapper{
		accessToken: tokenResponse.AccessToken,
		expiredAt:   time.Unix(tokenResponse.ExpiredAt, 0),
	}

	return nil
}

func (a *Authenticator) makeTokenRequest(apiKey string) (authResponse, error) {

	decodedResponse := authResponse{}
	request, err := http.GetAuthenticationRequest(a.ctx, apiKey)

	if err != nil {
		return decodedResponse, err
	}

	body, statusCode, err := http.DoRequestWithRetry(a.httpClient, request)

	if err != nil {
		return decodedResponse, err
	}

	if statusCode == 401 {
		return decodedResponse, errors.New("wrong apiKey, authentication failed with status code 401")
	}

	if statusCode != 200 {
		return decodedResponse, fmt.Errorf("authentication failed. failed with code %d. response: %s", statusCode, body)
	}

	err = json.Unmarshal(body, &decodedResponse)

	if err != nil {
		return decodedResponse, errors.New("failed to unmarshal authentication response")
	}

	return decodedResponse, nil
}
