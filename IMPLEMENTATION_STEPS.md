# مراحل دقیق پیاده‌سازی RBAC

## وضعیت فعلی ✅

تا اینجا این کارها انجام شده:

1. ✅ Database models: `UserRole`, `InboundAccess`
2. ✅ Database migration و seeder برای assign کردن admin role
3. ✅ RBAC Service برای مدیریت نقش‌ها و دسترسی‌ها
4. ✅ مستندسازی دقیق دسترسی‌های Vendor

---

## مرحله 1: تغییر ساختار Client برای ذخیره مالکیت

### مسئله:
در حال حاضر Client به صورت JSON در `Inbound.Settings` ذخیره می‌شود و جدول جداگانه‌ای ندارد.

### راه‌حل‌ها:

#### گزینه A (توصیه شده): جدول جدید `client_ownership`

```go
// database/model/rbac.go را آپدیت کنید:

// ClientOwnership tracks which user created each client
type ClientOwnership struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	InboundId   int    `json:"inboundId" gorm:"index;not null"`  // Foreign key to Inbound
	ClientEmail string `json:"clientEmail" gorm:"index;not null"` // Client's email (unique identifier)
	CreatedBy   int    `json:"createdBy" gorm:"index;not null"`   // User ID who created this client
	CreatedAt   int64  `json:"createdAt"`                          // Timestamp
}

func (ClientOwnership) TableName() string {
	return "client_ownership"
}
```

**مزایا:**
- بدون تغییر ساختار فعلی Client
- ساده و قابل نگهداری
- برای clientهای قدیمی هیچ مشکلی ایجاد نمی‌کند

#### گزینه B: فیلد `CreatedBy` به `Client` struct

```go
// در Client struct:
CreatedBy int `json:"createdBy,omitempty"` // User ID who created this client
```

**معایب:**
- Client در JSON ذخیره می‌شود، باید هر بار parse و update شود
- پیچیده‌تر برای مدیریت

**توصیه: گزینه A را انتخاب کنید**

---

## مرحله 2: آپدیت database/db.go

```go
// در initModels() بعد از &model.InboundAccess{} اضافه کنید:
&model.ClientOwnership{},
```

---

## مرحله 3: آپدیت RBAC Service

در `web/service/rbac.go` توابع زیر را اضافه کنید:

```go
// RecordClientOwnership records that a user created a client
func (s *RBACService) RecordClientOwnership(inboundId int, clientEmail string, userId int) error {
	db := database.GetDB()
	
	ownership := &model.ClientOwnership{
		InboundId:   inboundId,
		ClientEmail: clientEmail,
		CreatedBy:   userId,
		CreatedAt:   time.Now().Unix(),
	}
	
	return db.Create(ownership).Error
}

// CanEditClient checks if a user can edit/delete a specific client
func (s *RBACService) CanEditClient(userId int, inboundId int, clientEmail string) (bool, error) {
	// Admin can edit all clients
	isAdmin, err := s.IsAdmin(userId)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}
	
	// Vendor can only edit clients they created
	db := database.GetDB()
	var ownership model.ClientOwnership
	err = db.Where("inbound_id = ? AND client_email = ? AND created_by = ?", 
		inboundId, clientEmail, userId).First(&ownership).Error
	
	if database.IsNotFound(err) {
		return false, nil // Not owned by this user
	}
	
	if err != nil {
		return false, err
	}
	
	return true, nil
}

// GetClientsByVendor returns all clients created by a vendor for a specific inbound
func (s *RBACService) GetClientsByVendor(userId int, inboundId int) ([]string, error) {
	db := database.GetDB()
	var ownerships []model.ClientOwnership
	
	err := db.Where("inbound_id = ? AND created_by = ?", inboundId, userId).Find(&ownerships).Error
	if err != nil {
		return nil, err
	}
	
	clients := make([]string, len(ownerships))
	for i, o := range ownerships {
		clients[i] = o.ClientEmail
	}
	
	return clients, nil
}
```

---

## مرحله 4: آپدیت Inbound Controller

باید فایل `web/controller/inbound.go` را پیدا کنید و این تغییرات را اعمال کنید:

### 4.1: در تابع Add Client:

```go
func (a *InboundController) addInboundClient(c *gin.Context) {
	// ... کد فعلی ...
	
	// بعد از ساخت client موفقیت‌آمیز، مالکیت را ثبت کنید:
	userIdInt := session.GetLoginUser(c)
	if userIdInt != nil {
		rbacService := &service.RBACService{}
		_ = rbacService.RecordClientOwnership(inboundId, client.Email, *userIdInt)
	}
	
	// ...
}
```

### 4.2: در تابع Update Client:

```go
func (a *InboundController) updateInboundClient(c *gin.Context) {
	inboundId := // ... دریافت inbound ID
	clientEmail := // ... دریافت client email
	
	userIdInt := session.GetLoginUser(c)
	if userIdInt == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "Unauthorized"})
		return
	}
	
	// بررسی دسترسی
	rbacService := &service.RBACService{}
	canEdit, err := rbacService.CanEditClient(*userIdInt, inboundId, clientEmail)
	if err != nil || !canEdit {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "msg": "Access denied - You can only edit clients you created"})
		return
	}
	
	// ادامه عملیات update...
}
```

### 4.3: در تابع Delete Client:

```go
func (a *InboundController) delInboundClient(c *gin.Context) {
	inboundId := // ...
	clientEmail := // ...
	
	userIdInt := session.GetLoginUser(c)
	if userIdInt == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "Unauthorized"})
		return
	}
	
	// بررسی دسترسی
	rbacService := &service.RBACService{}
	canEdit, err := rbacService.CanEditClient(*userIdInt, inboundId, clientEmail)
	if err != nil || !canEdit {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "msg": "Access denied"})
		return
	}
	
	// حذف client...
	// سپس حذف ownership record:
	db := database.GetDB()
	db.Where("inbound_id = ? AND client_email = ?", inboundId, clientEmail).Delete(&model.ClientOwnership{})
}
```

### 4.4: در تابع List Clients:

```go
func (a *InboundController) getInboundClients(c *gin.Context) {
	inboundId := // ...
	
	userIdInt := session.GetLoginUser(c)
	if userIdInt == nil {
		return
	}
	
	// دریافت clients از inbound settings
	clients := // ... parse JSON
	
	// بررسی نقش برای فیلتر کردن
	rbacService := &service.RBACService{}
	isAdmin, _ := rbacService.IsAdmin(*userIdInt)
	
	if !isAdmin {
		// Vendor فقط clients خودش را می‌بیند
		ownedEmails, _ := rbacService.GetClientsByVendor(*userIdInt, inboundId)
		ownedMap := make(map[string]bool)
		for _, email := range ownedEmails {
			ownedMap[email] = true
		}
		
		// فیلتر کردن clients
		filteredClients := []Client{}
		for _, client := range clients {
			if ownedMap[client.Email] {
				filteredClients = append(filteredClients, client)
			}
		}
		clients = filteredClients
	}
	
	// بازگشت clients...
}
```

---

## مرحله 5: ایجاد Middleware برای محافظت از Routes

فایل `web/middleware/rbac.go` را بسازید:

```go
package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/mhsanaei/3x-ui/v2/web/service"
	"github.com/mhsanaei/3x-ui/v2/web/session"
)

// RequireAdmin ensures only admin users can access the route
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInt := session.GetLoginUser(c)
		if userIdInt == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "Unauthorized"})
			c.Abort()
			return
		}

		rbacService := &service.RBACService{}
		isAdmin, err := rbacService.IsAdmin(*userIdInt)
		if err != nil || !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "msg": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireVendorOrAdmin allows both vendor and admin users
func RequireVendorOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInt := session.GetLoginUser(c)
		if userIdInt == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("userId", *userIdInt)
		c.Next()
	}
}
```

---

## مرحله 6: آپدیت Routes در web/web.go

```go
// پیدا کردن بخش routes و آپدیت:

// Admin-only routes
adminGroup := g.Group("/panel")
adminGroup.Use(middleware.RequireAdmin())
{
	adminGroup.POST("/setting/update", xui.updateSetting)
	adminGroup.GET("/setting/all", xui.getAllSettings)
	// ... سایر admin routes
}

// Vendor and Admin routes
inboundGroup := g.Group("/panel/inbound")
inboundGroup.Use(middleware.RequireVendorOrAdmin())
{
	inboundGroup.POST("/list", xui.inbounds)
	inboundGroup.GET("/get/:id", xui.getInbound)
	
	// Client management
	inboundGroup.POST("/addClient", xui.addInboundClient)
	inboundGroup.POST("/updateClient/:id", xui.updateInboundClient)
	inboundGroup.POST("/delClient/:id", xui.delInboundClient)
	// ...
}

// Admin-only inbound operations
adminInboundGroup := g.Group("/panel/inbound")
adminInboundGroup.Use(middleware.RequireAdmin())
{
	adminInboundGroup.POST("/add", xui.addInbound)
	adminInboundGroup.POST("/update/:id", xui.updateInbound)
	adminInboundGroup.POST("/del/:id", xui.delInbound)
	// ...
}
```

---

## مرحله 7: تغییرات Frontend

### 7.1: اضافه کردن API برای دریافت نقش کاربر

در controller جدید `web/controller/user.go`:

```go
func (a *UserController) getUserRole(c *gin.Context) {
	userIdInt := session.GetLoginUser(c)
	if userIdInt == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false})
		return
	}
	
	rbacService := &service.RBACService{}
	role, err := rbacService.GetUserRole(*userIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"success": true, "role": role})
}
```

Route:
```go
g.GET("/panel/user/role", userController.getUserRole)
```

### 7.2: مخفی کردن UI Elements

در فایل HTML اصلی (مثلا `web/html/xui/index.html`):

```javascript
// دریافت نقش کاربر
fetch('/panel/user/role')
  .then(res => res.json())
  .then(data => {
    if (data.success && data.role === 'vendor') {
      hideVendorRestrictedUI();
    }
  });

function hideVendorRestrictedUI() {
  // مخفی کردن منوهای سمت چپ
  const elementsToHide = [
    '#theme-menu',           // Theme
    '#overview-menu',        // Overview
    '#panel-settings-menu',  // Panel Settings
    '#xray-configs-menu',    // Xray Configs
  ];
  
  elementsToHide.forEach(selector => {
    const el = document.querySelector(selector);
    if (el) el.style.display = 'none';
  });
  
  // مخفی کردن دکمه‌های صفحه Inbounds
  const inboundButtons = [
    '#add-inbound-btn',        // Add Inbound
    '#general-actions-btn',    // General Actions
    '.inbound-edit-btn',       // Edit buttons
    '.inbound-delete-btn',     // Delete buttons
    '.inbound-clone-btn',      // Clone buttons
    '.inbound-reset-traffic',  // Reset traffic buttons
  ];
  
  inboundButtons.forEach(selector => {
    document.querySelectorAll(selector).forEach(el => {
      el.style.display = 'none';
    });
  });
  
  // مخفی کردن آمار کلی
  document.querySelector('#total-stats')?.style.display = 'none';
}
```

### 7.3: محدود کردن عملیات Client

در صفحه لیست Clients:

```javascript
// بعد از دریافت لیست clients:
if (userRole === 'vendor') {
  // مخفی کردن دکمه‌های Edit/Delete برای clientهایی که خودش نساخته
  clients.forEach(client => {
    if (!client.canEdit) { // سرور باید flag بفرسته
      document.querySelector(`#edit-${client.email}`)?.remove();
      document.querySelector(`#delete-${client.email}`)?.remove();
    }
  });
}
```

---

## مرحله 8: ایجاد صفحه مدیریت Vendor (فقط برای Admin)

صفحه‌ای که Admin بتواند:
- لیست vendorها را ببیند
- vendor جدید بسازد
- inbound به vendor اختصاص دهد
- دسترسی vendor را لغو کند

(کد در `RBAC_IMPLEMENTATION_GUIDE.md` موجود است)

---

## مرحله 9: Testing

### 9.1: تست به عنوان Admin:
- باید همه چیز قابل مشاهده و ویرایش باشد
- بتواند vendor بسازد

### 9.2: تست به عنوان Vendor:
- فقط inboundهای اختصاصی را ببیند
- نتواند inbound بسازد/ویرایش کند
- بتواند client بسازد
- فقط clientهای خودش را ویرایش/حذف کند

### 9.3: تست امنیتی:
- سعی کنید با API مستقیم client دیگران را حذف کنید (Postman/curl)
- باید 403 Forbidden برگرداند

---

## 📦 Checklist نهایی:

- [ ] جدول `client_ownership` ایجاد شده
- [ ] RBAC Service آپدیت شده
- [ ] Controllerها بررسی مالکیت می‌کنند
- [ ] Middleware ایجاد شده
- [ ] Routes محافظت شده‌اند
- [ ] Frontend بر اساس نقش UI را محدود می‌کند
- [ ] صفحه مدیریت Vendor ساخته شده
- [ ] تست شده به عنوان Admin
- [ ] تست شده به عنوان Vendor
- [ ] تست امنیتی انجام شده

---

## 🔥 نکات مهم:

1. **حتماً قبل از deploy بر production، backup از database بگیرید**
2. **Migration روی دیتابیس قدیمی clientهایی که قبلاً ساخته شده‌اند ownership ندارند - می‌توانید برای آنها created_by=1 (admin) بگذارید**
3. **تست دقیق بر روی development environment**

---

**آماده شروع؟** بگو از کدوم مرحله شروع کنیم! 🚀
