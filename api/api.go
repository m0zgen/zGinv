// api/api.go
package api

import (
	"github.com/gofiber/fiber/v3"
	"zGinv/db"
)

func RegisterRoutes(r fiber.Router) {
	r.Get("/servers", listServers)
	r.Post("/servers", addServer)
	r.Put("/servers/:name", editServer)
	r.Delete("/servers/:name", deleteServer)
	r.Get("/groups", listGroups)
	r.Get("/servers/export", exportServers)
}

// listServers godoc
// @Summary Получить список серверов
// @Tags Servers
// @Produce json
// @Param project query string false "Фильтр по проекту"
// @Param region query string false "Фильтр по региону"
// @Param tag query string false "Фильтр по тегу"
// @Param group query string false "Фильтр по группе"
// @Success 200 {array} db.Server
// @Router /servers [get]
func listServers(c fiber.Ctx) error {
	var servers []db.Server
	query := db.DB

	if project := c.Query("project"); project != "" {
		query = query.Where("project = ?", project)
	}
	if region := c.Query("region"); region != "" {
		query = query.Where("region = ?", region)
	}
	if tag := c.Query("tag"); tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}
	if group := c.Query("group"); group != "" {
		var g db.Group
		if err := db.DB.Where("name = ?", group).First(&g).Error; err == nil {
			query = query.Where("group_id = ?", g.ID)
		}
	}

	query.Find(&servers)
	return c.JSON(servers)
}

// addServer godoc
// @Summary Добавить сервер
// @Tags Servers
// @Accept json
// @Produce json
// @Param server body db.Server true "Данные сервера"
// @Success 201 {object} db.Server
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Ошибка добавления"
// @Router /servers [post]
func addServer(c fiber.Ctx) error {
	var s db.Server
	if err := c.Bind().Body(&s); err != nil {
		return c.Status(400).SendString("Неверный формат запроса")
	}
	if s.Name == "" || s.Address == "" {
		return c.Status(400).SendString("Имя и адрес обязательны")
	}
	if err := db.DB.Create(&s).Error; err != nil {
		return c.Status(500).SendString("Ошибка добавления")
	}
	return c.Status(201).JSON(s)
}

// editServer godoc
// @Summary Обновить данные сервера
// @Tags Servers
// @Accept json
// @Produce json
// @Param name path string true "Имя сервера"
// @Param server body db.Server true "Обновлённые данные"
// @Success 200 {object} db.Server
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 404 {string} string "Сервер не найден"
// @Failure 500 {string} string "Ошибка сохранения"
// @Router /servers/{name} [put]
func editServer(c fiber.Ctx) error {
	name := c.Params("name")
	var s db.Server
	if err := db.DB.Where("name = ?", name).First(&s).Error; err != nil {
		return c.Status(404).SendString("Сервер не найден")
	}
	if err := c.Bind().Body(&s); err != nil {
		return c.Status(400).SendString("Неверный формат запроса")
	}
	if err := db.DB.Save(&s).Error; err != nil {
		return c.Status(500).SendString("Ошибка сохранения")
	}
	return c.JSON(s)
}

// deleteServer godoc
// @Summary Удалить сервер
// @Tags Servers
// @Produce plain
// @Param name path string true "Имя сервера"
// @Success 204 {string} string "Сервер удалён"
// @Failure 500 {string} string "Ошибка удаления"
// @Router /servers/{name} [delete]
func deleteServer(c fiber.Ctx) error {
	name := c.Params("name")
	if err := db.DB.Where("name = ?", name).Delete(&db.Server{}).Error; err != nil {
		return c.Status(500).SendString("Ошибка удаления")
	}
	return c.SendStatus(204)
}

// listGroups godoc
// @Summary Получить список групп
// @Tags Groups
// @Produce json
// @Success 200 {array} db.Group
// @Router /groups [get]
func listGroups(c fiber.Ctx) error {
	var groups []db.Group
	db.DB.Preload("Servers").Find(&groups)
	return c.JSON(groups)
}

// exportServers godoc
// @Summary Экспортировать список серверов
// @Tags Servers
// @Produce json
// @Success 200 {array} db.Server
// @Router /servers/export [get]
func exportServers(c fiber.Ctx) error {
	var servers []db.Server
	db.DB.Find(&servers)
	return c.JSON(servers)
}
