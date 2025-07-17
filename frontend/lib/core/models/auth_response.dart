import 'package:lunary_app/core/models/user.dart';

class AuthResponse {
  final String? message;
  final String token;
  final User? user;

  const AuthResponse({
    this.message,
    required this.token,
    this.user,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) {
    print('🔍 Парсинг AuthResponse: $json');
    print('🔍 Есть ли user в ответе: ${json.containsKey('user')}');
    print('🔍 User data: ${json['user']}');

    return AuthResponse(
      message: json['message'], // может быть null
      token: json['token'],
      user: json['user'] != null ? User.fromJson(json['user']) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'message': message,
      'token': token,
      'user': user?.toJson(),
    };
  }
}

class AuthRequest {
  final String phone;
  final String password;
  final String? email;

  const AuthRequest({
    required this.phone,
    required this.password,
    this.email,
  });

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = {
      'phone': phone,
      'password': password,
    };

    if (email != null) {
      data['email'] = email;
    }

    return data;
  }
}
