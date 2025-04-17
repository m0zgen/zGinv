package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/spf13/cobra"
	"zGinv/db"
)

var (
	exportGroup  string
	exportFormat string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Экспортирует серверы в разных форматах (ansible, csv, json)",
	Run: func(cmd *cobra.Command, args []string) {
		query := db.DB
		if exportGroup != "" {
			var g db.Group
			if err := db.DB.Where("name = ?", exportGroup).First(&g).Error; err != nil {
				fmt.Println("Группа не найдена:", exportGroup)
				return
			}
			query = query.Where("group_id = ?", g.ID)
		}

		var servers []db.Server
		query.Find(&servers)

		if len(servers) == 0 {
			fmt.Println("# Нет подходящих серверов.")
			return
		}

		switch exportFormat {
		case "ansible":
			if exportGroup != "" {
				fmt.Printf("[%s]\n", exportGroup)
			} else {
				fmt.Println("[all]")
			}
			for _, s := range servers {
				fmt.Printf("%s ansible_host=%s ansible_port=%d ansible_user=%s\n",
					s.Name, s.Address, s.Port, s.User)
			}

		case "csv":
			writer := csv.NewWriter(os.Stdout)
			writer.Write([]string{"Name", "Address", "Port", "User", "Project", "Region", "Tags", "GroupID", "Comment"})
			for _, s := range servers {
				writer.Write([]string{
					s.Name,
					s.Address,
					fmt.Sprintf("%d", s.Port),
					s.User,
					s.Project,
					s.Region,
					s.Tags,
					fmt.Sprintf("%d", s.GroupID),
					s.Comment,
				})
			}
			writer.Flush()

		case "json":
			data, err := json.MarshalIndent(servers, "", "  ")
			if err != nil {
				fmt.Println("Ошибка при сериализации:", err)
				return
			}
			fmt.Println(string(data))

		case "yaml":
			groupMap := map[string]map[string]map[string]struct{}{}

			var servers []db.Server
			db.DB.Where("group_id != 0").Find(&servers)

			for _, s := range servers {
				groupName := "" // нужно найти имя по GroupID
				var g db.Group
				if err := db.DB.First(&g, s.GroupID).Error; err == nil {
					groupName = g.Name
				}
				if groupName == "" {
					continue
				}

				if _, ok := groupMap[groupName]; !ok {
					groupMap[groupName] = map[string]map[string]struct{}{
						"hosts": {},
					}
				}
				groupMap[groupName]["hosts"][s.Name] = struct{}{}
			}

			full := map[string]any{
				"all": map[string]any{
					"children": groupMap,
				},
			}

			out, _ := yaml.Marshal(full)
			fmt.Println(string(out))

		default:
			fmt.Println("Поддерживаются форматы: ansible, csv, json")
		}
	},
}

func init() {
	exportCmd.Flags().StringVar(&exportGroup, "group", "", "Фильтр по имени группы")
	exportCmd.Flags().StringVar(&exportFormat, "format", "ansible", "Формат: ansible, csv, json, yaml")
	rootCmd.AddCommand(exportCmd)
}
