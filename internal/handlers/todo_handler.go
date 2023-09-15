package handlers

import (
	"strconv"

	"github.com/anilsenay/go-opentelemetry-example/internal/services"
	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	todoService *services.TodoService
}

func NewTodoHandler(s *services.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: s,
	}
}

func (h *TodoHandler) handleGetTodos(c *fiber.Ctx) error {
	todos, err := h.todoService.GetTodos(c.UserContext())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(todos)
}

func (h *TodoHandler) handleGetTodoById(c *fiber.Ctx) error {
	todoId := c.Params("id")
	id, err := strconv.Atoi(todoId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	todo, err := h.todoService.GetTodoById(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.JSON(todo)
}

func (h *TodoHandler) SetRoutes(app *fiber.App) {
	app.Get("/todos", h.handleGetTodos)
	app.Get("/todos/:id", h.handleGetTodoById)
}
