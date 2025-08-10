# Установка и настройка React Native приложения

## 🚀 Быстрый старт

### 1. Установка зависимостей

```bash
cd frontend
npm install
```

### 2. iOS (только macOS)

```bash
cd ios
pod install
cd ..
```

### 3. Запуск Metro bundler

```bash
npm start
```

### 4. Запуск на устройстве

В новом терминале:

```bash
# Android
npm run android

# iOS (только macOS) 
npm run ios
```

## 🔧 Настройка среды разработки

### Android

1. Установите Android Studio
2. Настройте Android SDK (API 31+)
3. Создайте виртуальное устройство или подключите физическое
4. Убедитесь что ANDROID_HOME настроен

### iOS (только macOS)

1. Установите Xcode из App Store
2. Установите Command Line Tools: `xcode-select --install`
3. Откройте Xcode и примите лицензии

## 📱 Тестирование на устройстве

### Android

1. Включите режим разработчика на устройстве
2. Включите USB отладку
3. Подключите устройство по USB
4. Запустите `npm run android`

### iOS

1. Подключите iPhone/iPad по USB
2. Доверьтесь компьютеру на устройстве
3. Откройте ios/LunaryMobile.xcworkspace в Xcode
4. Выберите ваше устройство и нажмите Run

## 🐛 Решение проблем

### Metro bundler не запускается

```bash
npx react-native start --reset-cache
```

### Android build ошибки

```bash
cd android
./gradlew clean
cd ..
npm run android
```

### iOS pod install ошибки

```bash
cd ios
rm -rf Pods Podfile.lock
pod install
cd ..
```

### TypeScript ошибки

```bash
npm run lint
npx tsc --noEmit
```

## 📚 Полезные команды

```bash
# Проверка TypeScript
npm run lint

# Форматирование кода
npm run format

# Тесты
npm test

# Очистка кэша
npm start -- --reset-cache

# Сборка release (Android)
cd android && ./gradlew assembleRelease
```
