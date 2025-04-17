// cmd/find.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"zGinv/db"
)

var (
	findName  string
	findGroup string
	findTag   string
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Поиск хостов по имени или тегу",
	Run: func(cmd *cobra.Command, args []string) {
		dbQuery := db.DB.Model(&db.Server{})

		if findName != "" {
			//dbQuery = dbQuery.Where("name LIKE ?", findName)
			pattern := strings.ReplaceAll(findName, "*", "%")
			dbQuery = dbQuery.Where("name LIKE ?", pattern)
		}
		if findGroup != "" {
			var g db.Group
			if err := db.DB.Where("name = ?", findGroup).First(&g).Error; err == nil {
				dbQuery = dbQuery.Where("group_id = ?", g.ID)
			}
		}
		if findTag != "" {
			dbQuery = dbQuery.Where("tags LIKE ?", "%"+findTag+"%")
		}

		var results []db.Server
		dbQuery.Find(&results)

		if len(results) == 0 {
			fmt.Println("Ничего не найдено")
			return
		}

		for _, s := range results {
			fmt.Printf("%-15s %-16s %-6d %-10s [tags: %s]\n",
				s.Name, s.Address, s.Port, s.Region, strings.ReplaceAll(s.Tags, ",", " "))
		}
	},
}

func init() {
	findCmd.Flags().StringVar(&findName, "name", "", "Поиск по имени (с поддержкой * шаблона)")
	findCmd.Flags().StringVar(&findGroup, "group", "", "Искать только внутри группы")
	findCmd.Flags().StringVar(&findTag, "tag", "", "Поиск по тегу (частичное совпадение)")
	rootCmd.AddCommand(findCmd)
}
