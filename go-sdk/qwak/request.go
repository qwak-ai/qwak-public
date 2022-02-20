package qwak

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/qwak-ai/qwak-public/go-sdk/qwak/http"
)

// PredictionRequest represents is a fluent API to build a prediction request on your model
type PredictionRequest struct {
	ModelId        string
	featuresVector []*FeatureVector
}

// NewPredictionRequest is a constructor of PredictionRequest fluent API
func NewPredictionRequest(modelId string) *PredictionRequest {
	return &PredictionRequest{ModelId: modelId}
}

// AddFeatureVector adding a new feature vector to your prediction request using fluent API
func (ir *PredictionRequest) AddFeatureVector(featureVector *FeatureVector) *PredictionRequest {
	ir.featuresVector = append(ir.featuresVector, featureVector)
	return ir
}

// AddFeatureVectors adding many new feature vector to your prediction request using fluent API
func (ir *PredictionRequest) AddFeatureVectors(featuresVector ...*FeatureVector) *PredictionRequest {
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

// PredictionResponse represnt a response from your model to a prediction request
type PredictionResponse struct {
	predictions []*PredictionResult
}

// GetPredictions is getting a resluts array from response
func (pr *PredictionResponse) GetPredictions() []*PredictionResult {
	return pr.predictions
}

// GetSinglePrediction returns a single result from a prediction response
func (pr *PredictionResponse) GetSinglePrediction() *PredictionResult {
	if len(pr.predictions) > 0 {
		return pr.predictions[0]
	}

	return nil
}

func responseFromRaw(results []byte) (*PredictionResponse, error) {

	var response []map[string]interface{}
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

// PredictionResult respresnts one result in a response for prediction request
type PredictionResult struct {
	valuesMap map[string]interface{}
}

// GetValueAsInt returning the value of column in a result converted to int.
// If conversion failed, an error returned
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

// GetValueAsFloat returning the value of column in a result converted to float.
// If conversion failed, an error returned
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

// GetValueAsString returning the value of column in a result converted to string.
// If conversion failed, an error returned
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

// GetValueAsArrayOfStrings returning the value of column in a result converted to array of strings.
// If conversion failed, an error returned
func (pr *PredictionResult) GetValueAsArrayOfStrings(columnName string) ([]string, error) {
	value, ok := pr.valuesMap[columnName]

	if !ok {
		return nil, errors.New("column is not exists")
	}

	parsedValue, ok := value.([]interface{})

	if !ok {
		return nil, errors.New("column value is not an array")
	}

	var result []string

	for idx, val := range parsedValue {
		parsedString, ok := val.(string)

		if !ok {
			return nil, fmt.Errorf("the value of %s with index %d is not a string", columnName, idx)
		}

		result = append(result, parsedString)
	}

	return result, nil
}

// GetValueAsInterface returning the value of column in a result without any conversion
func (pr *PredictionResult) GetValueAsInterface(columnName string) (interface{}, error) {
	value, ok := pr.valuesMap[columnName]

	if !ok {
		return nil, errors.New("column is not exists")
	}

	return value, nil
}

// FeatureVector represents a vector of features with their name and value
type FeatureVector struct {
	features []*feature
}

// NewFeatureVector is a constructor for FeatureVector with fluent API
func NewFeatureVector() *FeatureVector {
	return &FeatureVector{}
}

// WithFeature set a feature on a FeatureVector
func (fr *FeatureVector) WithFeature(name string, value interface{}) *FeatureVector {
	fr.features = append(fr.features, &feature{
		Name:  name,
		Value: value,
	})

	return fr
}

type feature struct {
	Name  string
	Value interface{}
}
