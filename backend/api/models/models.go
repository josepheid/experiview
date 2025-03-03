package models

type CreateExperimentRequest struct {
	Name        string `json:"name"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

type ExperimentItem struct {
	PK          string `dynamodbav:"PK" json:"id"`   // ID
	SK          string `dynamodbav:"SK" json:"date"` // Date
	Name        string `dynamodbav:"name" json:"name"`
	Description string `dynamodbav:"description" json:"description"`
	Type        string `dynamodbav:"type" json:"type"`
	CreatedAt   string `dynamodbav:"createdAt" json:"createdAt"`
	UpdatedAt   string `dynamodbav:"updatedAt" json:"updatedAt"`
}
