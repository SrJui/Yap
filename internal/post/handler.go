package post

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"
)

func CreatePost(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "create_post",
		"request_id", requestid.FromContext(c),
	)

	// get sess info
	sess := session.FromContext(c)

	// get user id from sess info
	userID, ok := sess.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	requestLog = requestLog.With("user_id", userID)

	// get content
	content := c.FormValue("content")
	if !isValidContent(content) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "content is invalid",
		})
	}

	// save post to db
	createdPost, err := Create(c.Context(), content, userID, nil)
	if err != nil {
		requestLog.Error("post creation failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "post creation failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPost)
}

func UpdatePost(c fiber.Ctx) error {

	return c.SendStatus(fiber.StatusOK)
}

func DeletePost(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "delete_post",
		"request_id", requestid.FromContext(c),
	)

	// get sess info
	sess := session.FromContext(c)

	// get user id from sess info
	userID, ok := sess.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	requestLog = requestLog.With("user_id", userID)

	// get id and validate
	postID, err := strconv.ParseUint(c.Params("id"), 10, strconv.IntSize)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Post ID must be a valid positive number",
		})
	}
	if postID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Post ID must be greater than 0",
		})
	}

	// delete post
	err = DeleteSpecificPost(c.Context(), userID, uint(postID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusOK)
	}
	if err != nil {
		requestLog.Error("failed to delete post", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete post",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func GetOnePost(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "get_one_post",
		"request_id", requestid.FromContext(c),
	)

	// get sess info
	sess := session.FromContext(c)

	// get user id from sess info
	userID, ok := sess.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	requestLog = requestLog.With("user_id", userID)

	// get id and validate
	postID, err := strconv.ParseUint(c.Params("id"), 10, strconv.IntSize)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Post ID must be a valid positive number",
		})
	}
	if postID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Post ID must be greater than 0",
		})
	}

	// search for post with id
	foundPost, err := FindByID(c.Context(), uint(postID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		requestLog.Error("error finding post with id", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "error finding the post",
		})
	}

	// return with post info
	return c.JSON(foundPost)
}

func GetAllPosts(c fiber.Ctx) error {
	requestLog := slog.With(
		"action", "get_all_posts",
		"request_id", requestid.FromContext(c),
	)

	// get sess info
	sess := session.FromContext(c)

	// get user id from sess info
	userID, ok := sess.Get("user_id").(uint)
	if !ok || userID == 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	requestLog = requestLog.With("user_id", userID)

	// get limit and validate
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "limit must be an integer between 1 and 100",
		})
	}

	// get offset
	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil || offset < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "offset must be a non-negative integer",
		})
	}

	// get posts
	posts, err := FindAll(c.Context(), limit, offset)
	if err != nil {
		requestLog.Error("post lookup failed", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "post lookup failed",
		})
	}

	return c.JSON(posts)
}
