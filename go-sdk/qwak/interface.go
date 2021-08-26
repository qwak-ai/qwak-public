package qwak


type RealTimeClient interface {
	Predict(predictionRequst *PredictionRequest) (*PredictionResponse, error)
}
