package services

import (
	"calendly/app/models"
	"calendly/app/repository"
	"calendly/lib/db"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MeetingCreateServiceInterface interface {
	Create(ctx context.Context, owner *models.User, params *CreateMeetingParams) (*models.Meeting, error)
}

func NewMeetingCreateService() MeetingCreateServiceInterface {
	return &meetingCreateSvc{
		meetingRepo:                repository.NewMeetingRepository(),
		userRepo:                   repository.NewUserRepository(),
		overlappingAvailabilitySvc: NewOverlappingSlotFetcher(),
	}
}

type CreateMeetingParams struct {
	Email              *string    `json:"email" validate:"required,email"`
	StartTime          *time.Time `json:"start_time" validate:"required"`
	EndTime            *time.Time `json:"end_time" validate:"required"`
	Title              *string    `json:"title" validate:"required"`
	MeetingDescription *string    `json:"meeting_description"`
}

type meetingCreateSvc struct {
	meetingRepo                repository.MeetingRepositoryInterface
	userRepo                   repository.UserRepositoryInterface
	overlappingAvailabilitySvc OverlappingSlotFetcherInterface
}

func (svc *meetingCreateSvc) Create(ctx context.Context, owner *models.User, params *CreateMeetingParams) (*models.Meeting, error) {
	participant, err := svc.userRepo.GetByEmail(ctx, db.Get(), *params.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if err := svc.validateMeetingTimeAvailability(ctx, owner, participant, *params.StartTime, *params.EndTime); err != nil {
		return nil, err
	}

	return svc.create(ctx, params, owner, participant)
}

func (svc *meetingCreateSvc) create(ctx context.Context, params *CreateMeetingParams, owner, participant *models.User) (*models.Meeting, error) {
	meeting := &models.Meeting{
		Title:              *params.Title,
		MeetingDescription: *params.MeetingDescription,

		StartTime: *params.StartTime,
		EndTime:   *params.EndTime,

		MeetingParticipants: []*models.MeetingParticipant{
			&models.MeetingParticipant{UserID: owner.ID, AcceptStatus: models.Accepted},
			&models.MeetingParticipant{UserID: participant.ID, AcceptStatus: models.MAY_BE},
		},
	}

	if err := svc.meetingRepo.Create(ctx, db.Get(), meeting); err != nil {
		return nil, err
	}

	return meeting, nil
}

func (svc *meetingCreateSvc) validateMeetingTimeAvailability(ctx context.Context, userA, userB *models.User, startTime, endTime time.Time) error {
	availableSlots, err := svc.overlappingAvailabilitySvc.Fetch(ctx, userA, userB,
		models.DateOnly{Time: time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)})
	if err != nil {
		return err
	}

	for _, slot := range availableSlots {
		if svc.isWithinSlot(ctx, startTime, endTime, slot) {
			return nil
		}
	}

	return ErrNoFreeSlot
}

func (svc *meetingCreateSvc) isWithinSlot(ctx context.Context, startTime, endTime time.Time, slot *slot) bool {
	startTimeOnly := models.TimeOnly{Time: time.Date(0, 0, 0, startTime.Hour(), startTime.Minute(), startTime.Second()+1, 0, time.UTC)}
	endTimeOnly := models.TimeOnly{Time: time.Date(0, 0, 0, endTime.Hour(), endTime.Minute(), endTime.Second()-1, 0, time.UTC)}

	return startTimeOnly.After(slot.StartTime) && endTimeOnly.Before(slot.EndTime)
}
