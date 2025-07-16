class Scenario {
  final int id;
  final String title;
  final String description;
  final String? detailedDescription;
  final String category;
  final String duration;
  final String participants;
  final String budget;
  final List<String> tags;
  final double rating;
  final List<String>? materials;
  final List<ScenarioStep>? steps;
  final DateTime createdAt;
  final DateTime? updatedAt;

  const Scenario({
    required this.id,
    required this.title,
    required this.description,
    this.detailedDescription,
    required this.category,
    required this.duration,
    required this.participants,
    required this.budget,
    required this.tags,
    required this.rating,
    this.materials,
    this.steps,
    required this.createdAt,
    this.updatedAt,
  });

  factory Scenario.fromJson(Map<String, dynamic> json) {
    return Scenario(
      id: json['id'] as int,
      title: json['title'] as String,
      description: json['description'] as String,
      detailedDescription: json['detailed_description'] as String?,
      category: json['category'] as String,
      duration: json['duration'] as String,
      participants: json['participants'] as String,
      budget: json['budget'] as String,
      tags: List<String>.from(json['tags'] as List? ?? []),
      rating: (json['rating'] as num?)?.toDouble() ?? 0.0,
      materials: json['materials'] != null 
          ? List<String>.from(json['materials'] as List)
          : null,
      steps: json['steps'] != null
          ? (json['steps'] as List)
              .map((step) => ScenarioStep.fromJson(step as Map<String, dynamic>))
              .toList()
          : null,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: json['updated_at'] != null
          ? DateTime.parse(json['updated_at'] as String)
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'detailed_description': detailedDescription,
      'category': category,
      'duration': duration,
      'participants': participants,
      'budget': budget,
      'tags': tags,
      'rating': rating,
      'materials': materials,
      'steps': steps?.map((step) => step.toJson()).toList(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
    };
  }

  Scenario copyWith({
    int? id,
    String? title,
    String? description,
    String? detailedDescription,
    String? category,
    String? duration,
    String? participants,
    String? budget,
    List<String>? tags,
    double? rating,
    List<String>? materials,
    List<ScenarioStep>? steps,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Scenario(
      id: id ?? this.id,
      title: title ?? this.title,
      description: description ?? this.description,
      detailedDescription: detailedDescription ?? this.detailedDescription,
      category: category ?? this.category,
      duration: duration ?? this.duration,
      participants: participants ?? this.participants,
      budget: budget ?? this.budget,
      tags: tags ?? this.tags,
      rating: rating ?? this.rating,
      materials: materials ?? this.materials,
      steps: steps ?? this.steps,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  String toString() {
    return 'Scenario(id: $id, title: $title, category: $category, rating: $rating)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Scenario && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}

class ScenarioStep {
  final int step;
  final String title;
  final String description;
  final String? duration;

  const ScenarioStep({
    required this.step,
    required this.title,
    required this.description,
    this.duration,
  });

  factory ScenarioStep.fromJson(Map<String, dynamic> json) {
    return ScenarioStep(
      step: json['step'] as int,
      title: json['title'] as String,
      description: json['description'] as String,
      duration: json['duration'] as String?,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'step': step,
      'title': title,
      'description': description,
      'duration': duration,
    };
  }

  ScenarioStep copyWith({
    int? step,
    String? title,
    String? description,
    String? duration,
  }) {
    return ScenarioStep(
      step: step ?? this.step,
      title: title ?? this.title,
      description: description ?? this.description,
      duration: duration ?? this.duration,
    );
  }

  @override
  String toString() {
    return 'ScenarioStep(step: $step, title: $title)';
  }
}

class ScenarioSearchResponse {
  final List<Scenario> scenarios;
  final int total;
  final int limit;
  final int offset;

  const ScenarioSearchResponse({
    required this.scenarios,
    required this.total,
    required this.limit,
    required this.offset,
  });

  factory ScenarioSearchResponse.fromJson(Map<String, dynamic> json) {
    return ScenarioSearchResponse(
      scenarios: (json['scenarios'] as List)
          .map((scenario) => Scenario.fromJson(scenario as Map<String, dynamic>))
          .toList(),
      total: json['total'] as int,
      limit: json['limit'] as int,
      offset: json['offset'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'scenarios': scenarios.map((scenario) => scenario.toJson()).toList(),
      'total': total,
      'limit': limit,
      'offset': offset,
    };
  }

  bool get hasMore => offset + scenarios.length < total;
}

class Feedback {
  final int id;
  final int scenarioId;
  final int userId;
  final int rating;
  final String? comment;
  final DateTime createdAt;

  const Feedback({
    required this.id,
    required this.scenarioId,
    required this.userId,
    required this.rating,
    this.comment,
    required this.createdAt,
  });

  factory Feedback.fromJson(Map<String, dynamic> json) {
    return Feedback(
      id: json['id'] as int,
      scenarioId: json['scenario_id'] as int,
      userId: json['user_id'] as int,
      rating: json['rating'] as int,
      comment: json['comment'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'scenario_id': scenarioId,
      'user_id': userId,
      'rating': rating,
      'comment': comment,
      'created_at': createdAt.toIso8601String(),
    };
  }

  Feedback copyWith({
    int? id,
    int? scenarioId,
    int? userId,
    int? rating,
    String? comment,
    DateTime? createdAt,
  }) {
    return Feedback(
      id: id ?? this.id,
      scenarioId: scenarioId ?? this.scenarioId,
      userId: userId ?? this.userId,
      rating: rating ?? this.rating,
      comment: comment ?? this.comment,
      createdAt: createdAt ?? this.createdAt,
    );
  }

  @override
  String toString() {
    return 'Feedback(id: $id, scenarioId: $scenarioId, rating: $rating)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Feedback && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
} 