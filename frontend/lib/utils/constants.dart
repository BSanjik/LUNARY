import 'package:flutter/material.dart';

// API Configuration
class ApiConstants {
  // Base URL для API Gateway
  static const String baseUrl = 'http://localhost:8080';
  
  // Auth endpoints
  static const String registerEndpoint = '/auth/register';
  static const String loginEndpoint = '/auth/login';
  static const String validateEndpoint = '/auth/validate';
  
  // Scenario endpoints
  static const String scenarioSearchEndpoint = '/scenarios/search';
  static const String scenarioDetailsEndpoint = '/scenarios';
  static const String scenarioFeedbackEndpoint = '/scenarios';
  
  // Headers
  static const Map<String, String> defaultHeaders = {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  };
  
  // Request timeouts
  static const Duration requestTimeout = Duration(seconds: 30);
  static const Duration connectTimeout = Duration(seconds: 15);
}

// App Colors
class AppColors {
  // Primary colors (inspired by moon and events)
  static const Color primary = Color(0xFF6366F1); // Indigo
  static const Color primaryLight = Color(0xFF818CF8);
  static const Color primaryDark = Color(0xFF4F46E5);
  
  // Secondary colors
  static const Color secondary = Color(0xFFF59E0B); // Amber
  static const Color secondaryLight = Color(0xFFFBBF24);
  static const Color secondaryDark = Color(0xFFD97706);
  
  // Neutral colors
  static const Color background = Color(0xFFFAFAFA);
  static const Color surface = Color(0xFFFFFFFF);
  static const Color surfaceDark = Color(0xFF1F2937);
  static const Color backgroundDark = Color(0xFF111827);
  
  // Text colors
  static const Color textPrimary = Color(0xFF1F2937);
  static const Color textSecondary = Color(0xFF6B7280);
  static const Color textLight = Color(0xFF9CA3AF);
  static const Color textOnDark = Color(0xFFF9FAFB);
  
  // Status colors
  static const Color success = Color(0xFF10B981);
  static const Color warning = Color(0xFFF59E0B);
  static const Color error = Color(0xFFEF4444);
  static const Color info = Color(0xFF3B82F6);
  
  // Card and component colors
  static const Color cardBackground = Color(0xFFFFFFFF);
  static const Color cardBorder = Color(0xFFE5E7EB);
  static const Color divider = Color(0xFFE5E7EB);
}

// App Dimensions
class AppDimensions {
  // Padding and margins
  static const double paddingXS = 4.0;
  static const double paddingS = 8.0;
  static const double paddingM = 16.0;
  static const double paddingL = 24.0;
  static const double paddingXL = 32.0;
  
  // Border radius
  static const double radiusS = 8.0;
  static const double radiusM = 12.0;
  static const double radiusL = 16.0;
  static const double radiusXL = 24.0;
  
  // Icon sizes
  static const double iconS = 16.0;
  static const double iconM = 24.0;
  static const double iconL = 32.0;
  static const double iconXL = 48.0;
  
  // Button heights
  static const double buttonHeightS = 36.0;
  static const double buttonHeightM = 48.0;
  static const double buttonHeightL = 56.0;
}

// App Strings
class AppStrings {
  // App name
  static const String appName = 'LUNARY';
  static const String appTagline = 'Планируйте идеальные мероприятия с ИИ';
  
  // Auth
  static const String login = 'Вход';
  static const String register = 'Регистрация';
  static const String logout = 'Выход';
  static const String name = 'Имя';
  static const String phone = 'Телефон';
  static const String password = 'Пароль';
  static const String confirmPassword = 'Подтвердите пароль';
  static const String forgotPassword = 'Забыли пароль?';
  static const String dontHaveAccount = 'Нет аккаунта?';
  static const String alreadyHaveAccount = 'Уже есть аккаунт?';
  
  // Navigation
  static const String home = 'Главная';
  static const String search = 'Поиск';
  static const String profile = 'Профиль';
  static const String favorites = 'Избранное';
  static const String settings = 'Настройки';
  
  // Scenarios
  static const String scenarios = 'Сценарии';
  static const String searchScenarios = 'Поиск сценариев...';
  static const String scenarioDetails = 'Детали сценария';
  static const String noScenariosFound = 'Сценарии не найдены';
  static const String rating = 'Рейтинг';
  static const String feedback = 'Отзыв';
  static const String sendFeedback = 'Оставить отзыв';
  
  // Common
  static const String loading = 'Загрузка...';
  static const String error = 'Ошибка';
  static const String success = 'Успешно';
  static const String retry = 'Повторить';
  static const String cancel = 'Отмена';
  static const String save = 'Сохранить';
  static const String delete = 'Удалить';
  static const String edit = 'Редактировать';
  
  // Validation messages
  static const String requiredField = 'Обязательное поле';
  static const String invalidPhone = 'Неверный формат телефона';
  static const String invalidPassword = 'Пароль должен содержать минимум 6 символов';
  static const String passwordsDoNotMatch = 'Пароли не совпадают';
}

// Regular expressions for validation
class AppRegex {
  static final RegExp phone = RegExp(r'^\d{10,15}$');
  static final RegExp email = RegExp(r'^[^@]+@[^@]+\.[^@]+$');
  static final RegExp password = RegExp(r'^.{6,}$');
}