package getexperiments

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/josepheid/experiview/api/models"
)

func TestGetExperimentsHandlerSuccess(t *testing.T) {
	// Setup mock response
	mockItems := []map[string]types.AttributeValue{
		{
			"PK":          &types.AttributeValueMemberS{Value: "EXPERIMENT#123"},
			"SK":          &types.AttributeValueMemberS{Value: "2023-01-01"},
			"type":        &types.AttributeValueMemberS{Value: "EXPERIMENT"},
			"name":        &types.AttributeValueMemberS{Value: "Test Experiment 1"},
			"description": &types.AttributeValueMemberS{Value: "Description 1"},
		},
		{
			"PK":          &types.AttributeValueMemberS{Value: "EXPERIMENT#456"},
			"SK":          &types.AttributeValueMemberS{Value: "2023-02-01"},
			"type":        &types.AttributeValueMemberS{Value: "EXPERIMENT"},
			"name":        &types.AttributeValueMemberS{Value: "Test Experiment 2"},
			"description": &types.AttributeValueMemberS{Value: "Description 2"},
		},
	}

	// Create mock Query function
	mockQuery := func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
		return &dynamodb.QueryOutput{
			Items: mockItems,
		}, nil
	}

	// Create handler
	logger := slog.Default()
	handler, err := NewHandler(logger, "TestTable", mockQuery)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/experiments", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.ExperimentItem
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 experiments, got %d", len(response))
	}

	if response[0].Name != "Test Experiment 1" {
		t.Errorf("Expected name 'Test Experiment 1', got '%s'", response[0].Name)
	}
}

func TestGetExperimentsWithDateFilters(t *testing.T) {
	// Setup mock response
	mockItems := []map[string]types.AttributeValue{
		{
			"PK":          &types.AttributeValueMemberS{Value: "EXPERIMENT#789"},
			"SK":          &types.AttributeValueMemberS{Value: "2023-03-15"},
			"type":        &types.AttributeValueMemberS{Value: "EXPERIMENT"},
			"name":        &types.AttributeValueMemberS{Value: "Date Filtered Experiment"},
			"description": &types.AttributeValueMemberS{Value: "Description with date filter"},
		},
	}

	// Create mock Query function
	mockQuery := func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
		// Validate that the date filters are being applied correctly
		// This is a more thorough test that verifies the query parameters
		if params.TableName == nil || *params.TableName != "TestTable" {
			t.Errorf("Expected table name 'TestTable', got %v", *params.TableName)
		}

		if params.IndexName == nil || *params.IndexName != "allExperimentsIndex" {
			t.Errorf("Expected index name 'allExperimentsIndex', got %v", *params.IndexName)
		}

		// We can't easily verify the expression attributes without parsing them,
		// but we can check that they exist
		if len(params.ExpressionAttributeNames) == 0 || len(params.ExpressionAttributeValues) == 0 {
			t.Error("Missing expression attributes")
		}

		return &dynamodb.QueryOutput{
			Items: mockItems,
		}, nil
	}

	// Create handler
	logger := slog.Default()
	handler, err := NewHandler(logger, "TestTable", mockQuery)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Create request with date filters
	req := httptest.NewRequest(http.MethodGet, "/experiments?startDate=2023-01-01&endDate=2023-12-31", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.ExperimentItem
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 experiment, got %d", len(response))
	}
}

func TestGetExperimentsWithNameFilter(t *testing.T) {
	// Setup mock response with multiple items
	mockItems := []map[string]types.AttributeValue{
		{
			"PK":          &types.AttributeValueMemberS{Value: "EXPERIMENT#123"},
			"SK":          &types.AttributeValueMemberS{Value: "2023-01-01"},
			"type":        &types.AttributeValueMemberS{Value: "EXPERIMENT"},
			"name":        &types.AttributeValueMemberS{Value: "Apple Experiment"},
			"description": &types.AttributeValueMemberS{Value: "Description 1"},
		},
		{
			"PK":          &types.AttributeValueMemberS{Value: "EXPERIMENT#456"},
			"SK":          &types.AttributeValueMemberS{Value: "2023-02-01"},
			"type":        &types.AttributeValueMemberS{Value: "EXPERIMENT"},
			"name":        &types.AttributeValueMemberS{Value: "Banana Experiment"},
			"description": &types.AttributeValueMemberS{Value: "Description 2"},
		},
	}

	// Create mock Query function
	mockQuery := func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
		return &dynamodb.QueryOutput{
			Items: mockItems,
		}, nil
	}

	// Create handler
	logger := slog.Default()
	handler, err := NewHandler(logger, "TestTable", mockQuery)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Create request with name filter
	req := httptest.NewRequest(http.MethodGet, "/experiments?filter=apple", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.ExperimentItem
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 filtered experiment, got %d", len(response))
	}

	if len(response) > 0 && response[0].Name != "Apple Experiment" {
		t.Errorf("Expected name 'Apple Experiment', got '%s'", response[0].Name)
	}
}

func TestGetExperimentsQueryError(t *testing.T) {
	// Create mock Query function that returns an error
	mockQuery := func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
		return nil, errors.New("DynamoDB error")
	}

	// Create handler
	logger := slog.Default()
	handler, err := NewHandler(logger, "TestTable", mockQuery)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/experiments", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok || message != "error querying all experiments" {
		t.Errorf("Expected message 'error querying all experiments', got %v", response["message"])
	}

	statusCode, ok := response["statusCode"].(float64)
	if !ok || int(statusCode) != http.StatusInternalServerError {
		t.Errorf("Expected statusCode %d, got %v", http.StatusInternalServerError, response["statusCode"])
	}
}

func TestGetExperimentsWithStartDateOnly(t *testing.T) {
	// Setup mock response
	mockItems := []map[string]types.AttributeValue{
		{
			"PK":          &types.AttributeValueMemberS{Value: "EXPERIMENT#123"},
			"SK":          &types.AttributeValueMemberS{Value: "2023-05-01"},
			"type":        &types.AttributeValueMemberS{Value: "EXPERIMENT"},
			"name":        &types.AttributeValueMemberS{Value: "Recent Experiment"},
			"description": &types.AttributeValueMemberS{Value: "After start date"},
		},
	}

	// Create mock Query function that validates parameters
	mockQuery := func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
		// Here we could add validation for the startDate being correctly applied
		// to the key condition expression, but it's complex to parse
		return &dynamodb.QueryOutput{
			Items: mockItems,
		}, nil
	}

	// Create handler
	logger := slog.Default()
	handler, err := NewHandler(logger, "TestTable", mockQuery)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Create request with only start date
	req := httptest.NewRequest(http.MethodGet, "/experiments?startDate=2023-01-01", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.ExperimentItem
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 experiment, got %d", len(response))
	}
}

func TestGetExperimentsUnmarshalError(t *testing.T) {
	// Setup invalid items that will cause an unmarshal error
	invalidItems := []map[string]types.AttributeValue{
		{
			"PK":   &types.AttributeValueMemberS{Value: "EXPERIMENT#123"},
			"SK":   &types.AttributeValueMemberS{Value: "2023-01-01"},
			"type": &types.AttributeValueMemberS{Value: "EXPERIMENT"},
			// Missing required fields or invalid types for unmarshaling
			"name": &types.AttributeValueMemberBOOL{Value: true}, // Expect string, got bool
		},
	}

	// Create mock Query function
	mockQuery := func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
		return &dynamodb.QueryOutput{
			Items: invalidItems,
		}, nil
	}

	// Create handler
	logger := slog.Default()
	handler, err := NewHandler(logger, "TestTable", mockQuery)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/experiments", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	message, ok := response["message"].(string)
	if !ok || message != "error unmarshalling results" {
		t.Errorf("Expected message 'error unmarshalling results', got %v", response["message"])
	}
}
