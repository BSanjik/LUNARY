# LUNARY Mobile App (React Native)

Мобильное приложение LUNARY, построенное на React Native с TypeScript.

## 🚀 Технологии

- **React Native 0.72** - кроссплатформенная разработка
- **TypeScript** - типизированный JavaScript
- **Redux Toolkit** - управление состоянием
- **React Navigation 6** - навигация между экранами
- **React Hook Form** - управление формами
- **Axios** - HTTP клиент
- **React Native Paper** - UI компоненты
- **Vector Icons** - иконки

## 📁 Структура проекта

```
frontend/
├── src/
│   ├── components/       # Переиспользуемые компоненты
│   ├── screens/         # Экраны приложения
│   │   ├── auth/        # Экраны аутентификации
│   │   └── main/        # Основные экраны
│   ├── navigation/      # Навигация
│   ├── services/        # API сервисы
│   ├── store/          # Redux store и slices
│   ├── types/          # TypeScript типы
│   ├── constants/      # Константы (тема, API)
│   └── utils/          # Утилиты
├── assets/             # Изображения и иконки
├── android/           # Android конфигурация
├── ios/              # iOS конфигурация
└── package.json      # Зависимости и scripts
```

## 🛠️ Установка и запуск

### Требования
- Node.js 16+
- React Native CLI
- Android Studio (для Android)
- Xcode (для iOS, только macOS)

### Установка зависимостей
```bash
cd frontend
npm install

# Для iOS (только macOS)
cd ios && pod install && cd ..
```

### Запуск Metro bundler
```bash
npm start
```

### Запуск на устройстве
```bash
# Android
npm run android

# iOS (только macOS)
npm run ios
```

## 🏗️ Архитектура

### State Management (Redux Toolkit)
- **authSlice** - управление аутентификацией
- **scenarioSlice** - управление сценариями
- Persist для сохранения auth состояния

### Navigation (React Navigation)
- **AuthNavigator** - экраны входа/регистрации
- **MainNavigator** - основные экраны с tab navigation
- **AppNavigator** - корневая навигация

### API Integration
- **apiService** - базовый HTTP клиент с interceptors
- **authService** - методы аутентификации
- **scenarioService** - методы работы со сценариями

## 📱 Экраны

### Аутентификация
- **LoginScreen** - вход в систему
- **RegisterScreen** - регистрация пользователя

### Основные экраны
- **HomeScreen** - главная с статистикой
- **ScenariosScreen** - список сценариев
- **ProfileScreen** - профиль пользователя
- **SettingsScreen** - настройки приложения

## 🎨 Дизайн система

### Цветовая палитра
```typescript
colors: {
  primary: '#6366F1',      // Основной цвет
  secondary: '#EC4899',    // Дополнительный
  success: '#10B981',      // Успех
  error: '#EF4444',        // Ошибка
  warning: '#F59E0B',      // Предупреждение
  // ... остальные цвета
}
```

### Типографика
- Размеры: xs(12) - xxxl(32)
- Веса: light(300) - bold(700)

### Отступы и скругления
- Spacing: xs(4) - xxl(48)
- Border radius: sm(4) - full(9999)

## 🔒 Безопасность

### Аутентификация
- JWT токены с автоматическим обновлением
- Безопасное хранение в AsyncStorage
- Автоматический logout при 401

### Валидация
- React Hook Form с validation rules
- Валидация телефонов в формате E.164
- Проверка силы пароля

## 📡 API Integration

### Base URL
- **Development**: `http://localhost:8081`
- **Production**: `https://api.lunary.app`

### Endpoints
```typescript
AUTH: {
  REGISTRATION: '/api/v1/registration',
  LOGIN: '/api/v1/login',
},
SCENARIOS: {
  LIST: '/api/v1/scenarios',
  DETAIL: '/api/v1/scenarios/:id',
}
```

### Error Handling
- Структурированные ошибки с кодами
- Пользовательские сообщения на русском
- Автоматический retry для network errors

## 🧪 Тестирование

```bash
# Запуск тестов
npm test

# Линтинг
npm run lint

# Type checking
npx tsc --noEmit
```

## 📦 Сборка

### Development
```bash
# Android Debug
npm run build:android

# iOS Debug (только macOS)
npm run build:ios
```

### Production
```bash
# Android Release
cd android && ./gradlew assembleRelease

# iOS Release (только macOS)
cd ios && xcodebuild -workspace LunaryMobile.xcworkspace \
  -scheme LunaryMobile -configuration Release archive
```

## 🔧 Конфигурация

### Metro Config
- Path mapping для удобных imports
- React Native Reanimated plugin

### Babel Config
- Module resolver для @aliases
- TypeScript preset

### ESLint
- React Native конфигурация
- TypeScript правила

## 📱 Features

### ✅ Реализовано
- 🔐 **Аутентификация** - вход/регистрация
- 🏠 **Главная** - dashboard с статистикой
- 📋 **Сценарии** - список с фильтрацией
- 👤 **Профиль** - информация о пользователе
- ⚙️ **Настройки** - конфигурация приложения
- 🎨 **Дизайн система** - единое оформление
- 📱 **Адаптивность** - поддержка разных экранов

### 🚧 В разработке
- 🔄 **Offline mode** - работа без интернета
- 🌙 **Темная тема** - dark mode
- 🔔 **Push уведомления** - Firebase integration
- 📷 **Загрузка файлов** - media service integration
- 🌐 **Локализация** - поддержка языков

## 🤝 Разработка

### Code Style
- **TypeScript** для всех компонентов
- **Functional components** с hooks
- **Arrow functions** для консистентности
- **Destructuring** для props

### Naming Conventions
- **PascalCase** для компонентов
- **camelCase** для функций и переменных
- **SCREAMING_SNAKE_CASE** для констант
- **kebab-case** для файлов assets

### Git Workflow
```bash
# Создание feature branch
git checkout -b feature/new-screen

# Commit с описанием
git commit -m "feat: add new screen component"

# Push и создание PR
git push origin feature/new-screen
```

## 📞 Поддержка

- **Issues**: [GitHub Issues](https://github.com/BSanjik/LUNARY/issues)
- **Email**: support@lunary.app
- **Документация**: [docs/](../docs/)

---

**Сделано с ❤️ командой LUNARY**
