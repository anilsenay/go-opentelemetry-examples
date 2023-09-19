package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/anilsenay/go-opentelemetry-example-2/internal/models"
	"github.com/anilsenay/go-opentelemetry-example-2/pkg/tracer"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TodoService struct {
}

func NewTodoService() *TodoService {
	return &TodoService{}
}

func (s *TodoService) GetTodos(ctx context.Context) ([]models.Todo, error) {
	defer tracer.Trace(&ctx, "TodoService-GetTodos")()

	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3003/todos", nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do request error: %w", err)
	}

	defer resp.Body.Close()

	var target []models.Todo
	err = json.NewDecoder(resp.Body).Decode(&target)

	return target, err
}

func (s *TodoService) GetTodoById(ctx context.Context, id int) (models.Todo, error) {
	defer tracer.Trace(&ctx, "TodoService-GetTodoById", trace.WithAttributes(attribute.Int("id", id)))()

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://localhost:3003/todo/%d", id), nil)
	if err != nil {
		return models.Todo{}, fmt.Errorf("create request error: %w", err)
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   10 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		return models.Todo{}, fmt.Errorf("do request error: %w", err)
	}

	defer resp.Body.Close()

	var target models.Todo
	err = json.NewDecoder(resp.Body).Decode(&target)

	return target, err
}
