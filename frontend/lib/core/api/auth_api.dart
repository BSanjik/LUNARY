import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../constants/api_constants.dart';
import '../models/user.dart';
import '../models/auth_response.dart';
import 'api_client.dart';

class AuthApi {
  final ApiClient _apiClient = ApiClient();

  // Очистка номера телефона (убираем пробелы и форматирование)
  String _cleanPhoneNumber(String phone) {
    return phone.replaceAll(' ', '').trim();
  }

  // Регистрация пользователя
  Future<AuthResponse> register({
    required String phone,
    required String password,
    String? email,
  }) async {
    try {
      final cleanPhone = _cleanPhoneNumber(phone);

      final request = AuthRequest(
        phone: cleanPhone,
        password: password,
        email: email,
      );

      if (kDebugMode) {
        print(
            '🚀 Отправка регистрации на: ${ApiConstants.baseUrl}${ApiConstants.authRegister}');
        print('📱 Телефон: $cleanPhone');
        print('📧 Email: $email');
        print('📦 Данные запроса: ${request.toJson()}');
      }

      final response = await _apiClient.dio.post(
        ApiConstants.authRegister,
        data: request.toJson(),
      );

      if (kDebugMode) {
        print('✅ Ответ сервера: ${response.statusCode}');
        print('📄 Данные ответа: ${response.data}');
      }

      if (response.statusCode == 200 || response.statusCode == 201) {
        // Проверяем, что токен не пустой
        final responseData = response.data;
        if (responseData['token'] == null ||
            responseData['token'].toString().isEmpty) {
          throw AuthException(
              'Сервер вернул пустой токен. Проверьте настройки backend.');
        }

        final authResponse = AuthResponse.fromJson(response.data);

        // Сохраняем токен
        await _apiClient.saveToken(authResponse.token);

        if (kDebugMode) {
          print('🔑 Токен сохранен');
        }

        return authResponse;
      } else {
        throw DioException(
          requestOptions: response.requestOptions,
          response: response,
          message: 'Ошибка регистрации',
        );
      }
    } on DioException catch (e) {
      if (kDebugMode) {
        print('❌ Ошибка DioException:');
        print('   Код: ${e.response?.statusCode}');
        print('   Сообщение: ${e.message}');
        print('   Данные ответа: ${e.response?.data}');
        print('   URL: ${e.requestOptions.uri}');
      }

      // Обработка ошибок от сервера
      String errorMessage = 'Ошибка регистрации';

      if (e.response?.statusCode == 409) {
        errorMessage = 'Пользователь с таким номером уже существует';
      } else if (e.response?.statusCode == 400) {
        errorMessage = 'Некорректные данные';
      } else if (e.response?.data != null &&
          e.response!.data['message'] != null) {
        errorMessage = e.response!.data['message'];
      } else if (e.type == DioExceptionType.connectionTimeout) {
        errorMessage = 'Превышено время ожидания подключения';
      } else if (e.type == DioExceptionType.receiveTimeout) {
        errorMessage = 'Превышено время ожидания ответа';
      } else if (e.type == DioExceptionType.connectionError) {
        errorMessage = 'Ошибка подключения к серверу';
      }

      throw AuthException(errorMessage);
    } catch (e) {
      if (kDebugMode) {
        print('❌ Неизвестная ошибка: $e');
      }
      throw AuthException('Ошибка сети');
    }
  }

  // Вход пользователя
  Future<AuthResponse> login({
    required String phone,
    required String password,
  }) async {
    try {
      final cleanPhone = _cleanPhoneNumber(phone);

      final request = AuthRequest(
        phone: cleanPhone,
        password: password,
      );

      if (kDebugMode) {
        print(
            '🚀 Отправка входа на: ${ApiConstants.baseUrl}${ApiConstants.authLogin}');
        print('📱 Телефон: $cleanPhone');
        print('📦 Данные запроса: ${request.toJson()}');
      }

      final response = await _apiClient.dio.post(
        ApiConstants.authLogin,
        data: request.toJson(),
      );

      if (kDebugMode) {
        print('✅ Ответ сервера: ${response.statusCode}');
        print('📄 Данные ответа: ${response.data}');
      }

      if (response.statusCode == 200 || response.statusCode == 201) {
        // Проверяем, что токен не пустой
        final responseData = response.data;
        if (responseData['token'] == null ||
            responseData['token'].toString().isEmpty) {
          throw AuthException(
              'Сервер вернул пустой токен. Проверьте настройки backend.');
        }

        final authResponse = AuthResponse.fromJson(response.data);

        // Сохраняем токен
        await _apiClient.saveToken(authResponse.token);

        if (kDebugMode) {
          print('🔑 Токен сохранен');
        }

        return authResponse;
      } else {
        throw DioException(
          requestOptions: response.requestOptions,
          response: response,
          message: 'Ошибка входа',
        );
      }
    } on DioException catch (e) {
      if (kDebugMode) {
        print('❌ Ошибка DioException:');
        print('   Код: ${e.response?.statusCode}');
        print('   Сообщение: ${e.message}');
        print('   Данные ответа: ${e.response?.data}');
        print('   URL: ${e.requestOptions.uri}');
      }

      String errorMessage = 'Ошибка входа';

      if (e.response?.statusCode == 401) {
        errorMessage = 'Неверный номер телефона или пароль';
      } else if (e.response?.statusCode == 400) {
        errorMessage = 'Некорректные данные';
      } else if (e.response?.data != null &&
          e.response!.data['message'] != null) {
        errorMessage = e.response!.data['message'];
      } else if (e.type == DioExceptionType.connectionTimeout) {
        errorMessage = 'Превышено время ожидания подключения';
      } else if (e.type == DioExceptionType.receiveTimeout) {
        errorMessage = 'Превышено время ожидания ответа';
      } else if (e.type == DioExceptionType.connectionError) {
        errorMessage = 'Ошибка подключения к серверу';
      }

      throw AuthException(errorMessage);
    } catch (e) {
      if (kDebugMode) {
        print('❌ Неизвестная ошибка: $e');
      }
      throw AuthException('Ошибка сети');
    }
  }

  // Выход
  Future<void> logout() async {
    await _apiClient.logout();
  }

  // Проверка авторизации
  Future<bool> isLoggedIn() async {
    return await _apiClient.isLoggedIn();
  }
}

// Кастомное исключение для аутентификации
class AuthException implements Exception {
  final String message;

  const AuthException(this.message);

  @override
  String toString() => 'AuthException: $message';
}
