# 📘 راهنمای مدیریت Vendor در 3X-UI Pajo

## 🌟 نسخه RBAC چیست؟

3X-UI Pajo نسخه توسعه‌یافته‌ای از 3X-UI است که قابلیت **Multi-Vendor** و **Role-Based Access Control (RBAC)** را اضافه کرده است.

### ویژگی‌های اصلی:

✅ **مدیریت چند فروشنده (Vendor)**: تا 50+ فروشنده مستقل  
✅ **کنترل دسترسی به Inbound**: هر Vendor فقط به Inbound های خودش دسترسی دارد  
✅ **ردیابی مالکیت Client**: هر Client متعلق به یک Vendor است  
✅ **داشبورد مجزا**: هر Vendor فقط اطلاعات خودش را می‌بیند  
✅ **مدیریت متمرکز Admin**: Admin می‌تواند همه Vendor ها را مدیریت کند

---

## 🚀 نصب

### نصب سریع:

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Farsimen/3x-ui-pajo/rbac-implementation/install.sh)
```

### گزینه‌های تعاملی:

1. **پورت پنل** (پیش‌فرض: 2053)
2. **نام کاربری Admin**
3. **رمز عبور Admin**
4. **مسیر نصب** (پیش‌فرض: /usr/local/x-ui)

---

## 👤 ایجاد Vendor جدید (Admin)

### از طریق API:

```bash
curl -X POST http://YOUR_IP:PORT/vendor/create \
  -H "Content-Type: application/json" \
  -d '{
    "username": "vendor1",
    "password": "securepass123",
    "inboundIds": [1, 2, 3]
  }'
```

### پارامترها:

- `username`: نام کاربری Vendor (حداقل 3 کاراکتر)
- `password`: رمز عبور (حداقل 4 کاراکتر)
- `inboundIds`: آرایه از ID های Inbound که Vendor به آن‌ها دسترسی دارد

---

## 📋 مدیریت Vendor ها

### لیست همه Vendor ها:

```bash
curl -X GET http://YOUR_IP:PORT/vendor/list
```

**خروجی نمونه:**

```json
[
  {
    "id": 2,
    "username": "vendor1",
    "role": "vendor",
    "inbounds": [1, 2, 3]
  },
  {
    "id": 3,
    "username": "vendor2",
    "role": "vendor",
    "inbounds": [4, 5]
  }
]
```

### اعطای دسترسی به Inbound:

```bash
curl -X POST http://YOUR_IP:PORT/vendor/grant \
  -H "Content-Type: application/json" \
  -d '{
    "vendorId": 2,
    "inboundId": 6
  }'
```

### لغو دسترسی به Inbound:

```bash
curl -X POST http://YOUR_IP:PORT/vendor/revoke \
  -H "Content-Type: application/json" \
  -d '{
    "vendorId": 2,
    "inboundId": 1
  }'
```

### حذف Vendor:

```bash
curl -X DELETE http://YOUR_IP:PORT/vendor/delete/2
```

---

## 🔐 احراز هویت و Session

### ورود (Login):

```bash
curl -X POST http://YOUR_IP:PORT/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "vendor1",
    "password": "securepass123"
  }'
```

### دریافت نقش کاربر:

```bash
curl -X POST http://YOUR_IP:PORT/getUserRole
```

**خروجی نمونه برای Vendor:**

```json
{
  "role": "vendor",
  "isAdmin": false,
  "isVendor": true,
  "inbounds": [1, 2, 3],
  "username": "vendor1"
}
```

---

## 🎯 محدودیت‌های Vendor

### دسترسی‌های Vendor:

✅ مشاهده و مدیریت Client های خودش  
✅ مشاهده آمار ترافیک Client های خودش  
✅ ایجاد/حذف/ویرایش Client ها در Inbound های مجاز  
✅ مشاهده لیست Inbound های مجاز خودش

### محدودیت‌های Vendor:

❌ دسترسی به Inbound های سایر Vendor ها  
❌ مشاهده/ویرایش Client های سایر Vendor ها  
❌ ایجاد یا حذف Inbound ها  
❌ تغییر تنظیمات پنل  
❌ مدیریت Vendor های دیگر  
❌ دسترسی به تنظیمات Xray

---

## 📊 ساختار دیتابیس

### جداول RBAC:

#### 1. `user_roles` - نقش‌های کاربر

| فیلد | نوع | توضیح |
|------|-----|-------|
| id | INTEGER | شناسه منحصر به فرد |
| user_id | INTEGER | شناسه کاربر |
| role | VARCHAR | نقش (admin/vendor) |
| created_at | DATETIME | تاریخ ایجاد |

#### 2. `inbound_accesses` - دسترسی به Inbound

| فیلد | نوع | توضیح |
|------|-----|-------|
| id | INTEGER | شناسه منحصر به فرد |
| user_id | INTEGER | شناسه کاربر |
| inbound_id | INTEGER | شناسه Inbound |
| created_at | DATETIME | تاریخ ایجاد |

---

## 🛠️ دستورات مفید

### مدیریت سرویس:

```bash
# شروع پنل
systemctl start x-ui

# توقف پنل
systemctl stop x-ui

# ریستارت پنل
systemctl restart x-ui

# مشاهده وضعیت
systemctl status x-ui

# مشاهده لاگ‌ها
journalctl -u x-ui -f
```

### دستورات CLI:

```bash
# منوی کنترل
x-ui

# تغییر رمز عبور
x-ui setting -username admin -password newpass

# تغییر پورت
x-ui setting -port 2087
```

---

## 🔧 عیب‌یابی

### مشکل: Xray متوقف می‌شود

```bash
# چک کردن لاگ Xray
journalctl -u x-ui -n 100 | grep -i xray

# چک کردن فایل‌های GeoIP
ls -la /usr/local/x-ui/bin/*.dat

# ریستارت دستی
systemctl restart x-ui
```

### مشکل: Vendor نمی‌تواند Client ایجاد کند

```bash
# چک کردن دسترسی‌های Vendor
curl -X GET http://localhost:2053/vendor/inbounds/VENDOR_ID

# اعطای دسترسی
curl -X POST http://localhost:2053/vendor/grant \
  -d '{"vendorId": 2, "inboundId": 1}'
```

### مشکل: پنل بالا نمی‌آید

```bash
# چک کردن پورت
ss -tulpn | grep :2053

# چک کردن فایروال
ufw status
ufw allow 2053/tcp

# چک کردن سرویس systemd
systemctl status x-ui
```

---

## 🌐 API Endpoints

### Admin-only Endpoints:

| متد | مسیر | توضیح |
|------|------|-------|
| POST | `/vendor/create` | ایجاد Vendor جدید |
| GET | `/vendor/list` | لیست همه Vendor ها |
| POST | `/vendor/grant` | اعطای دسترسی |
| POST | `/vendor/revoke` | لغو دسترسی |
| DELETE | `/vendor/delete/:id` | حذف Vendor |
| GET | `/vendor/inbounds/:id` | لیست Inbound های Vendor |

### Public Endpoints:

| متد | مسیر | توضیح |
|------|------|-------|
| POST | `/login` | ورود به سیستم |
| GET | `/logout` | خروج از سیستم |
| POST | `/getUserRole` | دریافت نقش کاربر |

---

## 🔒 امنیت

### توصیه‌های امنیتی:

1. **رمز عبور قوی**: حداقل 12 کاراکتر با ترکیب حروف، اعداد و نمادها
2. **فایروال**: فقط پورت‌های لازم را باز کنید
3. **HTTPS**: حتماً SSL فعال کنید
4. **به‌روزرسانی**: به طور منظم پنل را به‌روز کنید
5. **Backup**: از دیتابیس به طور منظم پشتیبان بگیرید

### فعال کردن فایروال:

```bash
ufw enable
ufw allow 22/tcp    # SSH
ufw allow 2053/tcp  # Panel
ufw allow 443/tcp   # Xray
ufw reload
```

---

## 📚 مثال کامل: راه‌اندازی سیستم Multi-Vendor

### گام 1: نصب پنل

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Farsimen/3x-ui-pajo/rbac-implementation/install.sh)
```

### گام 2: ایجاد Inbound های مختلف

```bash
# وارد پنل شوید و 5 Inbound مختلف بسازید:
# - Inbound 1: VLESS Reality (پورت 443)
# - Inbound 2: VMess WS (پورت 80)
# - Inbound 3: Trojan (پورت 8443)
# - Inbound 4: VLESS WS TLS (پورت 2096)
# - Inbound 5: Shadowsocks (پورت 9000)
```

### گام 3: ایجاد Vendor اول (دسترسی به Inbound 1 و 2)

```bash
curl -X POST http://YOUR_IP:2053/vendor/create \
  -H "Content-Type: application/json" \
  -d '{
    "username": "vendor_tehran",
    "password": "Tehran@2025!",
    "inboundIds": [1, 2]
  }'
```

### گام 4: ایجاد Vendor دوم (دسترسی به Inbound 3 و 4)

```bash
curl -X POST http://YOUR_IP:2053/vendor/create \
  -H "Content-Type: application/json" \
  -d '{
    "username": "vendor_mashhad",
    "password": "Mashhad@2025!",
    "inboundIds": [3, 4]
  }'
```

### گام 5: ایجاد Vendor سوم (دسترسی به Inbound 5)

```bash
curl -X POST http://YOUR_IP:2053/vendor/create \
  -H "Content-Type: application/json" \
  -d '{
    "username": "vendor_shiraz",
    "password": "Shiraz@2025!",
    "inboundIds": [5]
  }'
```

### گام 6: تست دسترسی Vendor

```bash
# Login به عنوان vendor_tehran
curl -X POST http://YOUR_IP:2053/login \
  -d '{"username": "vendor_tehran", "password": "Tehran@2025!"}'

# چک کردن نقش
curl -X POST http://YOUR_IP:2053/getUserRole
# Output: {"role":"vendor", "inbounds":[1,2]}
```

### گام 7: اعطای دسترسی اضافی

```bash
# فرض کنیم می‌خواهیم Inbound 3 را هم به vendor_tehran بدهیم
curl -X POST http://YOUR_IP:2053/vendor/grant \
  -d '{"vendorId": 2, "inboundId": 3}'
```

---

## 🎓 سناریوهای واقعی

### سناریو 1: فروشنده محلی

**نیاز**: شما یک سرویس‌دهنده هستید و می‌خواهید به 10 فروشنده محلی دسترسی بدهید.

**راه‌حل**:
1. برای هر فروشنده یک Inbound اختصاصی بسازید
2. Vendor ایجاد کنید و دسترسی به همان Inbound را بدهید
3. هر فروشنده فقط Client های خودش را می‌بیند

### سناریو 2: فروش عمده

**نیاز**: 5 فروشنده عمده دارید، هر کدام نیاز به چند Inbound دارند.

**راه‌حل**:
1. برای هر فروشنده چند Inbound بسازید
2. در هنگام ایجاد Vendor، آرایه‌ای از Inbound ها رو بدهید
3. Vendor می‌تواند در همه Inbound های خودش Client ایجاد کند

### سناریو 3: تست و Development

**نیاز**: می‌خواهید یک Vendor تستی داشته باشید.

**راه‌حل**:
```bash
# ایجاد Inbound تستی
# ایجاد Vendor تستی
curl -X POST http://localhost:2053/vendor/create \
  -d '{"username":"test_vendor", "password":"test123", "inboundIds":[99]}'

# تست کردن
# حذف Vendor
curl -X DELETE http://localhost:2053/vendor/delete/VENDOR_ID
```

---

## 🆘 پشتیبانی

### منابع:

- 🌐 **GitHub**: [https://github.com/Farsimen/3x-ui-pajo](https://github.com/Farsimen/3x-ui-pajo)
- 📖 **Documentation**: [IMPLEMENTATION_STEPS.md](IMPLEMENTATION_STEPS.md)
- 🔧 **Issues**: [GitHub Issues](https://github.com/Farsimen/3x-ui-pajo/issues)

### گزارش باگ:

اگر باگی پیدا کردید:

1. لاگ‌ها را جمع‌آوری کنید:
   ```bash
   journalctl -u x-ui -n 200 > x-ui-logs.txt
   systemctl status x-ui > x-ui-status.txt
   ```

2. Issue جدید در GitHub باز کنید

3. اطلاعات زیر را ضمیمه کنید:
   - نسخه سیستم‌عامل
   - نسخه 3X-UI Pajo
   - توضیح دقیق مشکل
   - لاگ‌ها

---

## 📝 License

3X-UI Pajo تحت لایسنس GPL-3.0 منتشر شده است.

---

## ❤️ تشکر

- تیم [3x-ui](https://github.com/mhsanaei/3x-ui) برای پنل اصلی
- جامعه ایرانی برای حمایت و feedback

---

**Made with ❤️ by [Farsimen](https://github.com/Farsimen)**
