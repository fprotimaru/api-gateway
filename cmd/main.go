package main

import (
	"log"

	"api_gateway/internal/controller"
	"api_gateway/protos/protos/crud_pb"
	"api_gateway/protos/protos/parser_pb"

	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: jsoniter.Marshal,
		JSONDecoder: jsoniter.Unmarshal,
	})

	parserServer, err := grpc.Dial("localhost:8001", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	parserService := parser_pb.NewPostParserServiceClient(parserServer)
	parserController := controller.NewParserController(parserService)

	app.Get("/start-parse", parserController.StartParsing)

	crudServer, err := grpc.Dial("localhost:8002", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	crudService := crud_pb.NewPostCRUDServiceClient(crudServer)
	crudController := controller.NewCRUDController(crudService)

	app.Get("/post/list", crudController.List)
	app.Get("/post/detail/:id", crudController.Detail)
	app.Get("/post/update/:id", crudController.Update)
	app.Get("/post/delete/:id", crudController.Delete)

	if err := app.Listen(":8000"); err != nil {
		log.Fatalln(err)
	}
}
