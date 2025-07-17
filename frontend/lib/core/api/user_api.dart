import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../models/user.dart';
import 'api_client.dart';

class UserApi {
  final ApiClient _apiClient = ApiClient();

  // Получение текущего пользователя (декодируем из JWT токена)
  Future<User?> getCurrentUser() async {
    try {
      // Пока backend не готов, создадим пользователя из localStorage телефона
      // Позже заменим на реальный API call

      if (kDebugMode) {
        print('🔍 Попытка получить данные пользователя...');
      }

      // Временно создаем пользователя с данными из последней регистрации
      // В реальности здесь должен быть API call к /user/me или декодинг JWT

      return const User(
        id: 1,
        phone: "+77001234567", // временно захардкодим
        email: null,
      );
    } catch (e) {
      if (kDebugMode) {
        print('❌ Ошибка получения пользователя: $e');
      }
      return null;
    }
  }
}
