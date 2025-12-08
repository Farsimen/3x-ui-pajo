[English](/README.md) | [فارسی](/README.fa_IR.md) | [العربية](/README.ar_EG.md) | [中文](/README.zh_CN.md) | [Español](/README.es_ES.md) | [Русский](/README.ru_RU.md)

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./media/3x-ui-dark.png">
    <img alt="3x-ui-pajo" src="./media/3x-ui-light.png">
  </picture>
</p>

<h1 align="center">3X-UI Pajo - Multi-Vendor Edition</h1>

<p align="center">
  <strong>Advanced Xray Panel with Role-Based Access Control for Multiple Vendors</strong>
</p>

[![Release](https://img.shields.io/github/v/release/Farsimen/3x-ui-pajo.svg)](https://github.com/Farsimen/3x-ui-pajo/releases)
[![Build](https://img.shields.io/github/actions/workflow/status/Farsimen/3x-ui-pajo/release.yml.svg)](https://github.com/Farsimen/3x-ui-pajo/actions)
[![GO Version](https://img.shields.io/github/go-mod/go-version/Farsimen/3x-ui-pajo.svg)](#)
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true)](https://www.gnu.org/licenses/gpl-3.0.en.html)

---

## 🌟 What's New in Pajo Edition?

This fork adds **enterprise-grade multi-vendor management** to the powerful 3x-ui panel:

### 🔒 Role-Based Access Control (RBAC)
- **Admin Role**: Full control over panel settings, Xray configs, and all inbounds
- **Vendor Role**: Limited access to assigned inbounds with client management capabilities
- Secure permission system with database-backed access control

### 👥 Multi-Vendor Support
- Create unlimited vendor accounts with isolated access
- Each vendor can only see and manage their assigned inbounds
- Vendors can create, edit, and delete only their own clients
- Admin maintains full oversight and control

### 🔐 Enhanced Security
- Client ownership tracking - vendors can only modify clients they created
- Middleware-based route protection
- Automatic role assignment on user creation
- Data isolation between vendors

### 🎨 Simplified Vendor UI
- Clean, intuitive interface for vendors
- Hidden admin-only features (Panel Settings, Xray Configs)
- Focus on essential client management tools
- QR code and config link generation

---

## 🚀 Quick Start

### Installation

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Farsimen/3x-ui-pajo/rbac-implementation/install.sh)
```

### Default Credentials

```
Username: admin
Password: admin
Port: 2053
```

**⚠️ Change default password immediately after first login!**

---

## 📚 Documentation

### For Administrators

#### Creating a Vendor Account

1. Login as admin
2. Navigate to **Vendor Management** section
3. Click **Create New Vendor**
4. Enter username and password
5. Select inbounds to assign
6. Click **Create**

#### Managing Vendor Access

- **Grant Access**: Assign additional inbounds to existing vendors
- **Revoke Access**: Remove inbound access from vendors
- **View Vendors**: See all vendor accounts and their assigned inbounds

### For Vendors

#### What Vendors Can Do

✅ View assigned inbounds  
✅ Create new clients (users)  
✅ Edit clients they created  
✅ Delete clients they created  
✅ Generate config links and QR codes  
✅ Export client URLs  
✅ Add clients in bulk  

#### What Vendors Cannot Do

❌ Create or modify inbounds  
❌ Access Panel Settings  
❌ View Xray Configs  
❌ See other vendors' clients  
❌ Modify clients created by admin or other vendors  
❌ View server statistics  

---

## 🛠️ Technical Details

### Architecture

```
┌─────────────────────────────────┐
│         Admin Dashboard           │
│  (Full Access + Vendor Mgmt)      │
└───────────────┬─────────────────┘
                │
        ┌───────┼───────┐
        │               │
   ┌────┴────┐     ┌────┴────┐
   │ Vendor 1 │     │ Vendor 2 │
   │ Inbound A│     │ Inbound B│
   │ Inbound C│     │ Inbound D│
   └──────────┘     └──────────┘
```

### Database Schema

**New Tables:**
- `user_roles` - Maps users to their roles (admin/vendor)
- `inbound_access` - Defines which vendors can access which inbounds
- `client_ownership` - Tracks which user created each client

### Technology Stack

- **Backend**: Go (Golang)
- **Database**: SQLite with GORM
- **Frontend**: HTML5, JavaScript, Bootstrap
- **Core**: Xray-core

---

## 💻 API Endpoints

### Vendor Management (Admin Only)

```bash
# List all vendors
GET /api/vendor/list

# Create new vendor
POST /api/vendor/create
{
  "username": "vendor1",
  "password": "securepass",
  "inboundIds": [1, 2, 3]
}

# Grant inbound access
POST /api/vendor/grant
{
  "vendorId": 5,
  "inboundId": 10
}

# Revoke inbound access
POST /api/vendor/revoke
{
  "vendorId": 5,
  "inboundId": 10
}
```

### User Role Info

```bash
# Get current user's role
GET /api/user/role
```

---

## 🛡️ Security Features

### Authentication & Authorization
- Session-based authentication
- Role-based middleware protection
- Password hashing with bcrypt

### Data Isolation
- Vendors can only access their assigned inbounds
- Client ownership verification on all edit/delete operations
- Database-level access control

### UI Protection
- Dynamic menu hiding based on user role
- Disabled buttons for unauthorized actions
- Client-side and server-side validation

---

## 📝 Changelog

### v2.0.0-pajo (RBAC Edition)

**Added:**
- ✅ Role-Based Access Control (RBAC) system
- ✅ Multi-vendor support with isolated access
- ✅ Client ownership tracking
- ✅ Vendor management UI
- ✅ Enhanced security middleware
- ✅ Comprehensive documentation in Persian and English

**Changed:**
- 🔄 Updated database schema with RBAC tables
- 🔄 Modified UI to support role-based visibility
- 🔄 Enhanced API with permission checks

---

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting PRs.

### Development Setup

```bash
# Clone repository
git clone https://github.com/Farsimen/3x-ui-pajo.git
cd 3x-ui-pajo

# Checkout RBAC branch
git checkout rbac-implementation

# Build
go build -o x-ui main.go

# Run
./x-ui
```

---

## 💬 Support & Community

- **Issues**: [GitHub Issues](https://github.com/Farsimen/3x-ui-pajo/issues)
- **Documentation**: [Wiki](https://github.com/Farsimen/3x-ui-pajo/wiki)
- **Original Project**: [3x-ui by MHSanaei](https://github.com/MHSanaei/3x-ui)

---

## 📜 License

This project is licensed under the GPL V3 License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

### Based On
- [3x-ui](https://github.com/MHSanaei/3x-ui) by MHSanaei
- [x-ui](https://github.com/vaxilu/x-ui) by vaxilu

### Special Thanks
- [alireza0](https://github.com/alireza0/) - Original X-UI developer
- [MHSanaei](https://github.com/MHSanaei/) - 3x-ui maintainer
- All contributors to the Xray and V2Ray projects

---

## ⚠️ Disclaimer

**This project is for personal and educational use only.**  
Do not use it for illegal purposes.  
Do not use it in production without proper security auditing.

The developers are not responsible for any misuse of this software.

---

<p align="center">
  Made with ❤️ for the community
</p>

<p align="center">
  <a href="#top">⬆️ Back to Top</a>
</p>
