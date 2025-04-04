// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gorm

import (
	"time"
)

const TableNameBootstrapNode = "bootstrap_nodes"

// BootstrapNode mapped from table <bootstrap_nodes>
type BootstrapNode struct {
	ID                       int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PeerID                   string    `gorm:"column:peer_id;not null" json:"peer_id"`
	MultiAddresses           string    `gorm:"column:multi_addresses;not null" json:"multi_addresses"`
	ProtocolVersion          string    `gorm:"column:protocol_version" json:"protocol_version"`
	Region                   string    `gorm:"column:region" json:"region"`
	IsActive                 bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt                time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	LastUpdated              time.Time `gorm:"column:last_updated;default:CURRENT_TIMESTAMP" json:"last_updated"`
	LastSuccessfulConnection time.Time `gorm:"column:last_successful_connection" json:"last_successful_connection"`
	FailureCount             int32     `gorm:"column:failure_count" json:"failure_count"`
}

// TableName BootstrapNode's table name
func (*BootstrapNode) TableName() string {
	return TableNameBootstrapNode
}
