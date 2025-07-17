import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../constants/api_constants.dart';

class ApiClient {
  static final ApiClient _instance = ApiClient._internal();
  factory ApiClient() => _instance;
  ApiClient._internal();

  late final Dio _dio;
  late final Dio _publicDio; // Для запросов без авторизации

  Dio get dio => _dio;
  Dio get publicDio => _publicDio; // Getter для публичных запросов

  void initialize() {
    // Обычный Dio с авторизацией
    _dio = Dio(BaseOptions(
      baseUrl: ApiConstants.baseUrl,
      connectTimeout: ApiConstants.connectTimeout,
      receiveTimeout: ApiConstants.receiveTimeout,
      headers: ApiConstants.headers,
    ));

    // Публичный Dio без авторизации
    _publicDio = Dio(BaseOptions(
      baseUrl: ApiConstants.baseUrl,
      connectTimeout: ApiConstants.connectTimeout,
      receiveTimeout: ApiConstants.receiveTimeout,
      headers: ApiConstants.headers,
    ));

    // Interceptor для автоматического добавления JWT токена (только для _dio)
    _dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) async {
        // Добавляем JWT токен в заголовки если есть
        final token = await _getStoredToken();
        if (token != null) {
          options.headers['Authorization'] = 'Bearer $token';
          print('🔑 Отправляю токен: ${token.substring(0, 20)}...');
        } else {
          print('⚠️ Токен не найден для запроса: ${options.uri}');
        }
        handler.next(options);
      },
      onError: (error, handler) async {
        // Если 401 - токен истек, можно перенаправить на логин
        if (error.response?.statusCode == 401) {
          await _clearStoredToken();
          // TODO: Navigate to login screen
        }
        handler.next(error);
      },
    ));

    // Interceptor для логирования (для обоих)
    _dio.interceptors.add(LogInterceptor(
      requestBody: true,
      responseBody: true,
    ));

    _publicDio.interceptors.add(LogInterceptor(
      requestBody: true,
      responseBody: true,
    ));
  }

  // Получение сохраненного токена
  Future<String?> _getStoredToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('jwt_token');
  }

  // Сохранение токена
  Future<void> saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('jwt_token', token);
  }

  // Удаление токена
  Future<void> _clearStoredToken() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('jwt_token');
  }

  // Проверка авторизации
  Future<bool> isLoggedIn() async {
    final token = await _getStoredToken();
    return token != null;
  }

  // Выход из аккаунта
  Future<void> logout() async {
    await _clearStoredToken();
  }
}
