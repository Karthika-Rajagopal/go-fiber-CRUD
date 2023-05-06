package main

import (
	"strconv"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
)

type Todo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos = []*Todo{
	{Id: 1, Name: "Walk the dog", Completed: false},
	{Id: 2, Name: "Walk the cat", Completed: false},
}

func main() {
	app := fiber.New()  //sets up the new instance of the fiber app

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("hello world")
	})

	SetupApiV1(app)  //to setup api endpoints

	err := app.Listen(3000)     //app is started on port 3000 using the Listen function. If there is an error during startup, the program panics
	if err != nil {
		panic(err)
	}
}

func SetupApiV1(app *fiber.App) {     //creates a new group for version 1 of the API and then calls the SetupTodosRoutes function to set up the todo-specific endpoints
	v1 := app.Group("/v1")

	SetupTodosRoutes(v1)
}

func SetupTodosRoutes(grp fiber.Router) {    //SetupTodosRoutes function sets up the following endpoints
	todosRoutes := grp.Group("/todos")
	todosRoutes.Get("/", GetTodos)
	todosRoutes.Post("/", CreateTodo)
	todosRoutes.Get("/:id", GetTodo)
	todosRoutes.Delete("/:id", DeleteTodo)
	todosRoutes.Patch("/:id", UpdateTodo)
}

func UpdateTodo(ctx *fiber.Ctx) {       //UpdateTodo function reads the ID of a specific todo item from the URL parameters, finds the corresponding todo item in the list, updates its properties based on the request body, and returns the updated item in the response body
	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	var body request         //the request body is defined using a struct with optional fields for the name and completed status of the todo item. The body is parsed using the ctx.BodyParser function, which attempts to deserialize the JSON request body into the request struct. If there is an error during parsing, an error response is returned. If there is no error, the corresponding todo item is updated and returned in the response body. If the requested ID is not found in the list of todos, a 404 Not Found response is returned
	err = ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
		return
	}

	var todo *Todo

	for _, t := range todos {
		if t.Id == id {
			todo = t
			break
		}
	}

	if todo == nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = *body.Completed
	}

	ctx.Status(fiber.StatusOK).JSON(todo)
}

func DeleteTodo(ctx *fiber.Ctx) {       //DeleteTodo function reads the ID of a specific todo item from the URL parameters, finds the corresponding todo item in the list, removes it from the list, and returns a 204 No Content response
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	for i, todo := range todos {
		if todo.Id == id {
			todos = append(todos[0:i], todos[i+1:]...)
			ctx.Status(fiber.StatusNoContent)
			return
		}
	}

	ctx.Status(fiber.StatusNotFound)
}

func GetTodo(ctx *fiber.Ctx) {               //GetTodo function reads the ID of a specific todo item from the URL parameters, finds the corresponding todo item in the list, and returns it in the response body
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	for _, todo := range todos {
		if todo.Id == id {
			ctx.Status(fiber.StatusOK).JSON(todo)
			return
		}
	}

	ctx.Status(fiber.StatusNotFound)
}

func CreateTodo(ctx *fiber.Ctx) {   //CreateTodo function reads a new todo item from the request body, adds it to the list of todos, and returns the new item in the response body
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return
	}

	todo := &Todo{
		Id:        len(todos) + 1,
		Name:      body.Name,
		Completed: false,
	}

	todos = append(todos, todo)

	ctx.Status(fiber.StatusCreated).JSON(todo)
}

func GetTodos(ctx *fiber.Ctx) {           //GetTodos function returns the entire list of todo items in the response body
	ctx.Status(fiber.StatusOK).JSON(todos)
}
