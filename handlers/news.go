package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"news/db"
	"news/models"
	"news/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var Id string

func EditNews(c *fiber.Ctx) error {
	Id = c.Params("Id")
	var news models.News

	// Парсинг тела запроса
	if err := c.BodyParser(&news); err != nil {
		utils.Log.Error("Failed to parse body: ", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Выполнение обновления
	query := fmt.Sprintf("UPDATE News SET Title = COALESCE(NULLIF('%s', ''), Title), Content = COALESCE(NULLIF('%s', ''), Content) WHERE Id = %s::numeric;", news.Title, news.Content, Id[1:])
	result, err := db.ReformDB.Exec(query)
	if err != nil {
		utils.Log.Error("Failed to update news: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update news"})
	}

	// Проверка количества затронутых строк
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.Log.Error("Failed to get rows affected: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update news"})
	}

	if rowsAffected == 0 {
		utils.Log.Warn("No news found with the given ID: ", Id)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No news found with the given ID"})
	}

	category_del_query := fmt.Sprintf("DELETE FROM newscategories WHERE newsid = %s::numeric", Id[1:])
	db.ReformDB.Exec(category_del_query)

	// Вставка категорий и проверка затронутых строк
	var category_result sql.Result
	for _, category := range news.Categories {
		utils.Log.Info("category - ", category)
		category_query := fmt.Sprintf("INSERT INTO newscategories (newsid, categoryid) VALUES (%s::numeric, %d) ON CONFLICT (newsid, categoryid) DO NOTHING", Id[1:], category)
		category_result, err = db.ReformDB.Exec(category_query)
		if err != nil {
			utils.Log.Error("Failed to update categories: ", err, category)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update categories"})
		}
		categoryRowsAffected, err := category_result.RowsAffected()
		if err != nil {
			utils.Log.Error("Failed to get rows affected: ", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update categories"})
		}

		if categoryRowsAffected == 0 {
			utils.Log.Warn("Failed to update categories")
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Failed to update categories"})
		}
	}

	utils.Log.Info("News updated successfully: ", Id)
	return c.JSON(fiber.Map{"message": "News updated successfully"})
}

// Функция отображения новостей с пагинацией
func ListNews(c *fiber.Ctx) error {
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		utils.Log.Warn("Invalid page number: ", pageStr)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid page number"})
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		utils.Log.Warn("Invalid limit number: ", limitStr)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid limit number"})
	}

	offset := (page - 1) * limit

	var newsList []models.News

	query := fmt.Sprintf("SELECT n.id, n.title, n.content, string_agg(nc.categoryid::text, ',') "+
		"FROM News n "+
		"JOIN newscategories nc "+
		"ON n.id = nc.newsid "+
		"GROUP BY n.id "+
		"ORDER BY Id "+
		"LIMIT %d "+
		"OFFSET %d", limit, offset)
	rows, err := db.ReformDB.Query(query)
	if err != nil {
		utils.Log.Error("Failed to retrieve news: ", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve news"})
	}
	defer rows.Close()

	for rows.Next() {
		var news models.News
		var categories string
		if err := rows.Scan(&news.Id, &news.Title, &news.Content, &categories); err != nil {
			utils.Log.Error("Failed to scan news: ", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve news"})
		}

		if categories != "" {
			categoryStrings := strings.Split(categories, ",")
			for _, categoryStr := range categoryStrings {
				categoryID, err := strconv.Atoi(categoryStr)
				if err == nil {
					news.Categories = append(news.Categories, categoryID)
				}
			}
		}

		newsList = append(newsList, news)
	}

	return c.JSON(fiber.Map{"News": newsList})
}
