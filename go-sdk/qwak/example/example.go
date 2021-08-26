package qwak

import (
	"log"
	"qwak.ai/inference-sdk"
)

func exampleOfUser() {
	realTimeClient := qwak.NewRealTimeClient(qwak.RealTimeClientOptions{
		ApiKey:      "your api key",
		Environment: "environment name",
	})


	predictionRequest := &qwak.PredictionRequest{
		featureVectors: []*qwak.FeatureVector{
			&qwak.FeatureVector{
				Features: map[string]interface{}{
					"feature_a": 5,
					"feature_b": "USA",
					"feature_c": 8.5,
				},
			},
			&qwak.FeatureVector{
				Features: map[string]interface{}{
					"feature_a": 4,
					"feature_b": "UK",
					"feature_c": 2.7,
				},
			},
		},
	}

	response, err := realTimeClient.Predict(predictionRequest)

	if err != nil {
		log.Printf("failed to infer Qwak real time client. error: %v", err.Error())
		return
	}

	log.Printf("response from qwak client", response)

}
