@echo off
echo ======= Проверка окружения React Native =======
echo.

echo Проверка Node.js:
node --version
echo.

echo Проверка npm:
npm --version
echo.

echo Проверка Java:
java -version
echo.

echo Проверка Android SDK (если установлен):
if exist "%ANDROID_HOME%" (
    echo Android SDK найден: %ANDROID_HOME%
    echo Проверка ADB:
    "%ANDROID_HOME%\platform-tools\adb.exe" version 2>nul
    if errorlevel 1 (
        echo ADB не найден в %ANDROID_HOME%\platform-tools\
    )
    echo.
    echo Проверка эмуляторов:
    "%ANDROID_HOME%\emulator\emulator.exe" -list-avds 2>nul
    if errorlevel 1 (
        echo Emulator не найден или нет созданных AVD
    )
) else (
    echo Android SDK не найден. Установите Android Studio и настройте ANDROID_HOME
)
echo.

echo Проверка React Native CLI:
npx react-native --version
echo.

echo ======= Рекомендации =======
echo 1. Убедитесь что установлен Android Studio
echo 2. Настройте переменную окружения ANDROID_HOME
echo 3. Добавьте в PATH: %%ANDROID_HOME%%\platform-tools и %%ANDROID_HOME%%\emulator
echo 4. Создайте Android Virtual Device (AVD) через Android Studio
echo 5. Запустите эмулятор перед npm run android
echo.
pause
