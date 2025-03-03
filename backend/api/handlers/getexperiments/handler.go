package getexperiments

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/josepheid/experiview/api/internal/respond"
	"github.com/josepheid/experiview/api/models"
)

type Handler struct {
	logger    *slog.Logger
	tableName string
	ddbc      *dynamodb.Client
}

func NewHandler(logger *slog.Logger, tableName string, ddbc *dynamodb.Client) (Handler, error) {
	return Handler{
		logger:    logger,
		tableName: tableName,
		ddbc:      ddbc,
	}, nil
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	experiments := []models.ExperimentItem{}

	// Extract date filters from query params
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")

	// Build KeyConditionExpression (this replaces the incorrect FilterExpression usage)
	keyCondition := expression.Key("type").Equal(expression.Value("EXPERIMENT"))

	if startDate != "" && endDate != "" {
		keyCondition = keyCondition.And(expression.Key("SK").Between(expression.Value(startDate), expression.Value(endDate)))
	} else if startDate != "" {
		keyCondition = keyCondition.And(expression.Key("SK").GreaterThanEqual(expression.Value(startDate)))
	} else if endDate != "" {
		keyCondition = keyCondition.And(expression.Key("SK").LessThanEqual(expression.Value(endDate)))
	}

	// Build expression
	expr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		h.logger.Error("error building expression", "error", err)
		respond.WithError(w, "error building expression", http.StatusInternalServerError)
		return
	}

	// Query DynamoDB
	data, err := h.ddbc.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(h.tableName),
		IndexName:                 aws.String("allExperimentsIndex"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		h.logger.Error("error querying all experiments", "error", err)
		respond.WithError(w, "error querying all experiments", http.StatusInternalServerError)
		return
	}

	// Unmarshal results
	err = attributevalue.UnmarshalListOfMaps(data.Items, &experiments)
	if err != nil {
		h.logger.Error("error unmarshalling results", "error", err)
		respond.WithError(w, "error unmarshalling results", http.StatusInternalServerError)
		return
	}

	// Apply Name Filtering (Optional, Done in Go)
	filterRes := r.URL.Query().Get("filter")
	h.logger.Info("incoming filter " + filterRes)
	filtered := []models.ExperimentItem{}
	if filterRes != "" {
		for _, exp := range experiments {
			h.logger.Info("incoming name: " + exp.Name)
			if strings.Contains(strings.ToLower(exp.Name), strings.ToLower(filterRes)) {
				filtered = append(filtered, exp)
			}
		}
	}

	// Return filtered results if available
	if len(filtered) != 0 {
		respond.WithJSON(w, filtered, http.StatusOK)
		return
	}

	// Return all results
	respond.WithJSON(w, experiments, http.StatusOK)
}
