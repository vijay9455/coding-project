package endpoints

import (
	"calendly/app/middleware"
	"calendly/app/models"
	"calendly/app/repository"
	"calendly/app/services"
	"calendly/lib/db"
	"calendly/lib/logger"
	"calendly/lib/web"
	"context"
	"errors"
)

func CreateMeeting(request *web.Request) (*web.JSONResponse, web.ErrorInterface) {
	currentUser := middleware.GetCurrentUser(request.Context())
	if currentUser == nil {
		logger.Error(request.Context(), "should not have reached here. current_user missing in authenticated endpoint", nil)
		return nil, web.ErrInternalServerError("something went wrong")
	}

	var params services.CreateMeetingParams
	err := request.ValidateBodyToStruct(&params)
	if err != nil {
		logger.Error(request.Context(), "error validation", map[string]any{"error": err})
		return nil, web.ErrValidationFailed(err.Error())
	}

	meeting, err := services.NewMeetingCreateService().Create(request.Context(), currentUser, &params)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return nil, web.ErrValidationFailed("unable to find user with given email")
		}

		if errors.Is(err, services.ErrNoFreeSlot) {
			return nil, web.ErrValidationFailed(err.Error())
		}

		logger.Error(request.Context(), "error while creating meeting", map[string]any{"error": err})
		return nil, web.ErrInternalServerError("something went wrong")
	}

	response := web.JSONResponse(buildMeeting(request.Context(), meeting))
	return &response, nil
}

func buildMeeting(ctx context.Context, meeting *models.Meeting) map[string]any {
	return map[string]any{
		"meeting": map[string]any{
			"id":                  meeting.ID,
			"title":               meeting.Title,
			"meeting_description": meeting.MeetingDescription,
			"start_time":          meeting.StartTime,
			"end_time":            meeting.EndTime,

			"participants": buildParticipants(ctx, meeting.MeetingParticipants),
		},
	}
}

func buildParticipants(ctx context.Context, participants []*models.MeetingParticipant) []map[string]any {
	userRepo := repository.NewUserRepository()
	participantsList := make([]map[string]any, len(participants))
	for idx, participant := range participants {
		user, _ := userRepo.GetByID(ctx, db.Get(), participant.UserID)
		participantsList[idx] = map[string]any{"accept_status": participant.AcceptStatus, "email": user.Email}
	}
	return participantsList
}
