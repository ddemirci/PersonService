package handler

import (
	"context"
	"encoding/json"
	"person-service/domain"
	"person-service/ports"
)

type PersonHandler struct {
	service ports.PersonService
}

func NewPersonHandler(service ports.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func (h *PersonHandler) Handle(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {

	method := getMethod(event)

	if method == "POST" {
		return h.handlePost(ctx, event)
	}

	if method == "GET" {
		return h.handleGet(ctx)
	}

	return map[string]interface{}{
		"statusCode": 400,
		"body":       "Unsupported method",
	}, nil
}

func (h *PersonHandler) handlePost(ctx context.Context, event map[string]interface{}) (map[string]interface{}, error) {

	body, ok := event["body"].(string)
	if !ok {
		return response(400, "Invalid body")
	}

	var person domain.Person
	if err := json.Unmarshal([]byte(body), &person); err != nil {
		return response(400, "Invalid JSON")
	}

	created, err := h.service.CreatePerson(ctx, person)
	if err != nil {
		return errorResponse(err)
	}

	bodyBytes, _ := json.Marshal(created)
	return response(200, string(bodyBytes))
}

func (h *PersonHandler) handleGet(ctx context.Context) (map[string]interface{}, error) {

	persons, err := h.service.ListPersons(ctx)
	if err != nil {
		return errorResponse(err)
	}

	body, _ := json.Marshal(persons)
	return response(200, string(body))
}

func getMethod(event map[string]interface{}) string {
	// REST API format
	if m, ok := event["httpMethod"].(string); ok {
		return m
	}

	// HTTP API (v2) format
	if rc, ok := event["requestContext"].(map[string]interface{}); ok {
		if http, ok := rc["http"].(map[string]interface{}); ok {
			if m, ok := http["method"].(string); ok {
				return m
			}
		}
	}

	return ""
}

func response(status int, body string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"statusCode": status,
		"body":       body,
	}, nil
}

func errorResponse(err error) (map[string]interface{}, error) {
	return map[string]interface{}{
		"statusCode": 500,
		"body":       err.Error(),
	}, nil
}
