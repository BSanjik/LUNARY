# Media Service - LUNARY (В разработке)

## 📝 Описание
Сервис для загрузки, хранения и обработки медиафайлов. Обрабатывает фотографии событий, аватары пользователей и другие медиа-ресурсы.

## ⚠️ Статус
**В разработке** - сервис еще не реализован, только заглушка

## 🚀 Планируемые эндпоинты

### POST /upload/image
Загрузка изображения

### GET /image/{image_id}
Получение изображения по ID

### DELETE /image/{image_id}
Удаление изображения

### POST /upload/avatar
Загрузка аватара пользователя

## 📋 Планируемая Postman Collection

### 1. Загрузка изображения
```json
{
  "method": "POST",
  "url": "http://localhost:8085/upload/image",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}"
  },
  "body": {
    "mode": "formdata",
    "formdata": [
      {
        "key": "image",
        "type": "file", 
        "src": "path/to/image.jpg"
      },
      {
        "key": "category",
        "value": "event_photo"
      }
    ]
  }
}
```

### 2. Получение изображения
```json
{
  "method": "GET",
  "url": "http://localhost:8085/image/123456",
  "headers": {
    "Authorization": "Bearer {{jwt_token}}"
  }
}
```

## 🗄️ Планируемая схема БД

```sql
CREATE TABLE media_files (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255),
    file_size BIGINT,
    mime_type VARCHAR(100),
    s3_url VARCHAR(500),
    user_id INTEGER,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔧 Планируемая конфигурация

```env
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=minioaccess
MINIO_SECRET_KEY=miniosecret
BUCKET_NAME=lunary-media
PORT=8085
```

## 🔮 Планируемые возможности
- [ ] Интеграция с MinIO для хранения файлов
- [ ] Автоматическое изменение размера изображений
- [ ] Поддержка различных форматов медиа
- [ ] CDN для быстрой доставки контента 