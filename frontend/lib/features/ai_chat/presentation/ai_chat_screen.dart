import 'package:flutter/material.dart';

class AiChatScreen extends StatelessWidget {
  const AiChatScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Row(
          children: [
            Icon(Icons.psychology, color: Colors.blue),
            SizedBox(width: 8),
            Text('AI Помощник'),
          ],
        ),
      ),
      body: const Center(
        child: Text(
          'AI-чат\n(В разработке)',
          textAlign: TextAlign.center,
          style: TextStyle(fontSize: 18),
        ),
      ),
    );
  }
}
