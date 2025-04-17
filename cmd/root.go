// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zGinv",
	Short: "zGinv — tool for centralized inventory and management of VPS servers",
	Long:  "CLI-tool for managing VPS servers with a focus on simplicity and efficiency.",
}

// Execute запускает корневую команду
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Ошибка:", err)
		os.Exit(1)
	}
}
