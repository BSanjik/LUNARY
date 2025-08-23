@echo off
echo ======= Установка Expo как альтернативы =======
echo.
echo Это создаст новый проект Expo для тестирования,
echo если React Native CLI вызывает проблемы
echo.
pause

echo Установка Expo CLI...
npm install -g @expo/cli
echo.

echo Создание нового Expo проекта...
cd ..
npx create-expo-app LunaryExpo --template blank-typescript
echo.

echo Переход в проект и запуск...
cd LunaryExpo
npm start
echo.

echo ======= Инструкции =======
echo 1. Установите Expo Go на телефон из App Store/Google Play
echo 2. Отсканируйте QR код который появится
echo 3. Приложение откроется на вашем телефоне
echo.
pause
