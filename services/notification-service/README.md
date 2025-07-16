# Notification Service - LUNARY (В разработке)

## 📝 Описание
Сервис для отправки push-уведомлений, email-рассылок и SMS. Обрабатывает все виды уведомлений в системе LUNARY.

## ⚠️ Статус
**В разработке** - сервис еще не реализован, только заглушка

## 🚀 Планируемые эндпоинты

### POST /push/send
Отправка push-уведомления

### POST /email/send
Отправка email уведомления

### POST /sms/send
Отправка SMS уведомления

### GET /notifications/{user_id}
Получение истории уведомлений пользователя

## 📋 Планируемая Postman Collection

### 1. Отправка push-уведомления
```json
{
  "method": "POST",
  "url": "http://localhost:8086/push/send",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}",
    "Content-Type": "application/json"
  },
  "body": {
    "user_id": 123,
    "title": "Новый сценарий!",
    "message": "Найден идеальный сценарий для вашего мероприятия",
    "data": {
      "scenario_id": 456,
      "category": "family"
    }
  }
}
```

### 2. Отправка email
```json
{
  "method": "POST",
  "url": "http://localhost:8086/email/send",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}",
    "Content-Type": "application/json"
  },
  "body": {
    "to": "user@example.com",
    "subject": "Добро пожаловать в LUNARY!",
    "template": "welcome",
    "data": {
      "user_name": "Иван",
      "scenarios_count": 15
    }
  }
}
```

### 3. Отправка SMS
```json
{
  "method": "POST",
  "url": "http://localhost:8086/sms/send", 
  "headers": {
    "Authorization": "Bearer {{jwt_token}}",
    "Content-Type": "application/json"
  },
  "body": {
    "phone": "+77001234567",
    "message": "Ваш код подтверждения: 123456"
  }
}
```

### 4. История уведомлений
```json
{
  "method": "GET",
  "url": "http://localhost:8086/notifications/123",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}"
  }
}
```

## 🗄️ Планируемая схема БД

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'push', 'email', 'sms'
    title VARCHAR(255),
    message TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'sent', 'failed'
    sent_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE push_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    device_token VARCHAR(500) NOT NULL,
    platform VARCHAR(20) NOT NULL, -- 'ios', 'android'
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE email_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    html_content TEXT NOT NULL,
    text_content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔧 Планируемая конфигурация

```env
# FCM для Push уведомлений
FCM_SERVER_KEY=your_fcm_server_key

# SMTP для Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=noreply@lunary.com
SMTP_PASS=your_app_password

# SMS провайдер
SMS_API_KEY=your_sms_api_key
SMS_SENDER=LUNARY

# Redis для очередей
REDIS_URL=redis://redis:6379

PORT=8086
```

## 🔮 Планируемые возможности
- [ ] Firebase Cloud Messaging для push-уведомлений
- [ ] SMTP интеграция для email рассылок
- [ ] SMS через внешний API (например, SMS.RU)
- [ ] Очереди уведомлений через Redis
- [ ] Шаблоны для различных типов уведомлений
- [ ] Аналитика доставляемости уведомлений
- [ ] Настройки пользователей (что получать, а что нет)

## 📱 Типы уведомлений

### Push-уведомления
- Новые сценарии по интересам
- Напоминания о предстоящих событиях
- Обновления в избранных сценариях

### Email
- Приветственное письмо
- Еженедельная сводка новых сценариев
- Уведомления о важных обновлениях

### SMS
- Коды подтверждения при регистрации
- Критически важные уведомления
- Восстановление пароля 