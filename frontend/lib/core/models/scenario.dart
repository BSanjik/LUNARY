class Step {
  final int id;
  final int scenarioId;
  final String text;
  final String place;
  final String time;
  final double price;

  const Step({
    required this.id,
    required this.scenarioId,
    required this.text,
    required this.place,
    required this.time,
    required this.price,
  });

  factory Step.fromJson(Map<String, dynamic> json) {
    return Step(
      id: json['id'],
      scenarioId: json['scenario_id'],
      text: json['text'],
      place: json['place'] ?? '',
      time: json['time'] ?? '',
      price: (json['price'] ?? 0).toDouble(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'scenario_id': scenarioId,
      'text': text,
      'place': place,
      'time': time,
      'price': price,
    };
  }
}

class Scenario {
  final int id;
  final String title;
  final String description;
  final String category;
  final List<Step> steps;
  final DateTime? createdAt;

  const Scenario({
    required this.id,
    required this.title,
    required this.description,
    required this.category,
    required this.steps,
    this.createdAt,
  });

  factory Scenario.fromJson(Map<String, dynamic> json) {
    return Scenario(
      id: json['id'],
      title: json['title'] ?? '',
      description: json['description'] ?? '',
      category: json['category'] ?? '',
      steps: (json['steps'] as List<dynamic>?)
              ?.map((stepJson) => Step.fromJson(stepJson))
              .toList() ??
          [],
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'description': description,
      'category': category,
      'steps': steps.map((step) => step.toJson()).toList(),
      'created_at': createdAt?.toIso8601String(),
    };
  }

  // Вспомогательные геттеры для совместимости с UI
  String get duration => steps.isNotEmpty
      ? '${steps.first.time} - ${steps.last.time}'
      : 'Не указано';

  String get participants => 'Для всех';

  String get budget =>
      steps
          .fold<double>(0, (sum, step) => sum + step.price)
          .toStringAsFixed(0) +
      ' тенге';

  List<String> get checklist => steps.map((step) => step.text).toList();
}

class ScenarioSearchRequest {
  final String query;

  const ScenarioSearchRequest({required this.query});

  Map<String, dynamic> toJson() {
    return {'query': query};
  }
}
