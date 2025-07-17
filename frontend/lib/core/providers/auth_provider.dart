import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../api/auth_api.dart';
import '../models/user.dart';
import '../models/auth_response.dart';
import '../api/user_api.dart';

// Состояние аутентификации
class AuthState {
  final bool isLoading;
  final bool isAuthenticated;
  final User? user;
  final String? error;

  const AuthState({
    this.isLoading = false,
    this.isAuthenticated = false,
    this.user,
    this.error,
  });

  AuthState copyWith({
    bool? isLoading,
    bool? isAuthenticated,
    User? user,
    String? error,
  }) {
    return AuthState(
      isLoading: isLoading ?? this.isLoading,
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      user: user ?? this.user,
      error: error,
    );
  }
}

// Notifier для управления состоянием аутентификации
class AuthNotifier extends StateNotifier<AuthState> {
  final AuthApi _authApi;

  AuthNotifier(this._authApi) : super(const AuthState()) {
    _checkAuthStatus();
  }

  // Проверка статуса аутентификации при запуске
  Future<void> _checkAuthStatus() async {
    state = state.copyWith(isLoading: true);

    try {
      final isLoggedIn = await _authApi.isLoggedIn();
      state = state.copyWith(
        isLoading: false,
        isAuthenticated: isLoggedIn,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isAuthenticated: false,
        error: 'Ошибка проверки аутентификации',
      );
    }
  }

  // Регистрация
  Future<void> register({
    required String phone,
    required String password,
    String? email,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final authResponse = await _authApi.register(
        phone: phone,
        password: password,
        email: email,
      );

      state = state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: authResponse.user,
        error: null,
      );
    } on AuthException catch (e) {
      state = state.copyWith(
        isLoading: false,
        isAuthenticated: false,
        error: e.message,
      );
      rethrow;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isAuthenticated: false,
        error: 'Неизвестная ошибка',
      );
      rethrow;
    }
  }

  // Вход
  Future<void> login({
    required String phone,
    required String password,
  }) async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final authResponse = await _authApi.login(
        phone: phone,
        password: password,
      );

      state = state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: authResponse.user,
        error: null,
      );
    } on AuthException catch (e) {
      state = state.copyWith(
        isLoading: false,
        isAuthenticated: false,
        error: e.message,
      );
      rethrow;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isAuthenticated: false,
        error: 'Неизвестная ошибка',
      );
      rethrow;
    }
  }

  // Выход
  Future<void> logout() async {
    state = state.copyWith(isLoading: true);

    try {
      await _authApi.logout();
      state = const AuthState(isAuthenticated: false);
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: 'Ошибка выхода',
      );
    }
  }

  // Очистка ошибки
  void clearError() {
    state = state.copyWith(error: null);
  }
}

// Провайдеры
final authApiProvider = Provider<AuthApi>((ref) => AuthApi());

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final authApi = ref.read(authApiProvider);
  return AuthNotifier(authApi);
});

// Удобные провайдеры для отдельных состояний
final isAuthenticatedProvider = Provider<bool>((ref) {
  return ref.watch(authProvider).isAuthenticated;
});

final isLoadingProvider = Provider<bool>((ref) {
  return ref.watch(authProvider).isLoading;
});

final currentUserProvider = Provider<User?>((ref) {
  return ref.watch(authProvider).user;
});
