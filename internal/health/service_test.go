package health

import (
	"context"
	"integration-go/internal/entity"
	"integration-go/internal/health/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name                  string
		mockDatabaseError     error
		mockRedisError        error
		expectedDatabaseState entity.HealthState
		expectedRedisState    entity.HealthState
		expectedHealthy       bool
	}{
		{
			name:                  "database and redis both healthy",
			mockDatabaseError:     nil,
			mockRedisError:        nil,
			expectedDatabaseState: entity.HealthStateOK,
			expectedRedisState:    entity.HealthStateOK,
			expectedHealthy:       true,
		},
		{
			name:                  "database unhealthy",
			mockDatabaseError:     assert.AnError,
			mockRedisError:        nil,
			expectedDatabaseState: entity.HealthStateFail,
			expectedRedisState:    entity.HealthStateOK,
			expectedHealthy:       false,
		},
		{
			name:                  "redis unhealthy",
			mockDatabaseError:     nil,
			mockRedisError:        assert.AnError,
			expectedDatabaseState: entity.HealthStateOK,
			expectedRedisState:    entity.HealthStateFail,
			expectedHealthy:       false,
		},
		{
			name:                  "both database and redis unhealthy",
			mockDatabaseError:     assert.AnError,
			mockRedisError:        assert.AnError,
			expectedDatabaseState: entity.HealthStateFail,
			expectedRedisState:    entity.HealthStateFail,
			expectedHealthy:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			mockRepo.On("CheckDatabase", mock.Anything).Return(tt.mockDatabaseError).Once()
			mockRepo.On("CheckRedis", mock.Anything).Return(tt.mockRedisError).Once()

			svc := &Service{repo: mockRepo}
			healthComponent, isHealthy := svc.Check(context.Background())

			assert.Equal(t, tt.expectedDatabaseState, healthComponent.Database)
			assert.Equal(t, tt.expectedRedisState, healthComponent.Redis)
			assert.Equal(t, tt.expectedHealthy, isHealthy)

			mockRepo.AssertExpectations(t)
		})
	}
}
