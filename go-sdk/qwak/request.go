package qwak

import (
	"encoding/json"
	"errors"
	"fmt"

	"qwak.ai/inference-sdk/http"
)

type PredictionRequest struct {
	ModelId        string
	featuresVector []*FeatureVector
}

func NewPredictionRequest(modelId string) *PredictionRequest {
	return &PredictionRequest{ModelId: modelId}
}

func (ir *PredictionRequest) AddFeatureVector(featureVector *FeatureVector) *PredictionRequest {
	ir.featuresVector = append(ir.featuresVector, featureVector)
	return ir
}
func (ir *PredictionRequest) AddFeaturesVector(featuresVector ...*FeatureVector) *PredictionRequest {
	ir.featuresVector = append(ir.featuresVector, featuresVector...)
	return ir
}

func (ir *PredictionRequest) asPandaOrientedDf() http.PandaOrientedDf {

	index := make([]int, len(ir.featuresVector))
	columnNextIdx := 0
	columnsIdxByName := map[string]int{}
	columnsData := make([][]interface{}, len(ir.featuresVector))

	// collect columns names and indeces
	for idx, vector := range ir.featuresVector {
		index[idx] = idx
		for _, feature := range vector.features {
			if _, ok := columnsIdxByName[feature.Name]; !ok {
				columnsIdxByName[feature.Name] = columnNextIdx
				columnNextIdx++
			}
		}
	}

	// collect values
	for idx, vector := range ir.featuresVector {
		columnsData[idx] = make([]interface{}, len(columnsIdxByName))

		for _, feature := range vector.features {
			columnsData[idx][columnsIdxByName[feature.Name]] = feature.Value
		}
	}

	columnsNames := make([]string, len(columnsIdxByName))

	for columnName, columnIdx := range columnsIdxByName {
		columnsNames[columnIdx] = columnName
	}

	return http.NewPandaOrientedDf(columnsNames, index, columnsData)
}

type PredictionResponse struct {
	predictions []*PredictionResult
}

func (pr *PredictionResponse) GetPredictions() []*PredictionResult {
	return pr.predictions
}

func (pr *PredictionResponse) GetSinglePrediction() *PredictionResult {
	if len(pr.predictions) > 0 {
		return pr.predictions[0]
	}

	return nil
}

func responseFromRaw(results []byte) (*PredictionResponse, error) {

	response := []map[string]interface{}{}
	err := json.Unmarshal(results, &response)

	if err != nil {
		return nil, fmt.Errorf("qwak client failed to predict: %s", err.Error())
	}

	predictionResponse := &PredictionResponse{}

	for _, result := range response {
		predictionResponse.predictions = append(predictionResponse.predictions, &PredictionResult{
			valuesMap: result,
		})
	}

	return predictionResponse, nil
}

type PredictionResult struct {
	valuesMap map[string]interface{}
}

func (pr *PredictionResult) GetValueAsInt(columnName string) (int, error) {
	value, ok := pr.valuesMap[columnName]

	if !ok {
		return 0, errors.New("column is not exists")
	}

	parsedValue, ok := value.(float64)

	if !ok {
		return 0, errors.New("column value is not a number")
	}

	return int(parsedValue), nil
}

func (pr *PredictionResult) GetValueAsFloat(columnName string) (float64, error) {
	value, ok := pr.valuesMap[columnName]

	if !ok {
		return 0, errors.New("column is not exists")
	}

	parsedValue, ok := value.(float64)

	if !ok {
		return 0, errors.New("column value is not a float")
	}

	return parsedValue, nil
}

func (pr *PredictionResult) GetValueAsString(columnName string) (string, error) {
	value, ok := pr.valuesMap[columnName]

	if !ok {
		return "", errors.New("column is not exists")
	}

	parsedValue, ok := value.(string)

	if !ok {
		return "", errors.New("column value is not a float")
	}

	return parsedValue, nil
}

type FeatureVector struct {
	features []*Feature
}

func NewFeatureVector() *FeatureVector {
	return &FeatureVector{}
}

func (fr *FeatureVector) WithFeature(name string, value interface{}) (*FeatureVector) {
	fr.features = append(fr.features, &Feature{
		Name: name,
		Value: value,
	})

	return fr
}



type Feature struct {
	Name  string
	Value interface{}
}
