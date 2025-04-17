package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"zGinv/db"
)

var (
	projectFilter string
	regionFilter  string
	groupFilter   string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show servers from the database",
	Run: func(cmd *cobra.Command, args []string) {
		query := db.DB.Preload("Group")

		if projectFilter != "" {
			query = query.Where("project = ?", projectFilter)
		}
		if regionFilter != "" {
			query = query.Where("region = ?", regionFilter)
		}
		if groupFilter != "" {
			// Найдём группу по имени
			var group db.Group
			if err := db.DB.Where("name = ?", groupFilter).First(&group).Error; err != nil {
				fmt.Println("Group not found:", groupFilter)
				return
			}
			query = query.Where("group_id = ?", group.ID)
		}

		var servers []db.Server
		query.Find(&servers)

		if len(servers) == 0 {
			fmt.Println("No servers found.")
			return
		}

		for _, s := range servers {
			tags := strings.ReplaceAll(s.Tags, ",", " ")
			fmt.Printf(
				"%-15s %-16s %-6d %-10s [project: %s] [group: %d] [tags: %s]\n",
				s.Name, s.Address, s.Port, s.Region, s.Project, s.GroupID, tags,
			)
		}
	},
}

func init() {
	listCmd.Flags().StringVar(&projectFilter, "project", "", "Filter by project")
	listCmd.Flags().StringVar(&regionFilter, "region", "", "Filter by region")
	listCmd.Flags().StringVar(&groupFilter, "group", "", "Filter by group name")
	rootCmd.AddCommand(listCmd)
}
