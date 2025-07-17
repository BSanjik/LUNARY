import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../api/scenario_api.dart';
import '../models/chat_message.dart';
import '../models/scenario.dart';

// Состояние AI чата
class AiChatState {
  final List<ChatMessage> messages;
  final bool isLoading;
  final String? error;

  const AiChatState({
    this.messages = const [],
    this.isLoading = false,
    this.error,
  });

  AiChatState copyWith({
    List<ChatMessage>? messages,
    bool? isLoading,
    String? error,
  }) {
    return AiChatState(
      messages: messages ?? this.messages,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

// Notifier для управления AI чатом
class AiChatNotifier extends StateNotifier<AiChatState> {
  final ScenarioApi _scenarioApi;

  AiChatNotifier(this._scenarioApi) : super(const AiChatState()) {
    _addWelcomeMessage();
  }

  // Добавляем приветственное сообщение
  void _addWelcomeMessage() {
    final welcomeMessage = ChatMessage.ai('Привет! Я AI помощник LUNARY. \n\n'
        'Опишите мероприятие, которое хотите организовать, и я найду для вас подходящий сценарий!\n\n'
        'Например: "семейный ужин на 6 человек" или "корпоратив для коллег"');

    state = state.copyWith(
      messages: [welcomeMessage],
    );
  }

  // Отправка сообщения пользователя и поиск сценария
  Future<void> sendMessage(String message) async {
    if (message.trim().isEmpty) return;

    // Добавляем сообщение пользователя
    final userMessage = ChatMessage.user(message);
    final updatedMessages = [...state.messages, userMessage];

    state = state.copyWith(
      messages: updatedMessages,
      isLoading: true,
      error: null,
    );

    try {
      // Ищем сценарий
      final scenario = await _scenarioApi.searchScenario(message);

      // Добавляем сценарий как сообщение
      final scenarioMessage = ChatMessage.scenario(scenario);
      final finalMessages = [...state.messages, scenarioMessage];

      state = state.copyWith(
        messages: finalMessages,
        isLoading: false,
      );
    } catch (e) {
      // Добавляем сообщение об ошибке
      final errorMessage = ChatMessage.ai(
          'К сожалению, не удалось найти подходящий сценарий. \n\n'
          'Попробуйте переформулировать запрос или использовать другие ключевые слова.\n\n'
          'Например: "романтический ужин", "детский праздник", "встреча друзей"');

      final finalMessages = [...state.messages, errorMessage];

      state = state.copyWith(
        messages: finalMessages,
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  // Очистка чата
  void clearChat() {
    state = const AiChatState();
    _addWelcomeMessage();
  }

  // Очистка ошибки
  void clearError() {
    state = state.copyWith(error: null);
  }
}

// Провайдеры
final scenarioApiProvider = Provider<ScenarioApi>((ref) => ScenarioApi());

final aiChatProvider =
    StateNotifierProvider<AiChatNotifier, AiChatState>((ref) {
  final scenarioApi = ref.read(scenarioApiProvider);
  return AiChatNotifier(scenarioApi);
});
