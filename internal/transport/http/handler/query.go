package handler

import (
	"net/http"
	"strconv"
	"strings"

	app "github.com/example/subscriptions-service/internal/service/subscription"
	"github.com/google/uuid"
)

func parseListFilter(r *http.Request) (app.ListFilter, error) {
	q := r.URL.Query()
	f := app.ListFilter{Limit: 50}

	if v := q.Get("user_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			return f, err
		}
		f.UserID = &id
	}

	if v := strings.TrimSpace(q.Get("service_name")); v != "" {
		f.ServiceName = &v
	}

	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return f, err
		}
		f.Limit = n
	}

	if v := q.Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return f, err
		}
		f.Offset = n
	}

	return f, nil
}

func toCreateCommand(req subscriptionRequest) (app.CreateCommand, error) {
	id, err := uuid.Parse(req.UserID)
	if err != nil {
		return app.CreateCommand{}, err
	}

	return app.CreateCommand{ServiceName: req.ServiceName, Price: req.Price, UserID: id, StartDate: req.StartDate, EndDate: req.EndDate}, nil
}

func toUpdateCommand(id uuid.UUID, req subscriptionRequest) (app.UpdateCommand, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return app.UpdateCommand{}, err
	}

	return app.UpdateCommand{ID: id, ServiceName: req.ServiceName, Price: req.Price, UserID: userID, StartDate: req.StartDate, EndDate: req.EndDate}, nil
}
