# Инструкция по тестированию приложения Lunary

## ✅ Исправленные проблемы

1. **Конфликт версий React** - понижен с 19.1.0 до 18.2.0
2. **Java heap memory error** - создан gradle.properties с оптимизированными настройками
3. **Babel и TypeScript конфигурация** - настроены алиасы и пути

## 🚀 Быстрый старт

### Вариант 1: Проверка окружения (рекомендуется сначала)
Запустите `check_environment.bat` для диагностики системы

### Вариант 2: Expo (самый простой)
Запустите `setup_expo_alternative.bat` для создания Expo проекта

### Вариант 3: React Native CLI (требует настройки Android)
Следуйте инструкциям ниже

## Шаги для тестирования:

### 1. Установка зависимостей
Выполните следующие команды в командной строке:

```bash
cd frontend
del package-lock.json
rmdir /s /q node_modules
npm install --legacy-peer-deps
```

Или запустите файл `install_deps.bat` который я создал.

### 2. Запуск Metro Bundler
```bash
cd frontend
npm start
```

Или запустите файл `start_metro.bat`.

### 3. Запуск на Android
В новом терминале:
```bash
cd frontend
npm run android
```

Или запустите файл `run_android.bat`.

### 4. Запуск на iOS (если у вас Mac)
```bash
cd frontend
npm run ios
```

## Исправленные проблемы:

✅ **Babel конфигурация** - добавлен module-resolver для алиасов путей
✅ **TypeScript конфигурация** - настроены baseUrl и paths
✅ **Версии React** - понижены с 19.1.0 до 18.2.0 для совместимости
✅ **Импорты** - исправлен импорт в AppNavigator.tsx

## Структура приложения:

- **Аутентификация**: LoginScreen, RegisterScreen
- **Основные экраны**: HomeScreen, ScenariosScreen, ProfileScreen, SettingsScreen
- **Навигация**: Stack и Tab навигация настроены
- **Состояние**: Redux + Redux Persist для auth и scenario

## Решение проблем Android:

### ❌ Проблема: Java Heap Memory Error
**Решение:** Я создал файл `frontend/android/gradle.properties` с оптимизированными настройками памяти:
- Уменьшил heap с 2048m до 1536m
- Добавил параметры для стабильной работы

### ❌ Проблема: ADB не найден
**Решение:** Установите Android SDK и добавьте в PATH:
1. Установите Android Studio
2. Откройте SDK Manager и установите:
   - Android SDK Platform-Tools
   - Android SDK Build-Tools
3. Добавьте в PATH:
   - `%ANDROID_HOME%\platform-tools` (для adb)
   - `%ANDROID_HOME%\tools` (для emulator)

### ❌ Проблема: Нет эмуляторов
**Решение:** Создайте Android эмулятор:
1. Откройте Android Studio → Tools → AVD Manager
2. Create Virtual Device
3. Выберите устройство (например, Pixel 4)
4. Выберите API Level (например, API 30)
5. Запустите эмулятор

### Альтернативы для тестирования:

#### 1. Подключите физическое устройство:
- Включите Developer Options и USB Debugging
- Подключите через USB
- Разрешите USB debugging

#### 2. Используйте Expo Go (если переключитесь на Expo):
```bash
npm install -g @expo/cli
npx create-expo-app --template blank-typescript
```

### Возможные дополнительные команды:

### Очистка кэша Metro:
```bash
cd frontend
npx react-native start --reset-cache
```

### Очистка Android build:
```bash
cd frontend/android
./gradlew clean
```

### Проверка окружения React Native:
```bash
npx react-native doctor
```

### Проверка зависимостей:
```bash
cd frontend
npm ls
```

### Запуск с дополнительной информацией:
```bash
cd frontend
npx react-native run-android --verbose
```

## Если остаются ошибки TypeScript:

1. Перезапустите TypeScript server в VS Code: `Ctrl+Shift+P` → "TypeScript: Restart TS Server"
2. Убедитесь, что все зависимости установлены корректно
3. Проверьте, что файлы tsconfig.json и babel.config.js содержат правильные настройки

## Тестирование компонентов:

Приложение включает:
- Redux store с auth и scenario slices
- React Navigation с условной навигацией (auth/main)
- Toast уведомления
- Кастомная тема с цветами и стилями
- Базовые сервисы для API, аутентификации и сценариев
