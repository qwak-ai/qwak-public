package qwak

import (
	"context"
	"errors"
	"fmt"
	"qwak.ai/inference-sdk/authentication"
	"qwak.ai/inference-sdk/http"
)

const (
	PREDICTION_PATH_URL_TEMPLATE = "/v1/%s/predict"
	PREDICTION_BASE_URL_TEMPLATE = "https://models.%s.qwak.ai"
)

type QwakRealTimeClient struct {
	authenticator *authentication.Authenticator
	httpClient    http.HttpClient
	environment   string
	context       context.Context
}

type RealTimeClientConfig struct {
	ApiKey      string
	Environment string
	Context     context.Context
	HttpClient  http.HttpClient
}

func NewRealTimeClient(options RealTimeClientConfig) (RealTimeClient, error) {

	if len(options.ApiKey) == 0 {
		return nil, errors.New("api key is missing")
	}

	if len(options.Environment) == 0 {
		return nil, errors.New("environment is missing")
	}

	if options.Context == nil {
		options.Context = context.Background()
	}

	if options.HttpClient == nil {
		options.HttpClient = http.GetDefaultHttpClient()
	}

	return &QwakRealTimeClient{
		authenticator: authentication.NewAuthenticator(&authentication.AuthenticatorOptions{
			ApiKey: options.ApiKey,
			Ctx:    options.Context,
			HttpClient: options.HttpClient,
		}),
		httpClient:  options.HttpClient,
		context:     options.Context,
		environment: options.Environment,
	}, nil
}


func getPredictionUrl(environment string, modelId string) string {
	return fmt.Sprintf(PREDICTION_BASE_URL_TEMPLATE, environment) +
		fmt.Sprintf(PREDICTION_PATH_URL_TEMPLATE, modelId)
}

func (c *QwakRealTimeClient) Predict(predictionRequst *PredictionRequest) (*PredictionResponse, error) {
	if len(predictionRequst.ModelId) == 0 {
		return nil, errors.New("model id is missing in request")
	}

	token, err := c.authenticator.GetToken()

	if err != nil {
		return nil, fmt.Errorf("qwak client failed to predict: %s", err.Error())
	}

	pandaOrientedDf := predictionRequst.asPandaOrientedDf()
	predictionUrl := getPredictionUrl(c.environment, predictionRequst.ModelId)
	request, err := http.GetPredictionRequest(c.context, predictionUrl, token, pandaOrientedDf)

	if err != nil {
		return nil, fmt.Errorf("qwak client failed to predict: %s", err.Error())
	}

	responseBody, statusCode, err := http.DoRequest(c.httpClient, request)

	if err != nil {
		return nil, fmt.Errorf("qwak client failed to predict: %s", err.Error())
	}

	if statusCode != 200 {
		return nil, fmt.Errorf("qwak client failed to predict: response with status code %d. response: %s", statusCode, responseBody)
	}

	response, err := responseFromRaw(responseBody)

	if err != nil {
		return nil, fmt.Errorf("qwak client failed to parse response from model: %s", err.Error())
	}

	return response, nil
}
