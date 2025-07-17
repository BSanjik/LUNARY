import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../constants/api_constants.dart';
import '../models/scenario.dart';
import 'api_client.dart';

class ScenarioApi {
  final ApiClient _apiClient = ApiClient();

  // Поиск сценариев (без авторизации)
  Future<Scenario> searchScenario(String query) async {
    try {
      final request = ScenarioSearchRequest(query: query);

      if (kDebugMode) {
        print('🔍 Поиск сценария (БЕЗ токена): "$query"');
        print(
            '🚀 Отправка на: ${ApiConstants.baseUrl}${ApiConstants.scenarioSearch}');
        print('📦 Данные запроса: ${request.toJson()}');
      }

      // Используем publicDio для запросов без авторизации
      final response = await _apiClient.publicDio.post(
        ApiConstants.scenarioSearch,
        data: request.toJson(),
      );

      if (kDebugMode) {
        print('✅ Ответ сервера: ${response.statusCode}');
        print('📄 Данные ответа: ${response.data}');
      }

      if (response.statusCode == 200) {
        final scenario = Scenario.fromJson(response.data);

        if (kDebugMode) {
          print('🎯 Найден сценарий: ${scenario.title}');
        }

        return scenario;
      } else {
        throw DioException(
          requestOptions: response.requestOptions,
          response: response,
          message: 'Ошибка поиска сценария',
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

      String errorMessage = 'Ошибка поиска сценария';

      if (e.response?.statusCode == 404) {
        errorMessage = 'Сценарий не найден. Попробуйте изменить запрос.';
      } else if (e.response?.statusCode == 400) {
        errorMessage = 'Некорректный запрос';
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

      throw ScenarioException(errorMessage);
    } catch (e) {
      if (kDebugMode) {
        print('❌ Неизвестная ошибка: $e');
      }
      throw ScenarioException('Ошибка сети');
    }
  }

  // Отправка отзыва о сценарии (без авторизации)
  Future<void> submitFeedback({
    required int scenarioId,
    required String userLogin,
    required int rating,
    String? comment,
  }) async {
    try {
      final feedbackData = {
        'scenario_id': scenarioId,
        'user_login': userLogin,
        'rating': rating,
        'comment': comment,
      };

      if (kDebugMode) {
        print('📝 Отправка отзыва для сценария $scenarioId (БЕЗ токена)');
        print('⭐ Рейтинг: $rating');
      }

      // Используем publicDio для запросов без авторизации
      final response = await _apiClient.publicDio.post(
        ApiConstants.scenarioFeedback,
        data: feedbackData,
      );

      if (kDebugMode) {
        print('✅ Отзыв отправлен успешно');
      }
    } catch (e) {
      if (kDebugMode) {
        print('❌ Ошибка отправки отзыва: $e');
      }
      throw ScenarioException('Ошибка отправки отзыва');
    }
  }
}

// Кастомное исключение для сценариев
class ScenarioException implements Exception {
  final String message;

  const ScenarioException(this.message);

  @override
  String toString() => 'ScenarioException: $message';
}

// Класс для запроса поиска сценария
class ScenarioSearchRequest {
  final String query;

  ScenarioSearchRequest({required this.query});

  Map<String, dynamic> toJson() {
    return {
      'query': query,
    };
  }
}
