// main.go
package main

import (
	"zGinv/cmd"
	"zGinv/db"
)

func main() {
	db.InitDB()
	cmd.Execute()
}
