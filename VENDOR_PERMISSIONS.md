# دسترسی‌های Vendor (فروشنده)

بر اساس تحلیل UI و نیازمندی‌های شما، دسترسی‌های فروشندگان به شرح زیر است:

## ✅ دسترسی‌های مجاز برای Vendor:

### 1. مشاهده Inbounds
- فروشنده فقط Inbound‌هایی که توسط Admin به او اختصاص داده شده را می‌بیند
- نمایش اطلاعات اساسی: ID, Menu, Status (Enabled/Disabled), Remark

### 2. مدیریت Clients (کاربران)
فروشنده می‌تواند برای Inbound‌های خودش:

#### عملیات مجاز:
- ✅ **Add Client**: ساخت کاربر جدید
- ✅ **Add Bulk**: ساخت کاربران به صورت دسته‌ای
- ✅ **Export All URLs**: دریافت لینک‌های کانفیگ همه کاربران
- ✅ **Export Inbound**: دریافت فایل کانفیگ Inbound
- ✅ **مشاهده لیست Clients**: مشاهده کاربران هر Inbound

#### عملیات روی هر Client:
- ✅ **QR Code**: مشاهده QR کد
- ✅ **Edit Client**: ویرایش کاربر (فقط کاربرانی که خودش ساخته)
- ✅ **Info**: مشاهده اطلاعات کاربر
- ✅ **Copy**: کپی لینک کانفیگ
- ✅ **Delete**: حذف کاربر (فقط کاربرانی که خودش ساخته)

## ❌ دسترسی‌های غیرمجاز برای Vendor:

### 1. مدیریت Inbound
- ❌ **Add Inbound**: ساخت Inbound جدید (فقط Admin)
- ❌ **Edit Inbound**: ویرایش تنظیمات Inbound (فقط Admin)
- ❌ **Delete Inbound**: حذف Inbound (فقط Admin)
- ❌ **General Actions**: عملیات کلی روی Inbound‌ها (فقط Admin)
- ❌ **Clone Inbound**: کپی کردن Inbound (فقط Admin)
- ❌ **Reset Traffic**: ریست کردن ترافیک Inbound (فقط Admin)

### 2. تنظیمات پنل
- ❌ **Panel Settings**: دسترسی به تنظیمات پنل
- ❌ **Xray Configs**: دسترسی به تنظیمات Xray
- ❌ **Theme**: تغییر تم

### 3. آمار کلی
- ❌ **Total Sent/Received**: مشاهده آمار کل سرور
- ❌ **Total Usage**: مشاهده مصرف کل سرور

### 4. عملیات خطرناک
- ❌ **Delete Depleted Clients**: حذف کاربران تمام‌شده
- ❌ **Reset Clients Traffic**: ریست کردن ترافیک کاربران به صورت دسته‌ای

## 🔒 محدودیت‌های امنیتی:

### 1. Data Isolation
- هر Vendor فقط Inbound‌های اختصاص داده شده به خودش را می‌بیند
- هر Vendor فقط Client‌هایی که خودش ساخته را می‌تواند ویرایش/حذف کند
- Client‌های ساخته شده توسط Admin یا Vendor دیگر قابل ویرایش نیستند

### 2. UI Restrictions
- تب‌های زیر از منوی سمت چپ مخفی می‌شوند:
  - Theme
  - Overview (آمار کلی)
  - Panel Settings
  - Xray Configs
  - Log Out (نمایش داده می‌شود)

- دکمه‌های زیر از صفحه Inbounds مخفی می‌شوند:
  - Add Inbound
  - General Actions
  - دکمه Edit روی Inbound (آیکون ویرایش)
  - دکمه Delete روی Inbound (آیکون حذف)
  - دکمه Clone روی Inbound
  - دکمه Reset Traffic روی Inbound

- دکمه‌های زیر از منوی Client نمایش داده می‌شود:
  - Add Client ✅
  - Add Bulk ✅
  - Export All URLs ✅
  - Export Inbound ✅

- دکمه‌های زیر از منوی Client مخفی می‌شوند:
  - Reset Clients Traffic ❌
  - Delete Depleted Clients ❌

### 3. API Restrictions
- همه API endpoint‌های Vendor با middleware احراز هویت محافظت می‌شوند
- هر درخواست API بررسی می‌شود که آیا Vendor به آن Inbound دسترسی دارد
- هر درخواست ویرایش/حذف Client بررسی می‌شود که آیا Client توسط همین Vendor ساخته شده

## 📋 جدول مقایسه دسترسی‌ها:

| عملیات | Admin | Vendor |
|---------|-------|--------|
| مشاهده همه Inbound‌ها | ✅ | ❌ (فقط اختصاصی) |
| ساخت Inbound | ✅ | ❌ |
| ویرایش Inbound | ✅ | ❌ |
| حذف Inbound | ✅ | ❌ |
| ساخت Client | ✅ | ✅ |
| ویرایش Client | ✅ | ✅ (فقط خودش) |
| حذف Client | ✅ | ✅ (فقط خودش) |
| مشاهده QR / Config | ✅ | ✅ |
| دسترسی به Panel Settings | ✅ | ❌ |
| دسترسی به Xray Configs | ✅ | ❌ |
| مشاهده آمار کلی | ✅ | ❌ |
| Reset Traffic | ✅ | ❌ |
| Clone Inbound | ✅ | ❌ |

## 🔄 نحوه پیاده‌سازی:

### Backend:
1. اضافه کردن فیلد `created_by` به Client model
2. بررسی مالکیت Client در تمام API‌های ویرایش/حذف
3. فیلتر کردن Inbound‌ها بر اساس دسترسی‌های Vendor

### Frontend:
1. دریافت نقش کاربر از API
2. مخفی کردن المان‌های UI بر اساس نقش
3. غیرفعال کردن دکمه‌های عملیات ممنوع

---

**نکته:** این مستند بر اساس تحلیل UI ارسالی شما تهیه شده است.
