package controller

import (
	"context"
	"database/sql"
	"errors"

	"api_gateway/internal/entity"
	"api_gateway/protos/protos/crud_pb"

	"github.com/gofiber/fiber/v2"
)

type CRUDController struct {
	service crud_pb.PostCRUDServiceClient
}

func NewCRUDController(service crud_pb.PostCRUDServiceClient) *CRUDController {
	return &CRUDController{service: service}
}

func (ct *CRUDController) List(c *fiber.Ctx) error {
	var query struct {
		Limit  int `query:"limit"`
		Offset int `query:"offset"`
	}
	if err := c.QueryParser(&query); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	posts, err := ct.service.List(context.Background(), &crud_pb.ListRequest{
		Limit:  int32(query.Limit),
		Offset: int32(query.Offset),
	})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).JSON(posts)
}

func (ct *CRUDController) Detail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(404)
	}

	post, err := ct.service.Detail(context.Background(), &crud_pb.DetailRequest{Id: int64(id)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.SendStatus(404)
		}
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).JSON(post)
}

func (ct *CRUDController) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(404)
	}

	var post entity.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	_, err = ct.service.Update(context.Background(), &crud_pb.UpdateRequest{
		Post: &crud_pb.Post{
			Id:     int64(id),
			UserId: int64(post.UserId),
			Title:  post.Title,
			Body:   post.Body,
		}})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(200)
}

func (ct *CRUDController) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.SendStatus(404)
	}

	_, err = ct.service.Delete(context.Background(), &crud_pb.DeleteRequest{Id: int64(id)})
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendStatus(204)
}
