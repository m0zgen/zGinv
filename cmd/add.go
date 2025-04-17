package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"zGinv/db"
)

var (
	name    string
	address string
	port    int
	user    string
	project string
	region  string
	tags    string
	group   string
	comment string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Добавляет сервер в базу",
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" || address == "" {
			fmt.Println("Поля --name и --address обязательны.")
			return
		}

		var groupID uint
		if group != "" {
			var g db.Group
			if err := db.DB.FirstOrCreate(&g, db.Group{Name: group}).Error; err != nil {
				fmt.Println("Ошибка при создании/поиске группы:", err)
				return
			}
			groupID = g.ID
		}

		server := db.Server{
			Name:    name,
			Address: address,
			Port:    port,
			User:    user,
			Project: project,
			Region:  region,
			Tags:    tags,
			Comment: comment,
			GroupID: groupID,
		}

		if err := db.DB.Create(&server).Error; err != nil {
			fmt.Println("Ошибка при добавлении сервера:", err)
			return
		}

		fmt.Println("Сервер добавлен:", name)
	},
}

func init() {
	addCmd.Flags().StringVar(&name, "name", "", "Уникальное имя сервера (обязательное)")
	addCmd.Flags().StringVar(&address, "address", "", "IP-адрес сервера (обязательное)")
	addCmd.Flags().IntVar(&port, "port", 22, "Порт SSH (по умолчанию 22)")
	addCmd.Flags().StringVar(&user, "user", "root", "Пользователь SSH")
	addCmd.Flags().StringVar(&project, "project", "", "Проект")
	addCmd.Flags().StringVar(&region, "region", "", "Регион")
	addCmd.Flags().StringVar(&tags, "tags", "", "Теги (через запятую)")
	addCmd.Flags().StringVar(&group, "group", "", "Группа (будет создана при необходимости)")
	addCmd.Flags().StringVar(&comment, "comment", "", "Комментарий")

	rootCmd.AddCommand(addCmd)
}
