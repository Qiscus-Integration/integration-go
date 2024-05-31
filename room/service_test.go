package room

import (
	"context"
	"fmt"
	"integration-go/entity"
	"integration-go/room/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var errUnexpected = fmt.Errorf("unexpected")

func TestGetRoomByID(t *testing.T) {
	mockRepo := mocks.NewRepository(t)

	t.Run("error get room", func(t *testing.T) {
		mockRepo.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errUnexpected).Once()

		svc := Service{repo: mockRepo}
		room, err := svc.GetRoomByID(context.Background(), 1)
		assert.Equal(t, fmt.Errorf("failed to find room: %w", errUnexpected), err)
		assert.Nil(t, room)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error get room - not found", func(t *testing.T) {
		mockRepo.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, gorm.ErrRecordNotFound).Once()

		svc := Service{repo: mockRepo}
		room, err := svc.GetRoomByID(context.Background(), 1)
		assert.Equal(t, &roomError{roomErrorNotFound}, err)
		assert.Nil(t, room)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success get room", func(t *testing.T) {
		mockRepo.On("FindByID", mock.Anything, mock.AnythingOfType("int64")).Return(&entity.Room{ID: 1}, nil).Once()

		svc := Service{repo: mockRepo}
		room, err := svc.GetRoomByID(context.Background(), 1)
		assert.Equal(t, int64(1), room.ID)
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

}
