package services

import (
	"context"

	"github.com/anilsenay/go-opentelemetry-example/models"
	"github.com/anilsenay/go-opentelemetry-example/repositories"
	"github.com/anilsenay/go-opentelemetry-example/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TodoService struct {
	todoRepository *repositories.TodoRepository
}

func NewTodoService(r *repositories.TodoRepository) *TodoService {
	return &TodoService{
		todoRepository: r,
	}
}

func (s *TodoService) GetTodos(ctx context.Context) ([]models.Todo, error) {
	defer tracer.Trace(&ctx, "TodoService-GetTodos")()

	return s.todoRepository.GetTodos(ctx)
}

func (s *TodoService) GetTodoById(ctx context.Context, id int) (models.Todo, error) {
	defer tracer.Trace(&ctx, "TodoService-GetTodoById", trace.WithAttributes(attribute.Int("id", id)))()

	return s.todoRepository.GetTodoById(ctx, id)
}
