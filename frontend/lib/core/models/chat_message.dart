import 'package:lunary_app/core/models/scenario.dart';

enum MessageType {
  user,
  ai,
  scenario,
}

class ChatMessage {
  final String id;
  final String content;
  final MessageType type;
  final DateTime timestamp;
  final Scenario? scenario; // Для сообщений со сценариями

  const ChatMessage({
    required this.id,
    required this.content,
    required this.type,
    required this.timestamp,
    this.scenario,
  });

  factory ChatMessage.user(String content) {
    return ChatMessage(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      content: content,
      type: MessageType.user,
      timestamp: DateTime.now(),
    );
  }

  factory ChatMessage.ai(String content) {
    return ChatMessage(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      content: content,
      type: MessageType.ai,
      timestamp: DateTime.now(),
    );
  }

  factory ChatMessage.scenario(Scenario scenario) {
    return ChatMessage(
      id: DateTime.now().millisecondsSinceEpoch.toString(),
      content: 'Найден подходящий сценарий: ${scenario.title}',
      type: MessageType.scenario,
      timestamp: DateTime.now(),
      scenario: scenario,
    );
  }
}
