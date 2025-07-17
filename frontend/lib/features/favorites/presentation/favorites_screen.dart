import 'package:flutter/material.dart';

class FavoritesScreen extends StatelessWidget {
  // ← Исправляем ProfileScreen на FavoritesScreen
  const FavoritesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Избранное'), // ← Исправляем заголовок
      ),
      body: const Center(
        child: Text(
          'Избранное\n(В разработке)', // ← Исправляем текст
          textAlign: TextAlign.center,
          style: TextStyle(fontSize: 18),
        ),
      ),
    );
  }
}
