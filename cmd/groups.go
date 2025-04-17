package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"zGinv/db"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Показывает список групп и количество серверов в каждой",
	Run: func(cmd *cobra.Command, args []string) {
		var groups []db.Group
		if err := db.DB.Preload("Servers").Order("name ASC").Find(&groups).Error; err != nil {
			fmt.Println("Ошибка при загрузке групп:", err)
			return
		}

		if len(groups) == 0 {
			fmt.Println("Нет групп в базе данных.")
			return
		}

		fmt.Printf("%-20s %-10s %-s\n", "Группа", "Серверов", "Комментарий")
		fmt.Println(strings.Repeat("-", 60))

		for _, g := range groups {
			fmt.Printf("%-20s %-10d %-s\n", g.Name, len(g.Servers), g.Comment)
		}
	},
}

func init() {
	rootCmd.AddCommand(groupsCmd)
}
