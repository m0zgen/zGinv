package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"zGinv/db"
)

var sshFile string

var importSshCmd = &cobra.Command{
	Use:   "import-ssh",
	Short: "Импортирует хосты из SSH-конфига",
	Run: func(cmd *cobra.Command, args []string) {
		if sshFile == "" {
			fmt.Println("Укажи путь к SSH-конфигу через --file")
			return
		}
		if err := importSSHConfig(sshFile); err != nil {
			fmt.Println("Ошибка импорта:", err)
		}
	},
}

func init() {
	importSshCmd.Flags().StringVar(&sshFile, "file", "", "Путь к SSH config (например, ~/.ssh/bld.conf)")
	rootCmd.AddCommand(importSshCmd)
}

func importSSHConfig(filePath string) error {
	path := os.ExpandEnv(filePath)
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer f.Close()

	// Извлекаем имя группы из имени файла
	groupName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	var g db.Group
	if err := db.DB.FirstOrCreate(&g, db.Group{Name: groupName}).Error; err != nil {
		return fmt.Errorf("не удалось создать/получить группу %s: %w", groupName, err)
	}
	groupID := g.ID

	scanner := bufio.NewScanner(f)
	var (
		current  ServerStub
		comment  string
		imported int
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			comment = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key, value := parts[0], strings.Join(parts[1:], " ")

		switch strings.ToLower(key) {
		case "host":
			if current.Name != "" {
				saveServer(current)
				imported++
			}
			current = ServerStub{
				Name:    value,
				Port:    22,
				User:    "root",
				Comment: comment,
				GroupID: groupID,
			}
			comment = ""
		case "hostname":
			current.Address = value
		case "port":
			if p, err := strconv.Atoi(value); err == nil {
				current.Port = p
			}
		case "user":
			current.User = value
		}
	}
	// сохраняем последний
	if current.Name != "" {
		saveServer(current)
		imported++
	}
	fmt.Printf("Импорт завершён. Добавлено %d серверов в группу '%s'.\n", imported, groupName)
	return nil
}

type ServerStub struct {
	Name    string
	Address string
	Port    int
	User    string
	Comment string
	GroupID uint
}

func saveServer(s ServerStub) {
	server := db.Server{
		Name:    s.Name,
		Address: s.Address,
		Port:    s.Port,
		User:    s.User,
		Comment: s.Comment,
		GroupID: s.GroupID,
		Tags:    "imported",
	}
	if err := db.DB.Create(&server).Error; err != nil {
		fmt.Printf("Ошибка при добавлении %s: %v\n", s.Name, err)
	}
}
