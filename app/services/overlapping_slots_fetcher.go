package services

import (
	"calendly/app/models"
	"calendly/lib/logger"
	"context"
	"time"
)

type OverlappingSlotsParams struct {
	Email *string          `json:"email" validate:"required,email"`
	Date  *models.DateOnly `json:"date" validate:"required"`
}

type OverlappingSlotFetcherInterface interface {
	Fetch(ctx context.Context, userA, userB *models.User, date models.DateOnly) ([]*slot, error)
}

func NewOverlappingSlotFetcher() OverlappingSlotFetcherInterface {
	return &overlappingSlotFetcherSvc{
		availabilityFetcher: NewAvailableSlotFetcherService(),
	}
}

type overlappingSlotFetcherSvc struct {
	availabilityFetcher AvailableSlotFetcherInterface
}

func (svc *overlappingSlotFetcherSvc) Fetch(ctx context.Context, userA, userB *models.User, date models.DateOnly) ([]*slot, error) {
	userASlotDetail, err := svc.availabilityFetcher.Fetch(ctx, userA, date, date.Add(24*time.Hour))
	if err != nil {
		logger.Error(ctx, "error while fetching user available slots", map[string]any{"error": err})
		return nil, err
	}

	userBSlotDetail, err := svc.availabilityFetcher.Fetch(ctx, userB, date, date.Add(24*time.Hour))
	if err != nil {
		logger.Error(ctx, "error while fetching user available slots", map[string]any{"error": err})
		return nil, err
	}

	return svc.fetchOverlappingSlots(ctx, userASlotDetail[0].Slots, userBSlotDetail[0].Slots), nil
}

func (svc *overlappingSlotFetcherSvc) fetchOverlappingSlots(ctx context.Context, userASlots, userBSlots []*slot) []*slot {
	if len(userASlots) <= 0 || len(userBSlots) <= 0 {
		return nil
	}

	var slots []*slot
	i, j := 0, 0

	for i < len(userASlots) && j < len(userBSlots) {
		start := maxTime(userASlots[i].StartTime, userBSlots[j].StartTime)
		end := minTime(userASlots[i].EndTime, userBSlots[j].EndTime)

		if start.Before(end) {
			slots = append(slots, &slot{StartTime: start, EndTime: end})
		}

		if userASlots[i].EndTime.Before(userBSlots[j].EndTime) {
			i++
		} else {
			j++
		}
	}

	return slots
}

func maxTime(s, e models.TimeOnly) models.TimeOnly {
	if s.After(e) {
		return s
	}
	return e
}

func minTime(s, e models.TimeOnly) models.TimeOnly {
	if s.Before(e) {
		return s
	}
	return e
}
