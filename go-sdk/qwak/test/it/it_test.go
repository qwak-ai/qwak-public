package it_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/qwak-ai/qwak-public/go-sdk/qwak"
	qwakhttp "github.com/qwak-ai/qwak-public/go-sdk/qwak/http"
	"github.com/qwak-ai/qwak-public/go-sdk/qwak/test/it"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	realTimeClient *qwak.RealTimeClient
	ctx            context.Context
	ApiKey         string
	Environment    string
	HttpMock       it.HttpClientMock
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &IntegrationTestSuite{})
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.ApiKey = "jwt-token"

}

func (s *IntegrationTestSuite) TestPredict() {
	// Given
	s.givenQwakClientWithMockedHttpClient()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
	})).Return(it.GetHttpReponse(it.GetAuthResponseWithLongExpiration(), 200), nil).Once()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == "https://models.donald.qwak.ai/v1/otf/predict" &&
			req.Header.Get("authorization") == "Bearer jwt-token"
	})).Return(it.GetHttpReponse(it.GetPredictionResult(), 200), nil).Once()

	// When
	predictionRequest := qwak.NewPredictionRequest("otf").AddFeatureVector(
		qwak.NewFeatureVector().
			WithFeature("State", "PPP").
			WithFeature("Account_Length", 82).
			WithFeature("Area_Code", 53).
			WithFeature("Int'l_Plan", 66).
			WithFeature("VMail_Plan", 85).
			WithFeature("VMail_Message", 23).
			WithFeature("Day_Mins", 1).
			WithFeature("Day_Calls", 9).
			WithFeature("Eve_Mins", 12.0).
			WithFeature("Eve_Calls", 4).
			WithFeature("Night_Mins", 31).
			WithFeature("Night_Calls", 12).
			WithFeature("Intl_Mins", 40).
			WithFeature("Intl_Calls", 15).
			WithFeature("CustServ_Calls", 64).
			WithFeature("Agitation_Level", 9),
	)

	response, err := s.realTimeClient.Predict(predictionRequest)

	// Then
	s.Assert().Equal(nil, err)
	value, err := response.GetSinglePrediction().GetValueAsInt("churn")
	s.Assert().Equal(nil, err)
	s.Assert().Equal(1, value)
	s.HttpMock.Mock.AssertExpectations(s.T())
}

func (s *IntegrationTestSuite) TestAuthenticationOnlyOnceForToken() {
	// Given
	s.givenQwakClientWithMockedHttpClient()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
	})).Return(it.GetHttpReponse(it.GetAuthResponseWithLongExpiration(), 200), nil).Once()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == "https://models.donald.qwak.ai/v1/otf/predict" &&
			req.Header.Get("authorization") == "Bearer jwt-token"
	})).Return(it.GetHttpReponse(it.GetPredictionResult(), 200), nil).Times(3)

	// When
	predictionRequest := qwak.NewPredictionRequest("otf").AddFeatureVector(
		qwak.NewFeatureVector().
			WithFeature("State", "PPP"),
	)

	s.realTimeClient.Predict(predictionRequest)
	s.realTimeClient.Predict(predictionRequest)
	s.realTimeClient.Predict(predictionRequest)

	// Then
	s.HttpMock.Mock.AssertExpectations(s.T())
}

func (s *IntegrationTestSuite) TestAuthenticationRefreshWhenExpired() {
	// Given
	s.givenQwakClientWithMockedHttpClient()

	// Auth requests
	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
	})).Once().Return(it.GetHttpReponse(it.GetAuthResponseWithExpiredDate(), 200), nil).
		On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
		})).Once().Return(it.GetHttpReponse(it.GetAuthResponseWithExpiredDate(), 200), nil).
		On("Do", mock.MatchedBy(func(req *http.Request) bool {
			return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
		})).Once().Return(it.GetHttpReponse(it.GetAuthResponseWithExpiredDate(), 200), nil)

	// Predict requests
	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {

		return req.URL.String() == "https://models.donald.qwak.ai/v1/otf/predict" &&
			req.Header.Get("Authorization") == "Bearer jwt-token"
	})).Return(it.GetHttpReponse(it.GetPredictionResult(), 200), nil).Once()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {

		return req.URL.String() == "https://models.donald.qwak.ai/v1/otf/predict" &&
			req.Header.Get("Authorization") == "Bearer jwt-token"
	})).Return(it.GetHttpReponse(it.GetPredictionResult(), 200), nil).Once()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {

		return req.URL.String() == "https://models.donald.qwak.ai/v1/otf/predict" &&
			req.Header.Get("Authorization") == "Bearer jwt-token"
	})).Return(it.GetHttpReponse(it.GetPredictionResult(), 200), nil).Once()

	// When
	predictionRequest := qwak.NewPredictionRequest("otf").AddFeatureVector(
		qwak.NewFeatureVector().
			WithFeature("State", "PPP"),
	)

	s.realTimeClient.Predict(predictionRequest)
	s.realTimeClient.Predict(predictionRequest)
	s.realTimeClient.Predict(predictionRequest)

	// Then
	s.HttpMock.Mock.AssertExpectations(s.T())
}

func (s *IntegrationTestSuite) TestRetryOnFailure() {
	// Given
	s.givenQwakClientWithMockedHttpClient()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
	})).Return(it.GetHttpReponse(it.GetAuthResponseWithLongExpiration(), 503), nil).Once()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
	})).Return(it.GetHttpReponse(it.GetAuthResponseWithLongExpiration(), 503), nil).Once()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
	})).Return(it.GetHttpReponse(it.GetAuthResponseWithLongExpiration(), 200), nil).Once()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == "https://models.donald.qwak.ai/v1/otf/predict" &&
			req.Header.Get("authorization") == "Bearer jwt-token"
	})).Return(it.GetHttpReponse(it.GetPredictionResult(), 200), nil).Times(3)

	// When
	predictionRequest := qwak.NewPredictionRequest("otf").AddFeatureVector(
		qwak.NewFeatureVector().
			WithFeature("State", "PPP"),
	)

	s.realTimeClient.Predict(predictionRequest)
	s.realTimeClient.Predict(predictionRequest)
	s.realTimeClient.Predict(predictionRequest)

	// Then
	s.HttpMock.Mock.AssertExpectations(s.T())
}

func (s *IntegrationTestSuite) TestAuthFailed() {
	// Given
	s.givenQwakClientWithMockedHttpClient()

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.String() == qwakhttp.DefaultAuthEndpointUri
	})).Return(it.GetHttpReponse(it.GetAuthResponseWithLongExpiration(), 401), nil).Once()

	// When
	predictionRequest := qwak.NewPredictionRequest("otf").AddFeatureVector(
		qwak.NewFeatureVector().
			WithFeature("State", "PPP"),
	)

	_, err := s.realTimeClient.Predict(predictionRequest)

	// Then
	s.Assert().NotEqual(nil, err)
	s.HttpMock.Mock.AssertExpectations(s.T())
}

//func (s *IntegrationTestSuite) givenRealClient() {
//	client, err := qwak.NewRealTimeClient(qwak.RealTimeClientConfig{
//		ApiKey:      s.ApiKey,
//		Environment: "donald",
//		Context:     s.ctx,
//	})
//
//	if err != nil {
//		s.Assert().Fail("client init failed", err)
//	}
//
//	s.realTimeClient = client
//}

func (s *IntegrationTestSuite) givenQwakClientWithMockedHttpClient() {
	client, err := qwak.NewRealTimeClient(qwak.RealTimeClientConfig{
		ApiKey:      s.ApiKey,
		Environment: "donald",
		Context:     s.ctx,
		HttpClient:  &s.HttpMock,
	})

	if err != nil {
		s.Assert().Fail("client init failed", err)
	}

	s.realTimeClient = client
}

func (s *IntegrationTestSuite) TearDownSuite() {
}

func (s *IntegrationTestSuite) SetupTest() {
}

func (s *IntegrationTestSuite) TearDownTest() {
}
