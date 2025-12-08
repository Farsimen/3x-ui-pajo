package middleware

import (
	"x-ui/web/service"

	"github.com/gin-gonic/gin"
)

type RBACMiddleware struct {
	rbacService service.RBACService
}

func NewRBACMiddleware() *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: service.RBACService{},
	}
}

// RequireAdmin checks if the current user has admin role
func (m *RBACMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, exists := c.Get("session")
		if !exists {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "Unauthorized",
			})
			c.Abort()
			return
		}

		sessionData := session.(gin.H)
		userID, ok := sessionData["id"].(int)
		if !ok {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "Invalid session",
			})
			c.Abort()
			return
		}

		isAdmin, err := m.rbacService.IsAdmin(userID)
		if err != nil || !isAdmin {
			c.JSON(403, gin.H{
				"success": false,
				"msg":     "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireVendorOrAdmin checks if user is vendor or admin
func (m *RBACMiddleware) RequireVendorOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, exists := c.Get("session")
		if !exists {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "Unauthorized",
			})
			c.Abort()
			return
		}

		sessionData := session.(gin.H)
		userID, ok := sessionData["id"].(int)
		if !ok {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "Invalid session",
			})
			c.Abort()
			return
		}

		isAdmin, _ := m.rbacService.IsAdmin(userID)
		isVendor, _ := m.rbacService.IsVendor(userID)

		if !isAdmin && !isVendor {
			c.JSON(403, gin.H{
				"success": false,
				"msg":     "Vendor or Admin access required",
			})
			c.Abort()
			return
		}

		// Store role in context for later use
		if isAdmin {
			c.Set("userRole", "admin")
		} else {
			c.Set("userRole", "vendor")
		}

		c.Next()
	}
}

// CheckInboundAccess checks if user has access to specific inbound
func (m *RBACMiddleware) CheckInboundAccess(inboundIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, exists := c.Get("session")
		if !exists {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "Unauthorized",
			})
			c.Abort()
			return
		}

		sessionData := session.(gin.H)
		userID, ok := sessionData["id"].(int)
		if !ok {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "Invalid session",
			})
			c.Abort()
			return
		}

		// Admin has access to all inbounds
		isAdmin, _ := m.rbacService.IsAdmin(userID)
		if isAdmin {
			c.Next()
			return
		}

		// Get inbound ID from params or query
		var inboundIDStr string
		if inboundIDParam != "" {
			inboundIDStr = c.Param(inboundIDParam)
		} else {
			inboundIDStr = c.Query("inboundId")
		}

		if inboundIDStr == "" {
			c.JSON(400, gin.H{
				"success": false,
				"msg":     "Inbound ID required",
			})
			c.Abort()
			return
		}

		// Parse inbound ID
		var inboundID int
		_, err := fmt.Sscanf(inboundIDStr, "%d", &inboundID)
		if err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"msg":     "Invalid inbound ID",
			})
			c.Abort()
			return
		}

		// Check access
		hasAccess, err := m.rbacService.CanAccessInbound(userID, inboundID)
		if err != nil || !hasAccess {
			c.JSON(403, gin.H{
				"success": false,
				"msg":     "Access denied to this inbound",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckClientOwnership checks if user owns or can access a client
func (m *RBACMiddleware) CheckClientOwnership(clientIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, exists := c.Get("session")
		if !exists {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "Unauthorized",
			})
			c.Abort()
			return
		}

		sessionData := session.(gin.H)
		userID, ok := sessionData["id"].(int)
		if !ok {
			c.JSON(401, gin.H{
				"success": false,
				"msg":     "Invalid session",
			})
			c.Abort()
			return
		}

		// Admin can access all clients
		isAdmin, _ := m.rbacService.IsAdmin(userID)
		if isAdmin {
			c.Next()
			return
		}

		// Get client ID
		var clientIDStr string
		if clientIDParam != "" {
			clientIDStr = c.Param(clientIDParam)
		} else {
			clientIDStr = c.Query("clientId")
		}

		if clientIDStr == "" {
			c.JSON(400, gin.H{
				"success": false,
				"msg":     "Client ID required",
			})
			c.Abort()
			return
		}

		// Check ownership
		isOwner, err := m.rbacService.IsClientOwner(userID, clientIDStr)
		if err != nil || !isOwner {
			c.JSON(403, gin.H{
				"success": false,
				"msg":     "You can only modify clients you created",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// InjectUserRole adds user role to context
func (m *RBACMiddleware) InjectUserRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, exists := c.Get("session")
		if !exists {
			c.Next()
			return
		}

		sessionData := session.(gin.H)
		userID, ok := sessionData["id"].(int)
		if !ok {
			c.Next()
			return
		}

		// Get user role
		isAdmin, _ := m.rbacService.IsAdmin(userID)
		isVendor, _ := m.rbacService.IsVendor(userID)

		if isAdmin {
			c.Set("userRole", "admin")
			c.Set("isAdmin", true)
		} else if isVendor {
			c.Set("userRole", "vendor")
			c.Set("isVendor", true)
			
			// Get vendor's accessible inbounds
			inbounds, _ := m.rbacService.GetUserInbounds(userID)
			c.Set("vendorInbounds", inbounds)
		} else {
			c.Set("userRole", "user")
		}

		c.Next()
	}
}
