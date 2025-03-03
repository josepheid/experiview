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
	"github.com/josepheid/experiview/api/models"
)

type Handler struct {
	logger    *slog.Logger
	PutItem   func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	tableName string
}

func NewHandler(logger *slog.Logger, putItem func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error), tableName string) (Handler, error) {
	return Handler{
		logger:    logger,
		PutItem:   putItem,
		tableName: tableName,
	}, nil
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request models.CreateExperimentRequest
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

	experimentItem := models.ExperimentItem{
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

	_, err = h.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(h.tableName),
		Item:      data,
	})

	if err != nil {
		h.logger.Error("error putting item", "error", err)
		respond.WithError(w, "error putting item", http.StatusInternalServerError)
		return
	}

	respond.WithJSON(w, experimentItem, http.StatusCreated)
}
