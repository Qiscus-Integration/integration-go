package room

import (
	"context"
	"fmt"
	"integration-go/internal/entity"
	"integration-go/internal/qismo"
	"integration-go/internal/room/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var errUnexpected = fmt.Errorf("unexpected")

func TestGetRoomByID(t *testing.T) {
	mockRepo := mocks.NewRepository(t)

	t.Run("error get room", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(mock.Anything, mock.AnythingOfType("int64")).Return(nil, errUnexpected).Once()

		svc := Service{repo: mockRepo}
		room, err := svc.GetRoomByID(context.Background(), 1)
		assert.Equal(t, fmt.Errorf("failed to find room: %w", errUnexpected), err)
		assert.Nil(t, room)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error get room - not found", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(mock.Anything, mock.AnythingOfType("int64")).Return(nil, gorm.ErrRecordNotFound).Once()

		svc := Service{repo: mockRepo}
		room, err := svc.GetRoomByID(context.Background(), 1)
		assert.Equal(t, &roomError{roomErrorNotFound}, err)
		assert.Nil(t, room)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success get room", func(t *testing.T) {
		mockRepo.EXPECT().FindByID(mock.Anything, mock.AnythingOfType("int64")).Return(&entity.Room{ID: 1}, nil).Once()

		svc := Service{repo: mockRepo}
		room, err := svc.GetRoomByID(context.Background(), 1)
		assert.Equal(t, int64(1), room.ID)
		assert.Nil(t, err)
		mockRepo.AssertExpectations(t)
	})

}
func TestCreateRoom(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockOmni := mocks.NewOmnichannel(t)

	req := &qismo.WebhookNewSessionRequest{
		IsNewSession: true,
		Payload: struct {
			Room struct {
				ID              string `json:"id"`
				IDStr           string `json:"id_str"`
				IsPublicChannel bool   `json:"is_public_channel"`
				Name            string `json:"name"`
				Options         string `json:"options"`
				Participants    []struct {
					Email string `json:"email"`
				} `json:"participants"`
				RoomAvatar string `json:"room_avatar"`
				TopicID    string `json:"topic_id"`
				TopicIDStr string `json:"topic_id_str"`
				Type       string `json:"type"`
			} `json:"room"`
		}{
			Room: struct {
				ID              string `json:"id"`
				IDStr           string `json:"id_str"`
				IsPublicChannel bool   `json:"is_public_channel"`
				Name            string `json:"name"`
				Options         string `json:"options"`
				Participants    []struct {
					Email string `json:"email"`
				} `json:"participants"`
				RoomAvatar string `json:"room_avatar"`
				TopicID    string `json:"topic_id"`
				TopicIDStr string `json:"topic_id_str"`
				Type       string `json:"type"`
			}{
				IDStr: "room-123",
			},
		},
		WebhookType: "new_session",
	}

	t.Run("error create omnichannel tag", func(t *testing.T) {
		mockOmni.EXPECT().CreateRoomTag(mock.Anything, req.Payload.Room.IDStr, req.Payload.Room.IDStr).
			Return(errUnexpected).Once()

		svc := Service{
			repo: mockRepo,
			omni: mockOmni,
		}

		err := svc.CreateRoom(context.Background(), req)
		assert.Equal(t, fmt.Errorf("failed to create omnichannel tag: %w", errUnexpected), err)

		mockOmni.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error save room", func(t *testing.T) {
		mockOmni.EXPECT().CreateRoomTag(mock.Anything, req.Payload.Room.IDStr, req.Payload.Room.IDStr).
			Return(nil).Once()

		mockRepo.EXPECT().Save(mock.Anything, &entity.Room{
			MultichannelRoomID: req.Payload.Room.IDStr,
		}).Return(errUnexpected).Once()

		svc := Service{
			repo: mockRepo,
			omni: mockOmni,
		}

		err := svc.CreateRoom(context.Background(), req)
		assert.Equal(t, fmt.Errorf("failed to save room: %w", errUnexpected), err)

		mockOmni.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success create room", func(t *testing.T) {
		mockOmni.EXPECT().CreateRoomTag(mock.Anything, req.Payload.Room.IDStr, req.Payload.Room.IDStr).
			Return(nil).Once()

		mockRepo.EXPECT().Save(mock.Anything, &entity.Room{
			MultichannelRoomID: req.Payload.Room.IDStr,
		}).Return(nil).Once()

		svc := Service{
			repo: mockRepo,
			omni: mockOmni,
		}

		err := svc.CreateRoom(context.Background(), req)
		assert.Nil(t, err)

		mockOmni.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
