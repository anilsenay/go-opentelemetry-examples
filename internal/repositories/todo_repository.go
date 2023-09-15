package repositories

import (
	"context"

	"github.com/anilsenay/go-opentelemetry-example/internal/models"
	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

func (r *TodoRepository) GetTodos(ctx context.Context) ([]models.Todo, error) {
	var todos []models.Todo
	result := r.db.WithContext(ctx).Table("todos").Select("*").Find(&todos)

	if result.Error != nil {
		return nil, result.Error
	}

	return todos, nil
}

func (r *TodoRepository) GetTodoById(ctx context.Context, id int) (models.Todo, error) {
	var todo models.Todo
	result := r.db.WithContext(ctx).Table("todos").Select("*").Where("id = ?", id).First(&todo)

	if result.Error != nil {
		return models.Todo{}, result.Error
	}

	return todo, nil
}
