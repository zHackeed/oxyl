package models

const (
	JSONData  ValidationType = "JSON_DATA"
	URIData   ValidationType = "URI_DATA"
	QueryData ValidationType = "QUERY_DATA"
	MixedData ValidationType = "MIXED_DATA" // JSON + Query params
)

type ValidationType string

type RequestRequirements struct {
	validationType ValidationType
	model          interface{}
}

func NewRequestRequirements(validationType ValidationType, model interface{}) *RequestRequirements {
	return &RequestRequirements{
		validationType: validationType,
		model:          model,
	}
}

func (r *RequestRequirements) GetValidationType() ValidationType {
	return r.validationType
}

func (r *RequestRequirements) GetModel() interface{} {
	return r.model
}
