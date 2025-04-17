// Package db db/server.go
package db

import (
	"time"
)

type Group struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"core-dns" description:"Server group name"`
	Comment   string    `json:"comment" example:"CoreDNS servers" description:"Description of the group"`
	CreatedAt time.Time `json:"created_at" swaggertype:"string" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" swaggertype:"string" format:"date-time"`
	Servers   []*Server `json:"servers,omitempty" gorm:"many2many:group_servers;"`
}

type Server struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"dns-kz-1" description:"Unique server name"`
	Address   string    `json:"address" example:"185.100.100.1" description:"Server IP address"`
	Port      int       `json:"port" example:"22" description:"SSH or service port"`
	User      string    `json:"user" example:"root" description:"Username for access"`
	Project   string    `json:"project" example:"openbld" description:"Project name"`
	Region    string    `json:"region" example:"kz" description:"Deployment region"`
	Tags      string    `json:"tags" example:"dns,edge" description:"Server tags"`
	Comment   string    `json:"comment" example:"Backup DNS" description:"Additional server comment"`
	GroupID   uint      `json:"group_id"`
	Groups    []*Group  `json:"-" gorm:"many2many:group_servers;"` // hidden from Swagger
	CreatedAt time.Time `json:"created_at" swaggertype:"string" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" swaggertype:"string" format:"date-time"`
}
