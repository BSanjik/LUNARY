import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/user.dart';
import '../services/auth_api.dart';

enum AuthStatus {
  initial,
  loading,
  authenticated,
  unauthenticated,
  error,
}

class AuthProvider extends ChangeNotifier {
  AuthStatus _status = AuthStatus.initial;
  User? _user;
  String? _token;
  String? _errorMessage;
  final AuthApiService _authApi = AuthApiService();

  // Getters
  AuthStatus get status => _status;
  User? get user => _user;
  String? get token => _token;
  String? get errorMessage => _errorMessage;
  bool get isAuthenticated => _status == AuthStatus.authenticated && _user != null;
  bool get isLoading => _status == AuthStatus.loading;

  /// Initialize auth state - check for existing token
  Future<void> initialize() async {
    _setStatus(AuthStatus.loading);
    
    try {
      final prefs = await SharedPreferences.getInstance();
      final savedToken = prefs.getString('auth_token');
      
      if (savedToken != null) {
        _token = savedToken;
        
        // Validate token with backend
        final user = await _authApi.validateToken(savedToken);
        if (user != null) {
          _user = user;
          _setStatus(AuthStatus.authenticated);
        } else {
          await _clearSession();
          _setStatus(AuthStatus.unauthenticated);
        }
      } else {
        _setStatus(AuthStatus.unauthenticated);
      }
    } catch (e) {
      debugPrint('Auth initialization error: $e');
      await _clearSession();
      _setStatus(AuthStatus.unauthenticated);
    }
  }

  /// Register new user
  Future<bool> register({
    required String name,
    required String phone,
    required String password,
  }) async {
    _setStatus(AuthStatus.loading);
    _clearError();

    try {
      final response = await _authApi.register(
        name: name,
        phone: phone,
        password: password,
      );

      if (response != null) {
        _user = response.user;
        _token = response.token;
        await _saveSession();
        _setStatus(AuthStatus.authenticated);
        return true;
      } else {
        _setError('Не удалось создать аккаунт');
        _setStatus(AuthStatus.unauthenticated);
        return false;
      }
    } catch (e) {
      _setError(_getErrorMessage(e));
      _setStatus(AuthStatus.unauthenticated);
      return false;
    }
  }

  /// Login user
  Future<bool> login({
    required String phone,
    required String password,
  }) async {
    _setStatus(AuthStatus.loading);
    _clearError();

    try {
      final response = await _authApi.login(
        phone: phone,
        password: password,
      );

      if (response != null) {
        _user = response.user;
        _token = response.token;
        await _saveSession();
        _setStatus(AuthStatus.authenticated);
        return true;
      } else {
        _setError('Неверный телефон или пароль');
        _setStatus(AuthStatus.unauthenticated);
        return false;
      }
    } catch (e) {
      _setError(_getErrorMessage(e));
      _setStatus(AuthStatus.unauthenticated);
      return false;
    }
  }

  /// Logout user
  Future<void> logout() async {
    _setStatus(AuthStatus.loading);
    
    try {
      await _clearSession();
      _setStatus(AuthStatus.unauthenticated);
    } catch (e) {
      debugPrint('Logout error: $e');
      // Even if there's an error, clear local session
      await _clearSession();
      _setStatus(AuthStatus.unauthenticated);
    }
  }

  /// Update user profile
  Future<bool> updateProfile({
    String? name,
    String? email,
    UserPreferences? preferences,
  }) async {
    if (!isAuthenticated || _user == null) return false;

    try {
      final updatedUser = _user!.copyWith(
        name: name ?? _user!.name,
        email: email ?? _user!.email,
        preferences: preferences ?? _user!.preferences,
        updatedAt: DateTime.now(),
      );

      // TODO: Call API to update user profile when User Service is ready
      // For now, just update locally
      _user = updatedUser;
      notifyListeners();
      
      return true;
    } catch (e) {
      _setError(_getErrorMessage(e));
      return false;
    }
  }

  /// Refresh user data
  Future<void> refreshUser() async {
    if (!isAuthenticated || _token == null) return;

    try {
      final user = await _authApi.validateToken(_token!);
      if (user != null) {
        _user = user;
        notifyListeners();
      }
    } catch (e) {
      debugPrint('Refresh user error: $e');
    }
  }

  /// Save session to local storage
  Future<void> _saveSession() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      if (_token != null) {
        await prefs.setString('auth_token', _token!);
      }
      if (_user != null) {
        await prefs.setString('user_data', _user!.toJson().toString());
      }
    } catch (e) {
      debugPrint('Save session error: $e');
    }
  }

  /// Clear session from local storage
  Future<void> _clearSession() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.remove('auth_token');
      await prefs.remove('user_data');
      _user = null;
      _token = null;
    } catch (e) {
      debugPrint('Clear session error: $e');
    }
  }

  /// Set auth status
  void _setStatus(AuthStatus status) {
    _status = status;
    notifyListeners();
  }

  /// Set error message
  void _setError(String message) {
    _errorMessage = message;
    notifyListeners();
  }

  /// Clear error message
  void _clearError() {
    _errorMessage = null;
  }

  /// Get user-friendly error message
  String _getErrorMessage(dynamic error) {
    if (error is String) return error;
    
    final errorStr = error.toString().toLowerCase();
    
    if (errorStr.contains('connection') || errorStr.contains('network')) {
      return 'Ошибка подключения к интернету';
    } else if (errorStr.contains('timeout')) {
      return 'Превышено время ожидания';
    } else if (errorStr.contains('user already exists')) {
      return 'Пользователь с таким номером уже существует';
    } else if (errorStr.contains('invalid credentials')) {
      return 'Неверный номер телефона или пароль';
    } else if (errorStr.contains('invalid phone')) {
      return 'Неверный формат номера телефона';
    } else if (errorStr.contains('password')) {
      return 'Пароль должен содержать минимум 6 символов';
    }
    
    return 'Произошла ошибка. Попробуйте еще раз';
  }

  /// Validate phone number format
  static bool isValidPhone(String phone) {
    final cleanPhone = phone.replaceAll(RegExp(r'[^\d]'), '');
    return cleanPhone.length >= 10 && cleanPhone.length <= 15;
  }

  /// Validate password
  static bool isValidPassword(String password) {
    return password.length >= 6;
  }

  /// Format phone number for display
  static String formatPhone(String phone) {
    final cleanPhone = phone.replaceAll(RegExp(r'[^\d]'), '');
    if (cleanPhone.length == 11 && cleanPhone.startsWith('8')) {
      return '+7 ${cleanPhone.substring(1, 4)} ${cleanPhone.substring(4, 7)}-${cleanPhone.substring(7, 9)}-${cleanPhone.substring(9)}';
    } else if (cleanPhone.length == 11 && cleanPhone.startsWith('7')) {
      return '+${cleanPhone.substring(0, 1)} ${cleanPhone.substring(1, 4)} ${cleanPhone.substring(4, 7)}-${cleanPhone.substring(7, 9)}-${cleanPhone.substring(9)}';
    }
    return phone;
  }

  /// Clean phone number for API
  static String cleanPhone(String phone) {
    return phone.replaceAll(RegExp(r'[^\d]'), '');
  }
}
