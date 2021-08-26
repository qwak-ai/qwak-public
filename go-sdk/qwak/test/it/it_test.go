package it_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	qwak "qwak.ai/inference-sdk"
	qwakhttp "qwak.ai/inference-sdk/http"
	"qwak.ai/inference-sdk/test/it"
)

type IntegrationTestSuite struct {
	suite.Suite
	realTimeClient qwak.RealTimeClient
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
	s.ApiKey = "6abc183a174c7991d3473f684635b62b@b5a414c542e94c629a8ba28fd4c63d2d"

}

func (s *IntegrationTestSuite) TestPredict() {
	s.givenRealClient()
	_, err := qwakhttp.GetAuthenticationRequest(s.ctx, s.ApiKey)

	if err != nil {
		s.Assertions.Fail(fmt.Sprintf("faile create request: %s", err.Error()))
	}

	s.HttpMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		fmt.Println(req.URL.String(), qwakhttp.DEFAULT_AUTH_ENDPOINT_URI)
		return req.URL.String() == qwakhttp.DEFAULT_AUTH_ENDPOINT_URI
	})).Return("hi", nil).Once()

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
	s.Assert().Equal(nil, err)
	value, err := response.GetSinglePrediction().GetValueAsInt("churn")
	s.Assert().Equal(nil, err)
	s.Assert().Equal(1, value)

	fmt.Println(response)
}

func (s *IntegrationTestSuite) givenRealClient() {
	client, err := qwak.NewRealTimeClient(qwak.RealTimeClientConfig{
		ApiKey:      s.ApiKey,
		Environment: "donald",
		Context:     s.ctx,
	})

	if err != nil {
		s.Assert().Fail("client init failed", err)
	}

	s.realTimeClient = client
}

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

// func getAuthReponse() *http.Response {
// 	return &http.Response{
// 		Body: ,
// 	}
// }

func (s *IntegrationTestSuite) TearDownSuite() {
}

func (s *IntegrationTestSuite) SetupTest() {
}

func (s *IntegrationTestSuite) TearDownTest() {
}
