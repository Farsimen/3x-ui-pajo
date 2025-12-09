# 3X-UI Pajo - Multi-Vendor Panel با RBAC

<div dir="rtl">

## 📋 درباره پروژه

نسخه توسعه یافته 3X-UI با قابلیت **مدیریت چند فروشنده** (Multi-Vendor) و **کنترل دسترسی مبتنی بر نقش** (RBAC).

با این پنل می‌توانید:
- ✅ تا 50 فروشنده همزمان داشته باشید
- ✅ به هر فروشنده فقط inbound های خاصی دسترسی بدهید
- ✅ هر فروشنده فقط client های inbound های خودش را ببیند
- ✅ مدیریت کامل vendor ها از طریق API

## 🚀 نصب سریع

### پیش‌نیازها
- سرور Ubuntu 20.04+ یا Debian 11+
- دسترسی root
- حداقل 1GB RAM

### نصب با یک دستور

```bash
bash <(curl -Ls https://raw.githubusercontent.com/Farsimen/3x-ui-pajo/rbac-implementation/install_rbac.sh)
```

نصب حدود 5-10 دقیقه طول می‌کشد.

### اطلاعات پیش‌فرض

```
URL: http://YOUR_IP:2087
Username: admin
Password: admin
```

⚠️ **مهم:** بلافاصله پس از نصب، پسورد را تغییر دهید!

## 📚 مستندات API

### ساخت Vendor جدید

```bash
curl -X POST http://YOUR_IP:2087/BASE_PATH/vendor/create \
  -H "Content-Type: application/json" \
  -H "Cookie: 3x-ui=YOUR_SESSION" \
  -d '{
    "username": "vendor1",
    "password": "secure_password",
    "inboundIds": [1, 2, 3]
  }'
```

### لیست Vendor ها

```bash
curl http://YOUR_IP:2087/BASE_PATH/vendor/list \
  -H "Cookie: 3x-ui=YOUR_SESSION"
```

### حذف Vendor

```bash
curl -X DELETE http://YOUR_IP:2087/BASE_PATH/vendor/delete/2 \
  -H "Cookie: 3x-ui=YOUR_SESSION"
```

### دادن دسترسی به Inbound

```bash
curl -X POST http://YOUR_IP:2087/BASE_PATH/vendor/grant \
  -H "Content-Type: application/json" \
  -H "Cookie: 3x-ui=YOUR_SESSION" \
  -d '{
    "vendorId": 2,
    "inboundId": 5
  }'
```

### گرفتن دسترسی از Inbound

```bash
curl -X POST http://YOUR_IP:2087/BASE_PATH/vendor/revoke \
  -H "Content-Type: application/json" \
  -H "Cookie: 3x-ui=YOUR_SESSION" \
  -d '{
    "vendorId": 2,
    "inboundId": 5
  }'
```

## 🔧 مدیریت سرویس

```bash
# مشاهده وضعیت
systemctl status x-ui

# ری‌استارت
systemctl restart x-ui

# مشاهده لاگ‌ها
journalctl -u x-ui -f

# توقف سرویس
systemctl stop x-ui

# شروع سرویس
systemctl start x-ui
```

## 📊 ساختار دیتابیس

### جدول `user_roles`
| فیلد | نوع | توضیحات |
|------|-----|----------|
| id | INTEGER | شناسه یکتا |
| user_id | INTEGER | شناسه کاربر |
| role | VARCHAR(20) | نقش (admin/vendor) |

### جدول `inbound_access`
| فیلد | نوع | توضیحات |
|------|-----|----------|
| id | INTEGER | شناسه یکتا |
| user_id | INTEGER | شناسه فروشنده |
| inbound_id | INTEGER | شناسه inbound |

## 🛡️ امنیت

- ✅ رمز عبور با bcrypt hash می‌شود
- ✅ Session-based authentication
- ✅ RBAC برای جداسازی دسترسی‌ها
- ✅ هر vendor فقط inbound های خودش را می‌بیند

## 🐛 عیب‌یابی

### پنل باز نمیشه

```bash
# چک کردن وضعیت
systemctl status x-ui

# چک کردن لاگ
journalctl -u x-ui -n 50

# چک کردن پورت
ss -tulpn | grep 2087
```

### API کار نمی‌کنه

1. مطمئن شوید لاگین کرده‌اید و cookie دارید
2. BASE_PATH را صحیح وارد کنید
3. لاگ‌های سرور را چک کنید

### Build خطا می‌ده

```bash
# نصب dependencies
apt-get install -y gcc build-essential

# Build با CGO
CGO_ENABLED=1 go build -o x-ui main.go
```

## 📞 پشتیبانی

- GitHub Issues: [Farsimen/3x-ui-pajo](https://github.com/Farsimen/3x-ui-pajo/issues)
- Telegram: در صورت نیاز

## 📝 لایسنس

GPL-3.0 - مطابق پروژه اصلی 3X-UI

## 🙏 قدردانی

- [MHSanaei/3x-ui](https://github.com/MHSanaei/3x-ui) - پروژه اصلی
- تمام contributor های 3X-UI

---

**ساخته شده با ❤️ برای جامعه ایرانی**

</div>
