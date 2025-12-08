package service

import (
	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"gorm.io/gorm"
)

type RBACService struct{}

// GetUserRole returns the role of a user by userId
func (s *RBACService) GetUserRole(userId int) (model.Role, error) {
	db := database.GetDB()
	var userRole model.UserRole

	err := db.Where("user_id = ?", userId).First(&userRole).Error
	if err != nil {
		if database.IsNotFound(err) {
			// Default to vendor role if not found
			return model.RoleVendor, nil
		}
		return "", err
	}

	return userRole.Role, nil
}

// IsAdmin checks if a user has admin role
func (s *RBACService) IsAdmin(userId int) (bool, error) {
	role, err := s.GetUserRole(userId)
	if err != nil {
		return false, err
	}
	return role == model.RoleAdmin, nil
}

// IsVendor checks if a user has vendor role
func (s *RBACService) IsVendor(userId int) (bool, error) {
	role, err := s.GetUserRole(userId)
	if err != nil {
		return false, err
	}
	return role == model.RoleVendor, nil
}

// AssignRole assigns a role to a user
func (s *RBACService) AssignRole(userId int, role model.Role) error {
	db := database.GetDB()

	// Check if user already has a role
	var existingRole model.UserRole
	err := db.Where("user_id = ?", userId).First(&existingRole).Error

	if err != nil {
		if database.IsNotFound(err) {
			// Create new role
			userRole := &model.UserRole{
				UserId: userId,
				Role:   role,
			}
			return db.Create(userRole).Error
		}
		return err
	}

	// Update existing role
	return db.Model(&existingRole).Update("role", role).Error
}

// GrantInboundAccess grants a vendor access to an inbound
func (s *RBACService) GrantInboundAccess(userId int, inboundId int) error {
	db := database.GetDB()

	// Check if access already exists
	var existing model.InboundAccess
	err := db.Where("user_id = ? AND inbound_id = ?", userId, inboundId).First(&existing).Error

	if database.IsNotFound(err) {
		// Create new access
		access := &model.InboundAccess{
			UserId:    userId,
			InboundId: inboundId,
		}
		return db.Create(access).Error
	}

	return err // Already exists or other error
}

// RevokeInboundAccess removes vendor's access to an inbound
func (s *RBACService) RevokeInboundAccess(userId int, inboundId int) error {
	db := database.GetDB()
	return db.Where("user_id = ? AND inbound_id = ?", userId, inboundId).Delete(&model.InboundAccess{}).Error
}

// GetVendorInbounds returns all inbound IDs a vendor has access to
func (s *RBACService) GetVendorInbounds(userId int) ([]int, error) {
	db := database.GetDB()
	var accesses []model.InboundAccess

	err := db.Where("user_id = ?", userId).Find(&accesses).Error
	if err != nil {
		return nil, err
	}

	inboundIds := make([]int, len(accesses))
	for i, access := range accesses {
		inboundIds[i] = access.InboundId
	}

	return inboundIds, nil
}

// CanAccessInbound checks if a user can access a specific inbound
// Admins have access to all inbounds
// Vendors only have access to assigned inbounds
func (s *RBACService) CanAccessInbound(userId int, inboundId int) (bool, error) {
	// Check if user is admin
	isAdmin, err := s.IsAdmin(userId)
	if err != nil {
		return false, err
	}

	if isAdmin {
		return true, nil // Admins have access to everything
	}

	// Check vendor access
	db := database.GetDB()
	var access model.InboundAccess
	err = db.Where("user_id = ? AND inbound_id = ?", userId, inboundId).First(&access).Error

	if database.IsNotFound(err) {
		return false, nil // No access
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// FilterInboundsByAccess filters inbounds based on user's access
// Returns all inbounds for admin, filtered for vendors
func (s *RBACService) FilterInboundsByAccess(userId int, inbounds []*model.Inbound) ([]*model.Inbound, error) {
	isAdmin, err := s.IsAdmin(userId)
	if err != nil {
		return nil, err
	}

	if isAdmin {
		return inbounds, nil // Return all for admin
	}

	// Get vendor's accessible inbound IDs
	inboundIds, err := s.GetVendorInbounds(userId)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	accessibleMap := make(map[int]bool)
	for _, id := range inboundIds {
		accessibleMap[id] = true
	}

	// Filter inbounds
	filtered := make([]*model.Inbound, 0)
	for _, inbound := range inbounds {
		if accessibleMap[inbound.Id] {
			filtered = append(filtered, inbound)
		}
	}

	return filtered, nil
}

// GetInboundsWithAccess returns inbounds with access control applied
func (s *RBACService) GetInboundsWithAccess(userId int) ([]*model.Inbound, error) {
	db := database.GetDB()

	isAdmin, err := s.IsAdmin(userId)
	if err != nil {
		return nil, err
	}

	var inbounds []*model.Inbound

	if isAdmin {
		// Admin gets all inbounds
		err = db.Model(&model.Inbound{}).Find(&inbounds).Error
	} else {
		// Vendor gets only assigned inbounds
		inboundIds, err := s.GetVendorInbounds(userId)
		if err != nil {
			return nil, err
		}

		if len(inboundIds) == 0 {
			return []*model.Inbound{}, nil // No access to any inbound
		}

		err = db.Model(&model.Inbound{}).Where("id IN ?", inboundIds).Find(&inbounds).Error
	}

	return inbounds, err
}

// GetAllVendors returns all users with vendor role
func (s *RBACService) GetAllVendors() ([]model.User, error) {
	db := database.GetDB()
	var userRoles []model.UserRole

	err := db.Where("role = ?", model.RoleVendor).Preload("User").Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	users := make([]model.User, len(userRoles))
	for i, ur := range userRoles {
		users[i] = ur.User
	}

	return users, nil
}

// CreateVendorWithInbounds creates a new vendor user and assigns inbound access
func (s *RBACService) CreateVendorWithInbounds(username, password string, inboundIds []int) (*model.User, error) {
	db := database.GetDB()

	// Start transaction
	return nil, db.Transaction(func(tx *gorm.DB) error {
		// Create user
		user := &model.User{
			Username: username,
			Password: password, // Should be hashed before calling this
		}

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		// Assign vendor role
		userRole := &model.UserRole{
			UserId: user.Id,
			Role:   model.RoleVendor,
		}

		if err := tx.Create(userRole).Error; err != nil {
			return err
		}

		// Grant inbound access
		for _, inboundId := range inboundIds {
			access := &model.InboundAccess{
				UserId:    user.Id,
				InboundId: inboundId,
			}
			if err := tx.Create(access).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
