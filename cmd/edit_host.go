// cmd/edit.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"zGinv/db"
)

var editCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Редактирует существующий сервер по имени",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		var server db.Server
		if err := db.DB.Where("name = ?", name).First(&server).Error; err != nil {
			fmt.Printf("Сервер %s не найден\n", name)
			return
		}

		updated := false

		if cmd.Flags().Changed("address") {
			server.Address, _ = cmd.Flags().GetString("address")
			updated = true
		}
		if cmd.Flags().Changed("port") {
			server.Port, _ = cmd.Flags().GetInt("port")
			updated = true
		}
		if cmd.Flags().Changed("user") {
			server.User, _ = cmd.Flags().GetString("user")
			updated = true
		}
		if cmd.Flags().Changed("project") {
			server.Project, _ = cmd.Flags().GetString("project")
			updated = true
		}
		if cmd.Flags().Changed("region") {
			server.Region, _ = cmd.Flags().GetString("region")
			updated = true
		}
		if cmd.Flags().Changed("tags") {
			server.Tags, _ = cmd.Flags().GetString("tags")
			updated = true
		}
		if cmd.Flags().Changed("comment") {
			server.Comment, _ = cmd.Flags().GetString("comment")
			updated = true
		}

		if updated {
			if err := db.DB.Save(&server).Error; err != nil {
				fmt.Println("Ошибка при сохранении:", err)
			} else {
				fmt.Println("Сервер обновлён.")
			}
		} else {
			fmt.Println("Ничего не изменено.")
		}
	},
}

func init() {
	editCmd.Flags().String("address", "", "Новый IP адрес")
	editCmd.Flags().Int("port", 22, "Новый порт")
	editCmd.Flags().String("user", "", "Новый SSH пользователь")
	editCmd.Flags().String("project", "", "Новый проект")
	editCmd.Flags().String("region", "", "Новый регион")
	editCmd.Flags().String("tags", "", "Обновить теги")
	editCmd.Flags().String("comment", "", "Комментарий")
	rootCmd.AddCommand(editCmd)
}
