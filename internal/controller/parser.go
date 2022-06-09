package controller

import (
	"context"

	"api_gateway/protos/protos/parser_pb"

	"github.com/gofiber/fiber/v2"
)

type ParserController struct {
	service parser_pb.PostParserServiceClient
}

func NewParserController(service parser_pb.PostParserServiceClient) *ParserController {
	return &ParserController{service: service}
}

func (ct *ParserController) StartParsing(c *fiber.Ctx) error {
	_, err := ct.service.ParseData(context.Background(), &parser_pb.Empty{})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).SendString("parsing finished")
}
