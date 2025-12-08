// Package model defines RBAC (Role-Based Access Control) structures
package model

// Role represents a user role in the system
type Role string

const (
	RoleAdmin  Role = "admin"  // Full access to all features
	RoleVendor Role = "vendor" // Limited access - can only manage assigned inbounds
)

// UserRole links users to their roles
type UserRole struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId int    `json:"userId" gorm:"index;not null"`              // Foreign key to User
	Role   Role   `json:"role" gorm:"type:varchar(20);default:vendor"` // User's role
	User   User   `json:"user" gorm:"foreignKey:UserId"`               // Relation to User
}

// InboundAccess maps which inbounds a vendor can access
// Admin users don't need entries here - they have access to everything
type InboundAccess struct {
	Id        int `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId    int `json:"userId" gorm:"index;not null"`    // Foreign key to User
	InboundId int `json:"inboundId" gorm:"index;not null"` // Foreign key to Inbound
}

// TableName overrides the default table name
func (UserRole) TableName() string {
	return "user_roles"
}

// TableName overrides the default table name
func (InboundAccess) TableName() string {
	return "inbound_access"
}
