package createexperiment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/josepheid/experiview/api/models"
)

// To make the tests work, we need to modify the Handler to accept an interface
// rather than a concrete *dynamodb.Client. Since we can't modify the original code,
// we'll use a different approach.

// Instead, we'll create our tests by embedding a dynamodb.Client struct
// and overriding just the methods we need.

type TestDynamoDBClient struct {
	dynamodb.Client
	putItemFunc func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func (m *TestDynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if m.putItemFunc != nil {
		return m.putItemFunc(ctx, params, optFns...)
	}
	return &dynamodb.PutItemOutput{}, nil
}

func TestServeHTTP_Success(t *testing.T) {
	// Setup
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	tableName := "test-table"

	var capturedInput *dynamodb.PutItemInput
	putItemFunc := func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		capturedInput = params
		return &dynamodb.PutItemOutput{}, nil
	}

	handler, _ := NewHandler(logger, putItemFunc, tableName)

	// Create request
	reqBody := models.CreateExperimentRequest{
		Name:        "Test Experiment",
		Description: "This is a test experiment",
		Date:        "2025-03-01",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/experiments", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Verify DynamoDB call
	if capturedInput == nil {
		t.Fatal("Expected PutItem to be called, but it wasn't")
	}

	if aws.ToString(capturedInput.TableName) != tableName {
		t.Errorf("Expected table name %s, got %s", tableName, aws.ToString(capturedInput.TableName))
	}

	// Check for required fields in the item
	requiredFields := []string{"PK", "SK", "name", "description", "type", "createdAt", "updatedAt"}
	for _, field := range requiredFields {
		if _, exists := capturedInput.Item[field]; !exists {
			t.Errorf("Field %s missing from DynamoDB item", field)
		}
	}

	// Verify response
	var response models.ExperimentItem
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Name != reqBody.Name {
		t.Errorf("Expected name %s, got %s", reqBody.Name, response.Name)
	}
	if response.Description != reqBody.Description {
		t.Errorf("Expected description %s, got %s", reqBody.Description, response.Description)
	}
	if response.SK != reqBody.Date {
		t.Errorf("Expected SK %s, got %s", reqBody.Date, response.SK)
	}
	if response.Type != "EXPERIMENT" {
		t.Errorf("Expected Type EXPERIMENT, got %s", response.Type)
	}

	// Verify UUID format for PK
	_, err = uuid.Parse(response.PK)
	if err != nil {
		t.Errorf("PK is not a valid UUID: %s", response.PK)
	}

	// Verify timestamps are in RFC3339 format
	_, err = time.Parse(time.RFC3339, response.CreatedAt)
	if err != nil {
		t.Errorf("CreatedAt is not in RFC3339 format: %s", response.CreatedAt)
	}

	_, err = time.Parse(time.RFC3339, response.UpdatedAt)
	if err != nil {
		t.Errorf("UpdatedAt is not in RFC3339 format: %s", response.UpdatedAt)
	}
}

func TestServeHTTP_InvalidRequestBody(t *testing.T) {
	// Setup
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	putItemFunc := func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return nil, errors.New("DynamoDB error")
	}
	tableName := "test-table"

	handler, _ := NewHandler(logger, putItemFunc, tableName)

	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/experiments", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if message, ok := response["message"].(string); !ok || message != "error decoding request body" {
		t.Errorf("Expected message 'error putting item', got %v", response["message"])
	}
}

func TestServeHTTP_DynamoDBPutItemError(t *testing.T) {
	// Setup
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	putItemFunc := func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return nil, errors.New("DynamoDB error")
	}

	tableName := "test-table"

	handler, _ := NewHandler(logger, putItemFunc, tableName)

	// Create request
	reqBody := models.CreateExperimentRequest{
		Name:        "Test Experiment",
		Description: "This is a test experiment",
		Date:        "2025-03-01",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/experiments", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Assert on the response
	if statusCode, ok := response["statusCode"].(float64); !ok || int(statusCode) != 500 {
		t.Errorf("Expected statusCode 500, got %v", response["statusCode"])
	}

	if message, ok := response["message"].(string); !ok || message != "error putting item" {
		t.Errorf("Expected message 'error putting item', got %v", response["message"])
	}

}

func TestMarshallUnmarshallExperiment(t *testing.T) {
	// Test marshaling an experiment item to DynamoDB attributes
	experimentID := uuid.New()
	now := time.Now()

	experimentItem := models.ExperimentItem{
		PK:          experimentID.String(),
		SK:          "2025-03-01",
		Name:        "Test Experiment",
		Description: "Test Description",
		Type:        "EXPERIMENT",
		CreatedAt:   now.Format(time.RFC3339),
		UpdatedAt:   now.Format(time.RFC3339),
	}

	// Marshal to DynamoDB attributes
	data, err := attributevalue.MarshalMap(experimentItem)
	if err != nil {
		t.Fatalf("Failed to marshal experiment item: %v", err)
	}

	// Check required fields
	requiredFields := []string{"PK", "SK", "name", "description", "type", "createdAt", "updatedAt"}
	for _, field := range requiredFields {
		if _, exists := data[field]; !exists {
			t.Errorf("Field %s missing after marshaling", field)
		}
	}

	// Unmarshal back to struct
	var unmarshalledItem models.ExperimentItem
	err = attributevalue.UnmarshalMap(data, &unmarshalledItem)
	if err != nil {
		t.Fatalf("Failed to unmarshal experiment item: %v", err)
	}

	// Verify fields match
	if unmarshalledItem.PK != experimentItem.PK {
		t.Errorf("PK mismatch: expected %s, got %s", experimentItem.PK, unmarshalledItem.PK)
	}
	if unmarshalledItem.SK != experimentItem.SK {
		t.Errorf("SK mismatch: expected %s, got %s", experimentItem.SK, unmarshalledItem.SK)
	}
	if unmarshalledItem.Name != experimentItem.Name {
		t.Errorf("Name mismatch: expected %s, got %s", experimentItem.Name, unmarshalledItem.Name)
	}
	if unmarshalledItem.Description != experimentItem.Description {
		t.Errorf("Description mismatch: expected %s, got %s", experimentItem.Description, unmarshalledItem.Description)
	}
	if unmarshalledItem.Type != experimentItem.Type {
		t.Errorf("Type mismatch: expected %s, got %s", experimentItem.Type, unmarshalledItem.Type)
	}
	if unmarshalledItem.CreatedAt != experimentItem.CreatedAt {
		t.Errorf("CreatedAt mismatch: expected %s, got %s", experimentItem.CreatedAt, unmarshalledItem.CreatedAt)
	}
	if unmarshalledItem.UpdatedAt != experimentItem.UpdatedAt {
		t.Errorf("UpdatedAt mismatch: expected %s, got %s", experimentItem.UpdatedAt, unmarshalledItem.UpdatedAt)
	}
}

func TestCreateExperimentEndToEnd(t *testing.T) {
	// This test simulates an end-to-end flow
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	tableName := "test-table"

	// Track what was stored
	var storedItem map[string]types.AttributeValue

	putItemFunc := func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		storedItem = params.Item
		return &dynamodb.PutItemOutput{}, nil
	}

	handler, _ := NewHandler(logger, putItemFunc, tableName)

	// Create a valid request with all fields
	reqBody := models.CreateExperimentRequest{
		Name:        "End-to-End Test",
		Description: "Testing complete flow",
		Date:        "2025-03-03",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/experiments", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	handler.ServeHTTP(w, req)

	// Verify successful response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	// Verify what was stored in DynamoDB
	if storedItem == nil {
		t.Fatal("Expected item to be stored in DynamoDB")
	}

	// Unmarshal the stored item
	var storedExperiment models.ExperimentItem
	err := attributevalue.UnmarshalMap(storedItem, &storedExperiment)
	if err != nil {
		t.Fatalf("Failed to unmarshal stored item: %v", err)
	}

	// Verify fields match what was requested
	if storedExperiment.Name != reqBody.Name {
		t.Errorf("Stored name mismatch: expected %s, got %s", reqBody.Name, storedExperiment.Name)
	}
	if storedExperiment.Description != reqBody.Description {
		t.Errorf("Stored description mismatch: expected %s, got %s", reqBody.Description, storedExperiment.Description)
	}
	if storedExperiment.SK != reqBody.Date {
		t.Errorf("Stored SK mismatch: expected %s, got %s", reqBody.Date, storedExperiment.SK)
	}
	if storedExperiment.Type != "EXPERIMENT" {
		t.Errorf("Stored Type mismatch: expected EXPERIMENT, got %s", storedExperiment.Type)
	}

	// Verify response matches what was stored
	var responseExperiment models.ExperimentItem
	err = json.Unmarshal(w.Body.Bytes(), &responseExperiment)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseExperiment.PK != storedExperiment.PK {
		t.Errorf("PK mismatch between response and stored item")
	}
	if responseExperiment.Name != storedExperiment.Name {
		t.Errorf("Name mismatch between response and stored item")
	}
	if responseExperiment.Description != storedExperiment.Description {
		t.Errorf("Description mismatch between response and stored item")
	}
}
