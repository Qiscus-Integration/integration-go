package resolver

import (
	"context"
	"fmt"
	"integration-go/internal/entity"
	"integration-go/internal/resolver/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errUnexpected = fmt.Errorf("unexpected")

func TestResolvedOmnichannelRoom(t *testing.T) {
	mockRoomRepo := mocks.NewRoomRepository(t)
	mockOmni := mocks.NewOmnichannel(t)

	t.Run("error fetch rooms", func(t *testing.T) {
		mockRoomRepo.EXPECT().Fetch(mock.Anything).Return(nil, errUnexpected).Once()

		svc := Service{
			roomRepo: mockRoomRepo,
			omni:     mockOmni,
		}

		err := svc.ResolvedOmnichannelRoom(context.Background())
		assert.Equal(t, fmt.Errorf("failed to fetch rooms: %w", errUnexpected), err)

		mockRoomRepo.AssertExpectations(t)
		mockOmni.AssertExpectations(t)
	})

	t.Run("skip room less than 10 minutes", func(t *testing.T) {
		rooms := []*entity.Room{
			{
				MultichannelRoomID: "room-123",
				CreatedAt:          time.Now().Add(-5 * time.Minute),
			},
		}

		mockRoomRepo.EXPECT().Fetch(mock.Anything).Return(rooms, nil).Once()

		svc := Service{
			roomRepo: mockRoomRepo,
			omni:     mockOmni,
		}

		err := svc.ResolvedOmnichannelRoom(context.Background())
		assert.Nil(t, err)

		mockRoomRepo.AssertExpectations(t)
		mockOmni.AssertExpectations(t)
	})

	t.Run("error resolved room but continue process", func(t *testing.T) {
		rooms := []*entity.Room{
			{
				MultichannelRoomID: "room-123",
				CreatedAt:          time.Now().Add(-15 * time.Minute),
			},
			{
				MultichannelRoomID: "room-456",
				CreatedAt:          time.Now().Add(-20 * time.Minute),
			},
		}

		mockRoomRepo.EXPECT().Fetch(mock.Anything).Return(rooms, nil).Once()

		// First room
		mockOmni.EXPECT().ResolvedRoom(mock.Anything, "room-123").Return(errUnexpected).Once()

		// Second room
		mockOmni.EXPECT().ResolvedRoom(mock.Anything, "room-456").Return(nil).Once()
		mockRoomRepo.EXPECT().DeleteBy(mock.Anything, map[string]interface{}{
			"multichannel_room_id": "room-456",
		}).Return(nil).Once()

		svc := Service{
			roomRepo: mockRoomRepo,
			omni:     mockOmni,
		}

		err := svc.ResolvedOmnichannelRoom(context.Background())
		assert.Nil(t, err)

		mockRoomRepo.AssertExpectations(t)
		mockOmni.AssertExpectations(t)
	})

	t.Run("error delete room but continue process", func(t *testing.T) {
		rooms := []*entity.Room{
			{
				MultichannelRoomID: "room-123",
				CreatedAt:          time.Now().Add(-15 * time.Minute),
			},
			{
				MultichannelRoomID: "room-456",
				CreatedAt:          time.Now().Add(-20 * time.Minute),
			},
		}

		mockRoomRepo.EXPECT().Fetch(mock.Anything).Return(rooms, nil).Once()

		// First room
		mockOmni.EXPECT().ResolvedRoom(mock.Anything, "room-123").Return(nil).Once()
		mockRoomRepo.EXPECT().DeleteBy(mock.Anything, map[string]interface{}{
			"multichannel_room_id": "room-123",
		}).Return(errUnexpected).Once()

		// Second room
		mockOmni.EXPECT().ResolvedRoom(mock.Anything, "room-456").Return(nil).Once()
		mockRoomRepo.EXPECT().DeleteBy(mock.Anything, map[string]interface{}{
			"multichannel_room_id": "room-456",
		}).Return(nil).Once()

		svc := Service{
			roomRepo: mockRoomRepo,
			omni:     mockOmni,
		}

		err := svc.ResolvedOmnichannelRoom(context.Background())
		assert.Nil(t, err)

		mockRoomRepo.AssertExpectations(t)
		mockOmni.AssertExpectations(t)
	})

	t.Run("success resolve all rooms", func(t *testing.T) {
		rooms := []*entity.Room{
			{
				MultichannelRoomID: "room-123",
				CreatedAt:          time.Now().Add(-15 * time.Minute),
			},
			{
				MultichannelRoomID: "room-456",
				CreatedAt:          time.Now().Add(-20 * time.Minute),
			},
		}

		mockRoomRepo.EXPECT().Fetch(mock.Anything).Return(rooms, nil).Once()

		// First room
		mockOmni.EXPECT().ResolvedRoom(mock.Anything, "room-123").Return(nil).Once()
		mockRoomRepo.EXPECT().DeleteBy(mock.Anything, map[string]interface{}{
			"multichannel_room_id": "room-123",
		}).Return(nil).Once()

		// Second room
		mockOmni.EXPECT().ResolvedRoom(mock.Anything, "room-456").Return(nil).Once()
		mockRoomRepo.EXPECT().DeleteBy(mock.Anything, map[string]interface{}{
			"multichannel_room_id": "room-456",
		}).Return(nil).Once()

		svc := Service{
			roomRepo: mockRoomRepo,
			omni:     mockOmni,
		}

		err := svc.ResolvedOmnichannelRoom(context.Background())
		assert.Nil(t, err)

		mockRoomRepo.AssertExpectations(t)
		mockOmni.AssertExpectations(t)
	})
}
