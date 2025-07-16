import 'dart:convert';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'package:logger/logger.dart';

import '../models/user.dart';
import '../utils/constants.dart';

class AuthApiService {
  static final Logger _logger = Logger();
  static const Duration _timeout = ApiConstants.requestTimeout;

  /// Register new user
  Future<AuthResponse?> register({
    required String name,
    required String phone,
    required String password,
  }) async {
    try {
      _logger.i('Registering user: $name, $phone');
      
      final url = Uri.parse('${ApiConstants.baseUrl}${ApiConstants.registerEndpoint}');
      
      final response = await http.post(
        url,
        headers: ApiConstants.defaultHeaders,
        body: jsonEncode({
          'name': name,
          'phone': phone,
          'password': password,
        }),
      ).timeout(_timeout);

      _logger.d('Register response status: ${response.statusCode}');
      _logger.d('Register response body: ${response.body}');

      if (response.statusCode == 201) {
        final data = jsonDecode(response.body) as Map<String, dynamic>;
        return AuthResponse.fromJson(data);
      } else if (response.statusCode == 409) {
        throw 'User already exists';
      } else if (response.statusCode == 400) {
        final data = jsonDecode(response.body) as Map<String, dynamic>;
        throw data['error']?['message'] ?? 'Invalid data';
      } else {
        throw 'Server error: ${response.statusCode}';
      }
    } on SocketException {
      _logger.e('Network error during registration');
      throw 'Network connection error';
    } catch (e) {
      _logger.e('Registration error: $e');
      rethrow;
    }
  }

  /// Login user
  Future<AuthResponse?> login({
    required String phone,
    required String password,
  }) async {
    try {
      _logger.i('Logging in user: $phone');
      
      final url = Uri.parse('${ApiConstants.baseUrl}${ApiConstants.loginEndpoint}');
      
      final response = await http.post(
        url,
        headers: ApiConstants.defaultHeaders,
        body: jsonEncode({
          'phone': phone,
          'password': password,
        }),
      ).timeout(_timeout);

      _logger.d('Login response status: ${response.statusCode}');
      _logger.d('Login response body: ${response.body}');

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body) as Map<String, dynamic>;
        return AuthResponse.fromJson(data);
      } else if (response.statusCode == 401) {
        throw 'Invalid credentials';
      } else if (response.statusCode == 400) {
        final data = jsonDecode(response.body) as Map<String, dynamic>;
        throw data['error']?['message'] ?? 'Invalid data';
      } else {
        throw 'Server error: ${response.statusCode}';
      }
    } on SocketException {
      _logger.e('Network error during login');
      throw 'Network connection error';
    } catch (e) {
      _logger.e('Login error: $e');
      rethrow;
    }
  }

  /// Validate JWT token
  Future<User?> validateToken(String token) async {
    try {
      _logger.i('Validating token');
      
      final url = Uri.parse('${ApiConstants.baseUrl}${ApiConstants.validateEndpoint}');
      
      final response = await http.get(
        url,
        headers: {
          ...ApiConstants.defaultHeaders,
          'Authorization': 'Bearer $token',
        },
      ).timeout(_timeout);

      _logger.d('Validate token response status: ${response.statusCode}');

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body) as Map<String, dynamic>;
        if (data['valid'] == true && data['user'] != null) {
          return User.fromJson(data['user'] as Map<String, dynamic>);
        }
      } else if (response.statusCode == 401) {
        _logger.w('Token is invalid');
        return null;
      }
      
      return null;
    } on SocketException {
      _logger.e('Network error during token validation');
      throw 'Network connection error';
    } catch (e) {
      _logger.e('Token validation error: $e');
      return null; // Don't throw on token validation errors
    }
  }

  /// Test connection to auth service
  Future<bool> testConnection() async {
    try {
      _logger.i('Testing connection to auth service');
      
      final url = Uri.parse('${ApiConstants.baseUrl}/health');
      
      final response = await http.get(
        url,
        headers: ApiConstants.defaultHeaders,
      ).timeout(const Duration(seconds: 10));

      _logger.d('Health check response status: ${response.statusCode}');
      
      return response.statusCode == 200;
    } catch (e) {
      _logger.e('Connection test error: $e');
      return false;
    }
  }
}

/// Auth response model for register/login
class AuthResponse {
  final String message;
  final User user;
  final String token;

  const AuthResponse({
    required this.message,
    required this.user,
    required this.token,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    return AuthResponse(
      message: json['message'] as String,
      user: User.fromJson(json['user'] as Map<String, dynamic>),
      token: json['token'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'message': message,
      'user': user.toJson(),
      'token': token,
    };
  }

  @override
  String toString() {
    return 'AuthResponse(message: $message, user: ${user.name}, token: ${token.substring(0, 10)}...)';
  }
}