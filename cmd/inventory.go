package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"zGinv/db"
)

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Dynamic inventory для Ansible (формат JSON)",
	Run: func(cmd *cobra.Command, args []string) {
		var servers []db.Server
		if err := db.DB.Find(&servers).Error; err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка при получении серверов:", err)
			os.Exit(1)
		}

		result := make(map[string]interface{})
		all := make(map[string]interface{})
		hosts := make([]string, 0)
		hostvars := make(map[string]map[string]interface{})

		for _, s := range servers {
			hosts = append(hosts, s.Name)
			hostvars[s.Name] = map[string]interface{}{
				"ansible_host": s.Address,
				"ansible_port": s.Port,
				"ansible_user": s.User,
			}
		}

		all["hosts"] = hosts
		all["vars"] = map[string]interface{}{}
		result["all"] = all
		result["_meta"] = map[string]interface{}{
			"hostvars": hostvars,
		}

		json.NewEncoder(os.Stdout).Encode(result)
	},
}

func init() {
	rootCmd.AddCommand(inventoryCmd)
}
