package it

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HttpClientMock struct {
	mock.Mock
}

func (hcm *HttpClientMock) Do(request *http.Request) (*http.Response, error) {
	args := hcm.Mock.MethodCalled("Do", request)
	
	return args.Get(0).(*http.Response), args.Error(1)
}


