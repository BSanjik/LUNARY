class User {
  final int id;
  final String name;
  final String phone;
  final String? email;
  final DateTime createdAt;
  final DateTime? updatedAt;
  final UserPreferences? preferences;
  final UserStatistics? statistics;

  const User({
    required this.id,
    required this.name,
    required this.phone,
    this.email,
    required this.createdAt,
    this.updatedAt,
    this.preferences,
    this.statistics,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as int,
      name: json['name'] as String,
      phone: json['phone'] as String,
      email: json['email'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: json['updated_at'] != null 
          ? DateTime.parse(json['updated_at'] as String)
          : null,
      preferences: json['preferences'] != null
          ? UserPreferences.fromJson(json['preferences'] as Map<String, dynamic>)
          : null,
      statistics: json['statistics'] != null
          ? UserStatistics.fromJson(json['statistics'] as Map<String, dynamic>)
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'phone': phone,
      'email': email,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
      'preferences': preferences?.toJson(),
      'statistics': statistics?.toJson(),
    };
  }

  User copyWith({
    int? id,
    String? name,
    String? phone,
    String? email,
    DateTime? createdAt,
    DateTime? updatedAt,
    UserPreferences? preferences,
    UserStatistics? statistics,
  }) {
    return User(
      id: id ?? this.id,
      name: name ?? this.name,
      phone: phone ?? this.phone,
      email: email ?? this.email,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      preferences: preferences ?? this.preferences,
      statistics: statistics ?? this.statistics,
    );
  }

  @override
  String toString() {
    return 'User(id: $id, name: $name, phone: $phone, email: $email)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is User && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}

class UserPreferences {
  final List<String> favoriteCategories;
  final String budgetRange;
  final bool notificationsEnabled;
  final String? language;
  final String? theme;

  const UserPreferences({
    required this.favoriteCategories,
    required this.budgetRange,
    required this.notificationsEnabled,
    this.language,
    this.theme,
  });

  factory UserPreferences.fromJson(Map<String, dynamic> json) {
    return UserPreferences(
      favoriteCategories: List<String>.from(json['favorite_categories'] as List),
      budgetRange: json['budget_range'] as String,
      notificationsEnabled: json['notifications_enabled'] as bool,
      language: json['language'] as String?,
      theme: json['theme'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'favorite_categories': favoriteCategories,
      'budget_range': budgetRange,
      'notifications_enabled': notificationsEnabled,
      'language': language,
      'theme': theme,
    };
  }

  UserPreferences copyWith({
    List<String>? favoriteCategories,
    String? budgetRange,
    bool? notificationsEnabled,
    String? language,
    String? theme,
  }) {
    return UserPreferences(
      favoriteCategories: favoriteCategories ?? this.favoriteCategories,
      budgetRange: budgetRange ?? this.budgetRange,
      notificationsEnabled: notificationsEnabled ?? this.notificationsEnabled,
      language: language ?? this.language,
      theme: theme ?? this.theme,
    );
  }
}

class UserStatistics {
  final int scenariosViewed;
  final int feedbackGiven;
  final int eventsPlanned;
  final int favoriteScenarios;

  const UserStatistics({
    required this.scenariosViewed,
    required this.feedbackGiven,
    required this.eventsPlanned,
    required this.favoriteScenarios,
  });

  factory UserStatistics.fromJson(Map<String, dynamic> json) {
    return UserStatistics(
      scenariosViewed: json['scenarios_viewed'] as int,
      feedbackGiven: json['feedback_given'] as int,
      eventsPlanned: json['events_planned'] as int,
      favoriteScenarios: json['favorite_scenarios'] as int? ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'scenarios_viewed': scenariosViewed,
      'feedback_given': feedbackGiven,
      'events_planned': eventsPlanned,
      'favorite_scenarios': favoriteScenarios,
    };
  }

  UserStatistics copyWith({
    int? scenariosViewed,
    int? feedbackGiven,
    int? eventsPlanned,
    int? favoriteScenarios,
  }) {
    return UserStatistics(
      scenariosViewed: scenariosViewed ?? this.scenariosViewed,
      feedbackGiven: feedbackGiven ?? this.feedbackGiven,
      eventsPlanned: eventsPlanned ?? this.eventsPlanned,
      favoriteScenarios: favoriteScenarios ?? this.favoriteScenarios,
    );
  }
}