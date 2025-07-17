import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/providers/ai_chat_provider.dart';
import '../../../core/models/chat_message.dart';
import '../../../shared/responsive/adaptive_widgets/adaptive_scaffold.dart';
import '../widgets/chat_message_widget.dart';
import '../widgets/chat_input_widget.dart';
import '../widgets/scenario_card_widget.dart';

class AiChatScreen extends ConsumerWidget {
  const AiChatScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final chatState = ref.watch(aiChatProvider);
    final scrollController = ScrollController();

    // Автоматическая прокрутка к последнему сообщению
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (scrollController.hasClients) {
        scrollController.animateTo(
          scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });

    return AdaptiveScaffold(
      useScrollView: false,
      appBar: AppBar(
        centerTitle: true,
        title: const Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.psychology, color: Colors.blue),
            SizedBox(width: 8),
            Text('AI Помощник'),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () {
              ref.read(aiChatProvider.notifier).clearChat();
            },
            tooltip: 'Очистить чат',
          ),
        ],
      ),
      body: Column(
        children: [
          // Список сообщений
          Expanded(
            child: chatState.messages.isEmpty
                ? const Center(
                    child: CircularProgressIndicator(),
                  )
                : ListView.builder(
                    controller: scrollController,
                    padding: const EdgeInsets.all(16),
                    itemCount: chatState.messages.length,
                    itemBuilder: (context, index) {
                      final message = chatState.messages[index];

                      if (message.type == MessageType.scenario &&
                          message.scenario != null) {
                        return Padding(
                          padding: const EdgeInsets.only(bottom: 16),
                          child: ScenarioCardWidget(
                            scenario: message.scenario!,
                            onFeedback: (rating, comment) {
                              // TODO: Implement feedback
                            },
                          ),
                        );
                      } else {
                        return ChatMessageWidget(
                          message: message,
                          key: ValueKey(message.id),
                        );
                      }
                    },
                  ),
          ),

          // Индикатор загрузки
          if (chatState.isLoading)
            Container(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Container(
                    width: 40,
                    height: 40,
                    decoration: const BoxDecoration(
                      color: Colors.blue,
                      shape: BoxShape.circle,
                    ),
                    child: const Icon(
                      Icons.psychology,
                      color: Colors.white,
                      size: 20,
                    ),
                  ),
                  const SizedBox(width: 12),
                  const Text('AI думает...'),
                  const SizedBox(width: 8),
                  const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  ),
                ],
              ),
            ),

          // Поле ввода
          ChatInputWidget(
            onSend: (message) {
              ref.read(aiChatProvider.notifier).sendMessage(message);
            },
            enabled: !chatState.isLoading,
          ),
        ],
      ),
    );
  }
}
