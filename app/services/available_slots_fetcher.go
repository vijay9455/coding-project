package services

import (
	"calendly/app/models"
	"calendly/app/repository"
	"calendly/lib/db"
	"calendly/lib/logger"
	"context"
	"time"
)

type AvailableSlotFetcherInterface interface {
	Fetch(ctx context.Context, user *models.User, startDate, endDate models.DateOnly) ([]response, error)
}

type availableSlotFetcherSvc struct {
	uaRepo      repository.UserAvailabilityRepositoryInterface
	meetingRepo repository.MeetingRepositoryInterface
}

type AvailableSlotsParams struct {
	StartDate *models.DateOnly `json:"start_date" validate:"required"`
	EndDate   *models.DateOnly `json:"end_date" validate:"required"`
}

type slot struct {
	StartTime models.TimeOnly `json:"start_time"`
	EndTime   models.TimeOnly `json:"end_time"`
}

type response struct {
	Date      models.DateOnly `json:"date"`
	Available bool            `json:"available"`
	Slots     []*slot         `json:"slots"`
}

func NewAvailableSlotFetcherService() AvailableSlotFetcherInterface {
	return &availableSlotFetcherSvc{
		uaRepo:      repository.NewUserAvailabilityRepository(),
		meetingRepo: repository.NewMeetingRepository(),
	}
}

func (svc *availableSlotFetcherSvc) Fetch(ctx context.Context, user *models.User, startDate, endDate models.DateOnly) ([]response, error) {
	responses := []response{}
	availabilityMap, err := svc.buildAvailabilityMap(ctx, user)
	if err != nil {
		logger.Error(ctx, "error while availability MeetingMap", nil)
		return nil, err
	}

	meetingMap, err := svc.buildMeetingMap(ctx, user, startDate, endDate)
	if err != nil {
		logger.Error(ctx, "error while building MeetingMap", nil)
		return nil, err
	}

	for d := startDate; d.Before(endDate.Time); d = d.Add(24 * time.Hour) {
		if d.Time.Add(24 * time.Hour).Before(time.Now()) {
			responses = append(responses, response{Date: d, Available: false})
		} else {
			availabilities := availabilityMap[int64(d.Time.Weekday())]
			slots := svc.buildAvailableSlots(ctx, d, availabilities, meetingMap[d.Format(models.DateFormat)])
			responses = append(responses, response{Date: d, Available: len(slots) > 0, Slots: slots})
		}
	}
	return responses, nil
}

func (svc *availableSlotFetcherSvc) buildAvailabilityMap(ctx context.Context, user *models.User) (map[int64][]*models.UserAvailability, error) {
	availabilities, err := svc.uaRepo.GetByUserID(ctx, db.Get(), user.ID)
	if err != nil {
		return nil, err
	}

	availabilityMap := make(map[int64][]*models.UserAvailability)
	for _, availability := range availabilities {
		if _, exist := availabilityMap[availability.DayOfWeek]; !exist {
			availabilityMap[availability.DayOfWeek] = make([]*models.UserAvailability, 0)
		}

		availabilityMap[availability.DayOfWeek] = append(availabilityMap[availability.DayOfWeek], availability)
	}

	return availabilityMap, nil
}

func (svc *availableSlotFetcherSvc) buildMeetingMap(ctx context.Context, user *models.User, startDate, endDate models.DateOnly) (map[string][]*models.Meeting, error) {
	meetings, err := svc.meetingRepo.GetForUser(ctx, db.Get(), user.ID, startDate.Time, endDate.Time.Add(24*time.Hour))
	if err != nil {
		return nil, err
	}

	meetingMap := make(map[string][]*models.Meeting)
	for _, meeting := range meetings {
		dateStr := meeting.StartTime.Format(models.DateFormat)
		if _, exist := meetingMap[dateStr]; !exist {
			meetingMap[dateStr] = make([]*models.Meeting, 0)
		}
		meetingMap[dateStr] = append(meetingMap[dateStr], meeting)

		endDateStr := meeting.EndTime.Format(models.DateFormat)
		if endDateStr != dateStr {
			if _, exist := meetingMap[endDateStr]; !exist {
				meetingMap[endDateStr] = make([]*models.Meeting, 0)
			}
			meetingMap[endDateStr] = append(meetingMap[endDateStr], meeting)
		}
	}
	return meetingMap, nil
}

func (svc *availableSlotFetcherSvc) buildAvailableSlots(_ context.Context, date models.DateOnly, availabilities []*models.UserAvailability, meetings []*models.Meeting) []*slot {
	slots := make([]*slot, 0)
	for _, a := range availabilities {
		availabilityStartTime := time.Date(date.Year(), date.Month(), date.Day(), a.StartTime.Hour(), a.StartTime.Minute(), 0, 0, time.UTC)
		availabilityEndTime := time.Date(date.Year(), date.Month(), date.Day(), a.EndTime.Hour(), a.EndTime.Minute(), 0, 0, time.UTC)
		for _, meeting := range meetings {
			if availabilityStartTime.Before(meeting.StartTime) {
				slotEndTime := min(availabilityEndTime, meeting.StartTime)
				slots = append(slots, &slot{
					StartTime: models.TimeOnly{Time: time.Date(0, 0, 0, availabilityStartTime.Hour(), availabilityStartTime.Minute(), 0, 0, time.UTC)},
					EndTime:   models.TimeOnly{Time: time.Date(0, 0, 0, slotEndTime.Hour(), slotEndTime.Minute(), 0, 0, time.UTC)},
				})

				availabilityStartTime = meeting.EndTime
			} else {
				availabilityStartTime = meeting.EndTime
			}
		}

		if availabilityStartTime.Before(availabilityEndTime) {
			slots = append(slots, &slot{
				StartTime: models.TimeOnly{Time: time.Date(0, 0, 0, availabilityStartTime.Hour(), availabilityStartTime.Minute(), 0, 0, time.UTC)},
				EndTime:   models.TimeOnly{Time: time.Date(0, 0, 0, availabilityEndTime.Hour(), availabilityEndTime.Minute(), 0, 0, time.UTC)},
			})
		}
	}
	return slots
}

func min(s, e time.Time) time.Time {
	if s.Before(e) {
		return s
	}
	return e
}
