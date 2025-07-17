import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'app/lunary_app.dart';
import 'core/api/api_client.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Инициализируем API клиент
  ApiClient().initialize();

  runApp(
    const ProviderScope(
      child: LunaryApp(),
    ),
  );
}
