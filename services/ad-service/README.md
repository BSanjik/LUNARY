# Ad Service - LUNARY (В разработке)

## 📝 Описание
Сервис рекламы и монетизации. Управляет показом рекламных объявлений в приложении и партнерскими предложениями для сценариев.

## ⚠️ Статус
**В разработке** - сервис еще не реализован, только заглушка

## 🚀 Планируемые эндпоинты

### GET /ads/banner
Получение баннерной рекламы

### GET /ads/scenario/{scenario_id}
Получение рекламы для конкретного сценария

### POST /ads/click
Регистрация клика по рекламе

### GET /ads/stats
Статистика показов и кликов

## 📋 Планируемая Postman Collection

### 1. Получение баннерной рекламы
```json
{
  "method": "GET",
  "url": "http://localhost:8084/ads/banner",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}"
  },
  "params": [
    {
      "key": "category",
      "value": "family"
    },
    {
      "key": "user_location", 
      "value": "almaty"
    }
  ]
}
```

**Ожидаемый ответ:**
```json
{
  "id": "ad_12345",
  "type": "banner",
  "title": "Лучшие рестораны Алматы",
  "image_url": "https://example.com/banner.jpg",
  "click_url": "https://partner.com/restaurants",
  "partner": "FoodDelivery"
}
```

### 2. Реклама для сценария
```json
{
  "method": "GET",
  "url": "http://localhost:8084/ads/scenario/123",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}"
  }
}
```

**Ожидаемый ответ:**
```json
{
  "id": "ad_67890",
  "type": "scenario_partner",
  "title": "Скидка 20% на доставку еды",
  "description": "Идеально для семейного ужина!",
  "promo_code": "FAMILY20",
  "partner": "Glovo",
  "category": "food_delivery"
}
```

### 3. Регистрация клика
```json
{
  "method": "POST",
  "url": "http://localhost:8084/ads/click",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}",
    "Content-Type": "application/json"
  },
  "body": {
    "ad_id": "ad_12345",
    "user_id": 123,
    "scenario_id": 456
  }
}
```

## 🗄️ Планируемая схема БД

```sql
CREATE TABLE advertisements (
    id SERIAL PRIMARY KEY,
    partner_name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    click_url VARCHAR(500) NOT NULL,
    category VARCHAR(100),
    target_location VARCHAR(100),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ad_impressions (
    id SERIAL PRIMARY KEY,
    ad_id INTEGER REFERENCES advertisements(id),
    user_id INTEGER,
    scenario_id INTEGER,
    shown_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ad_clicks (
    id SERIAL PRIMARY KEY,
    ad_id INTEGER REFERENCES advertisements(id),
    user_id INTEGER,
    scenario_id INTEGER,
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔧 Планируемая конфигурация

```env
# Партнерские интеграции
GLOVO_API_KEY=your_glovo_api_key
CHOCOTRAVEL_API_KEY=your_travel_api_key

# Аналитика
GOOGLE_ANALYTICS_ID=your_ga_id

# Настройки рекламы
MAX_ADS_PER_USER_PER_DAY=10
AD_ROTATION_INTERVAL=300

PORT=8084
```

## 🤝 Планируемые партнеры

### Доставка еды
- **Glovo** - доставка ресторанов
- **Wolt** - премиум доставка
- **Yandex.Еда** - крупнейшая платформа

### Развлечения
- **ivi** - фильмы для домашних вечеров
- **Kinopark** - билеты в кино
- **Escape rooms** - квесты для компаний

### Путешествия  
- **Chocotravel** - туры и отели
- **Booking.com** - бронирование жилья
- **Aviata** - авиабилеты

## 🎯 Таргетинг рекламы

### По категориям сценариев
- `family` → Семейные рестораны, детские товары
- `corporate` → Кейтеринг, конференц-залы
- `friends` → Развлечения, активный отдых

### По геолокации
- Алматы → Местные рестораны и сервисы
- Астана → Региональные предложения
- Другие города → Общенациональные бренды

## 🔮 Планируемые возможности
- [ ] Интеграция с партнерскими API
- [ ] Персонализированный показ рекламы
- [ ] A/B тестирование рекламных креативов
- [ ] Аналитика эффективности кампаний
- [ ] Система промокодов для партнеров
- [ ] Автоматическое планирование показов 