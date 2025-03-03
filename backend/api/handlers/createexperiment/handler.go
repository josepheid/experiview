package createexperiment

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/josepheid/experiview/api/internal/respond"
)

type Handler struct {
	logger    *slog.Logger
	ddbc      *dynamodb.Client
	tableName string
}

type CreateExperimentRequest struct {
	Name        string `json:"name"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

type ExperimentItem struct {
	PK          string `dynamodbav:"PK" json:"PK"` // ID
	SK          string `dynamodbav:"SK" json:"SK"` // Date
	Name        string `dynamodbav:"name" json:"name"`
	Description string `dynamodbav:"description" json:"description"`
	Type        string `dynamodbav:"type" json:"type"`
	CreatedAt   string `dynamodbav:"createdAt" json:"createdAt"`
	UpdatedAt   string `dynamodbav:"updatedAt" json:"updatedAt"`
}

type CheckoutSessionResponse struct {
	URL string `json:"url"`
}

func NewHandler(logger *slog.Logger, ddbc *dynamodb.Client, tableName string) (Handler, error) {
	return Handler{
		logger:    logger,
		ddbc:      ddbc,
		tableName: tableName,
	}, nil
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request CreateExperimentRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	h.logger.Info("Incoming request", "requestBody", request)
	if err != nil {
		h.logger.Error("error decoding request body", "error", err)
		respond.WithError(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	experimentID := uuid.New()

	now := time.Now()
	createdAt, updatedAt := now, now

	experimentItem := ExperimentItem{
		PK:          experimentID.String(),
		SK:          request.Date,
		Name:        request.Name,
		Description: request.Description,
		Type:        "EXPERIMENT",
		CreatedAt:   createdAt.Format(time.RFC3339),
		UpdatedAt:   updatedAt.Format(time.RFC3339),
	}

	data, err := attributevalue.MarshalMap(experimentItem)

	if err != nil {
		h.logger.Error("error marshalling job item", "error", err)
		respond.WithError(w, "error marshalling job item", http.StatusInternalServerError)
		return
	}

	_, err = h.ddbc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(h.tableName),
		Item:      data,
	})

	if err != nil {
		h.logger.Error("error putting item", "error", err)
		respond.WithError(w, "error putting item", http.StatusInternalServerError)
		return
	}

	respond.WithJSON(w, nil, http.StatusCreated)
}
