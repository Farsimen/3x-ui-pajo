package controller

import (
	"strconv"
	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/logger"
	"github.com/mhsanaei/3x-ui/v2/web/service"
	"github.com/mhsanaei/3x-ui/v2/web/session"
	"github.com/mhsanaei/3x-ui/v2/util/crypto"

	"github.com/gin-gonic/gin"
)

type VendorController struct {
	BaseController
	
	rbacService service.RBACService
}

func NewVendorController(g *gin.RouterGroup) *VendorController {
	a := &VendorController{}
	a.initRouter(g)
	return a
}

func (a *VendorController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/vendor")
	g.Use(a.checkLogin)

	// Vendor management endpoints (admin only)
	g.POST("/create", a.createVendor)
	g.GET("/list", a.listVendors)
	g.POST("/grant", a.grantAccess)
	g.POST("/revoke", a.revokeAccess)
	g.DELETE("/delete/:id", a.deleteVendor)
	g.GET("/inbounds/:id", a.getVendorInbounds)
}

// createVendor creates a new vendor account and assigns inbounds
func (a *VendorController) createVendor(c *gin.Context) {
	// Check if user is admin
	if !session.IsLogin(c) {
		jsonMsg(c, "Unauthorized", nil)
		return
	}

	user := session.GetLoginUser(c)
	if user == nil {
		jsonMsg(c, "Invalid session", nil)
		return
	}

	isAdmin, err := a.rbacService.IsAdmin(user.Id)
	if err != nil || !isAdmin {
		jsonMsg(c, "Admin access required", nil)
		return
	}

	var req struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		InboundIDs []int  `json:"inboundIds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		jsonMsg(c, "Invalid request: "+err.Error(), nil)
		return
	}

	// Hash password
	hashedPassword, err := crypto.HashPasswordAsBcrypt(req.Password)
	if err != nil {
		jsonMsg(c, "Failed to hash password: "+err.Error(), nil)
		return
	}

	// Create user account directly in database
	db := database.GetDB()
	newUser := &model.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	err = db.Create(newUser).Error
	if err != nil {
		jsonMsg(c, "Failed to create user: "+err.Error(), nil)
		return
	}

	// Assign vendor role
	err = a.rbacService.AssignRole(newUser.Id, "vendor")
	if err != nil {
		// Rollback user creation
		db.Delete(newUser)
		jsonMsg(c, "Failed to assign vendor role: "+err.Error(), nil)
		return
	}

	// Grant access to inbounds
	for _, inboundID := range req.InboundIDs {
		err = a.rbacService.GrantInboundAccess(newUser.Id, inboundID)
		if err != nil {
			logger.Warning("Failed to grant access to inbound", inboundID, ":", err)
		}
	}

	jsonObj(c, map[string]interface{}{
		"success": true,
		"message": "Vendor created successfully",
		"vendorId": newUser.Id,
	}, nil)
}

// listVendors returns all vendor accounts (admin only)
func (a *VendorController) listVendors(c *gin.Context) {
	// Check if user is admin
	if !session.IsLogin(c) {
		jsonMsg(c, "Unauthorized", nil)
		return
	}

	user := session.GetLoginUser(c)
	if user == nil {
		jsonMsg(c, "Invalid session", nil)
		return
	}

	isAdmin, err := a.rbacService.IsAdmin(user.Id)
	if err != nil || !isAdmin {
		jsonMsg(c, "Admin access required", nil)
		return
	}

	// Get all vendors
	var vendors []struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	db := database.GetDB()
	err = db.Table("users").
		Select("users.id, users.username, user_roles.role").
		Joins("INNER JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role = ?", "vendor").
		Scan(&vendors).Error

	if err != nil {
		jsonMsg(c, "Failed to fetch vendors: "+err.Error(), nil)
		return
	}

	// Get inbound access for each vendor
	type VendorDetail struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Role      string `json:"role"`
		Inbounds  []int  `json:"inbounds"`
	}

	var vendorDetails []VendorDetail
	for _, vendor := range vendors {
		inbounds, err := a.rbacService.GetUserInbounds(vendor.ID)
		if err != nil {
			inbounds = []int{}
		}
		vendorDetails = append(vendorDetails, VendorDetail{
			ID:       vendor.ID,
			Username: vendor.Username,
			Role:     vendor.Role,
			Inbounds: inbounds,
		})
	}

	jsonObj(c, vendorDetails, nil)
}

// grantAccess grants a vendor access to an inbound (admin only)
func (a *VendorController) grantAccess(c *gin.Context) {
	// Check if user is admin
	if !session.IsLogin(c) {
		jsonMsg(c, "Unauthorized", nil)
		return
	}

	user := session.GetLoginUser(c)
	if user == nil {
		jsonMsg(c, "Invalid session", nil)
		return
	}

	isAdmin, err := a.rbacService.IsAdmin(user.Id)
	if err != nil || !isAdmin {
		jsonMsg(c, "Admin access required", nil)
		return
	}

	var req struct {
		VendorID  int `json:"vendorId" binding:"required"`
		InboundID int `json:"inboundId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		jsonMsg(c, "Invalid request: "+err.Error(), nil)
		return
	}

	err = a.rbacService.GrantInboundAccess(req.VendorID, req.InboundID)
	if err != nil {
		jsonMsg(c, "Failed to grant access: "+err.Error(), nil)
		return
	}

	jsonObj(c, map[string]interface{}{
		"success": true,
		"message": "Access granted successfully",
	}, nil)
}

// revokeAccess revokes a vendor's access to an inbound (admin only)
func (a *VendorController) revokeAccess(c *gin.Context) {
	// Check if user is admin
	if !session.IsLogin(c) {
		jsonMsg(c, "Unauthorized", nil)
		return
	}

	user := session.GetLoginUser(c)
	if user == nil {
		jsonMsg(c, "Invalid session", nil)
		return
	}

	isAdmin, err := a.rbacService.IsAdmin(user.Id)
	if err != nil || !isAdmin {
		jsonMsg(c, "Admin access required", nil)
		return
	}

	var req struct {
		VendorID  int `json:"vendorId" binding:"required"`
		InboundID int `json:"inboundId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		jsonMsg(c, "Invalid request: "+err.Error(), nil)
		return
	}

	err = a.rbacService.RevokeInboundAccess(req.VendorID, req.InboundID)
	if err != nil {
		jsonMsg(c, "Failed to revoke access: "+err.Error(), nil)
		return
	}

	jsonObj(c, map[string]interface{}{
		"success": true,
		"message": "Access revoked successfully",
	}, nil)
}

// deleteVendor deletes a vendor account (admin only)
func (a *VendorController) deleteVendor(c *gin.Context) {
	// Check if user is admin
	if !session.IsLogin(c) {
		jsonMsg(c, "Unauthorized", nil)
		return
	}

	user := session.GetLoginUser(c)
	if user == nil {
		jsonMsg(c, "Invalid session", nil)
		return
	}

	isAdmin, err := a.rbacService.IsAdmin(user.Id)
	if err != nil || !isAdmin {
		jsonMsg(c, "Admin access required", nil)
		return
	}

	vendorIDStr := c.Param("id")
	vendorID, err := strconv.Atoi(vendorIDStr)
	if err != nil {
		jsonMsg(c, "Invalid vendor ID", nil)
		return
	}

	// Check if user is actually a vendor
	isVendor, err := a.rbacService.IsVendor(vendorID)
	if err != nil || !isVendor {
		jsonMsg(c, "User is not a vendor", nil)
		return
	}

	// Delete user directly from database (will cascade delete roles and accesses)
	db := database.GetDB()
	err = db.Delete(&model.User{}, vendorID).Error
	if err != nil {
		jsonMsg(c, "Failed to delete vendor: "+err.Error(), nil)
		return
	}

	jsonObj(c, map[string]interface{}{
		"success": true,
		"message": "Vendor deleted successfully",
	}, nil)
}

// getVendorInbounds returns all inbounds accessible by a vendor
func (a *VendorController) getVendorInbounds(c *gin.Context) {
	// Check if user is admin
	if !session.IsLogin(c) {
		jsonMsg(c, "Unauthorized", nil)
		return
	}

	user := session.GetLoginUser(c)
	if user == nil {
		jsonMsg(c, "Invalid session", nil)
		return
	}

	isAdmin, err := a.rbacService.IsAdmin(user.Id)
	if err != nil || !isAdmin {
		jsonMsg(c, "Admin access required", nil)
		return
	}

	vendorIDStr := c.Param("id")
	vendorID, err := strconv.Atoi(vendorIDStr)
	if err != nil {
		jsonMsg(c, "Invalid vendor ID", nil)
		return
	}

	inbounds, err := a.rbacService.GetUserInbounds(vendorID)
	if err != nil {
		jsonMsg(c, "Failed to fetch inbounds: "+err.Error(), nil)
		return
	}

	jsonObj(c, inbounds, nil)
}
