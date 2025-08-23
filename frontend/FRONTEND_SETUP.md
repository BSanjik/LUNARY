# 📱 LUNARY Frontend Setup

## 🚀 Быстрый старт

### 1. Установка зависимостей
```bash
cd frontend
npm install --legacy-peer-deps
```

### 2. Настройка API endpoints
Проверьте конфигурацию в `src/constants/api.ts`:
- **API Gateway**: `http://localhost:8080`
- **Auth Service**: `http://localhost:8081` 
- **Scenario Service**: `http://localhost:8082`

### 3. Запуск приложения

#### Android:
```bash
npm run android
```

#### iOS:
```bash
npm run ios
```

#### Metro bundler:
```bash
npm start
```

## 🧪 Тестирование API

В приложении есть встроенная кнопка **"Тест API"** на главном экране, которая:
- Проверяет health check всех сервисов
- Тестирует регистрацию и логин
- Тестирует поиск сценариев и отзывы

### Ручное тестирование API:
```javascript
import { ApiTest, checkBackendHealth } from './src/utils/apiTest';

// Проверить здоровье сервисов
await checkBackendHealth();

// Полный тест API
await ApiTest.runFullTest();
```

## 📋 Backend API Endpoints

### Auth Service (8081):
- `POST /api/v1/registration` - Регистрация
- `POST /api/v1/login` - Логин
- `GET /health` - Health check

### Scenario Service (8082):
- `POST /search` - Поиск сценариев
- `POST /feedback` - Отправка отзыва
- `GET /health` - Health check

### API Gateway (8080):
- Проксирует все запросы к сервисам
- Требует JWT токен для защищенных endpoints

## 🔧 Настройка разработки

### Network Security (Android)
Если backend запущен на localhost, убедитесь что в Android разрешен HTTP трафик:

```xml
<!-- android/app/src/main/AndroidManifest.xml -->
<application
  android:usesCleartextTraffic="true"
  ...>
```

### iOS Simulator
По умолчанию localhost доступен из iOS симулятора.

## 🎯 Готовые компоненты

- ✅ **Аутентификация** - экраны логина и регистрации
- ✅ **Навигация** - Tab Navigator с экранами
- ✅ **State Management** - Redux Toolkit
- ✅ **API Service** - Axios с interceptors
- ✅ **Безопасность** - Keychain для токенов
- ✅ **UI Library** - React Native Paper

## 🚨 Troubleshooting

### Конфликт зависимостей:
```bash
npm install --legacy-peer-deps
```

### Metro cache:
```bash
npx react-native start --reset-cache
```

### Android build issues:
```bash
cd android && ./gradlew clean && cd ..
```

## 📱 Поддерживаемые платформы

- ✅ Android 6.0+ (API 23+)
- ✅ iOS 11.0+
- 🔄 React Native 0.80.2

