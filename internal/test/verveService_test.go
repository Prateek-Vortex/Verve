package test

import (
	"Verve/internal/model/entity"
	"Verve/internal/model/request"
	"Verve/internal/service"
	"context"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockVerveRepository struct {
	mock.Mock
}

func (m *MockVerveRepository) Save(ctx context.Context, entity entity.VerveEntity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockVerveRepository) GetUniqueCount(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return int64(args.Int(0)), args.Error(1)
}

// Mock RestClient
type MockRestClient struct {
	mock.Mock
}

func (m *MockRestClient) Post(path string, body interface{}) (*http.Response, error) {
	args := m.Called(path, body)
	return &http.Response{StatusCode: http.StatusOK}, args.Error(1)
}

func (m *MockRestClient) Get(path string) (*http.Response, error) {
	args := m.Called(path)
	return &http.Response{StatusCode: http.StatusOK}, args.Error(1)
}

func (m *MockRestClient) Put(path string, body interface{}) (*http.Response, error) {
	args := m.Called(path, body)
	return &http.Response{StatusCode: http.StatusOK}, args.Error(1)
}

func (m *MockRestClient) Delete(path string) (*http.Response, error) {
	args := m.Called(path)
	return &http.Response{StatusCode: http.StatusOK}, args.Error(1)
}

// Mock Event
type MockEvent struct {
	mock.Mock
}

func (m *MockEvent) Publish(ctx context.Context, topic string, message interface{}) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func (m *MockEvent) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	args := m.Called(ctx, topic, handler)
	return args.Error(0)
}

func (m *MockEvent) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestSaveAndPost(t *testing.T) {
	// Setup
	mockRepo := new(MockVerveRepository)
	mockRestClient := new(MockRestClient)
	mockEvent := new(MockEvent)
	logger := slog.Default()

	service := service.NewImplVerveService(mockRepo, mockRestClient, logger, mockEvent)

	// Test case 1: Successful save and post
	t.Run("successful save and post", func(t *testing.T) {
		ctx := context.Background()
		req := request.VerveRequest{
			Id:  "123",
			Url: "http://test.com",
		}

		mockRepo.On("Save", ctx, mock.AnythingOfType("entity.VerveEntity")).Return(nil)
		mockRepo.On("GetUniqueCount", mock.Anything).Return(1, nil)
		mockRestClient.On("Post", req.Url, mock.Anything).Return(req.Url, nil)

		err := service.SaveAndPost(ctx, req)
		assert.NoError(t, err)

		// Wait for goroutine to complete
		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
		mockRestClient.AssertExpectations(t)
	})
}

func TestLogUniqueCountEveryMinute(t *testing.T) {
	// Setup
	mockRepo := new(MockVerveRepository)
	mockRestClient := new(MockRestClient)
	mockEvent := new(MockEvent)
	logger := slog.Default()

	service := service.NewImplVerveService(mockRepo, mockRestClient, logger, mockEvent)

	t.Run("logs count successfully", func(t *testing.T) {
		// Create context with shorter timeout for testing
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Setup mock expectations
		mockRepo.On("GetUniqueCount", mock.Anything).Return(5, nil).Maybe()

		// Create done channel to signal test completion
		done := make(chan bool)

		// Start service in goroutine
		go func() {
			service.LogUniqueCountEveryMinute(ctx)
			done <- true
		}()

		// Wait for either context cancellation or test completion
		select {
		case <-done:
			// Test completed
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Test timed out")
		}

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestSendUniqueCountEveryMinute(t *testing.T) {
	// Setup
	mockRepo := new(MockVerveRepository)
	mockRestClient := new(MockRestClient)
	mockEvent := new(MockEvent)
	logger := slog.Default()

	service := service.NewImplVerveService(mockRepo, mockRestClient, logger, mockEvent)

	t.Run("sends count successfully", func(t *testing.T) {
		// Create shorter context for testing
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Setup mock expectations
		mockRepo.On("GetUniqueCount", mock.Anything).Return(5, nil).Maybe()
		mockEvent.On("Publish", mock.Anything, "unique_count", mock.Anything).Return(nil).Maybe()

		// Channel to track test completion
		done := make(chan bool)

		// Start service in goroutine
		go func() {
			service.SendUniqueCountEveryMinute(ctx)
			done <- true
		}()

		// Wait for completion or timeout
		select {
		case <-done:
			// Test completed normally
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Test timed out")
		}

		// Verify expectations were met
		mockRepo.AssertExpectations(t)
		mockEvent.AssertExpectations(t)
	})
}
