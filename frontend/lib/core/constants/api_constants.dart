import 'dart:io';

class ApiConstants {
  // Base URL для API Gateway (Android эмулятор использует 10.0.2.2 вместо localhost)
  static String get baseUrl {
    if (Platform.isAndroid) {
      return 'http://10.0.2.2:8080';
    } else {
      return 'http://localhost:8080';
    }
  }

  // Auth endpoints
  static const String authLogin = '/auth/login';
  static const String authRegister = '/auth/registration';

  // Scenario endpoints
  static const String scenarioSearch = '/scenario/search';
  static const String scenarioFeedback = '/scenario/feedback';

  // Headers
  static const Map<String, String> headers = {
    'Content-Type': 'application/json',
  };

  // Timeouts
  static const Duration connectTimeout = Duration(seconds: 10);
  static const Duration receiveTimeout = Duration(seconds: 15);
}
