import 'package:flutter/material.dart';
import 'package:mask_text_input_formatter/mask_text_input_formatter.dart';
import 'package:provider/provider.dart';

// --- Data Models ---

enum AuthType {
  seller,
  buyer,
}

class SellerData extends ChangeNotifier {
  String _sellerName;
  String _companyName;
  String _phoneNumber;
  String _userPassword;

  SellerData({
    String sellerName = '',
    String companyName = '',
    String phoneNumber = '',
    String userPassword = '',
  })  : _sellerName = sellerName,
        _companyName = companyName,
        _phoneNumber = phoneNumber,
        _userPassword = userPassword;

  String get sellerName => _sellerName; 
  String get companyName => _companyName;
  String get phoneNumber => _phoneNumber;
  String get userPassword => _userPassword;

  void updateSellerName(String name) {
    _sellerName = name;
    notifyListeners();
  }

  void updateCompanyName(String name) {
    _companyName = name;
    notifyListeners();
  }

  void updatePhoneNumber(String name) {
    _phoneNumber = name;
    notifyListeners();
  }
  
  void updateUserPassword(String name) {
    _userPassword = name;
    notifyListeners();
  }
}


class Preference {
  final String name;
  bool isSelected;

  Preference({required this.name, this.isSelected = false});

  void toggleSelected() {
    isSelected = !isSelected;
  }
}

class BuyerData extends ChangeNotifier {
  String _buyerName;
  String _buyerPhoneNumber;
  String _buyerEmail;
  String _buyerPassword;
  final List<Preference> _preferences;

  BuyerData({
    String buyerName = '',
    String buyerPhoneNumber = '',
    String buyerEmail = '',
    String buyerPassword = '',
    List<Preference>? initialPreferences,
  })  : _buyerName = buyerName,
        _buyerPhoneNumber = buyerPhoneNumber,
        _buyerEmail = buyerEmail,
        _buyerPassword = buyerPassword,
        _preferences = initialPreferences ??
            [
              Preference(name: 'Электроника'),
              Preference(name: 'Одежда'),
              Preference(name: 'Книги'),
              Preference(name: 'Спорт'),
              Preference(name: 'Дом и сад'),
              Preference(name: 'Авто'),
              Preference(name: 'Детские товары'),
              Preference(name: 'Продукты'),
              Preference(name: 'Косметика'),
              Preference(name: 'Путешествия'),
            ];

  String get buyerName => _buyerName;
  String get buyerPhoneNumber => _buyerPhoneNumber;
  String get buyerEmail => _buyerEmail;
  String get buyerPassword => _buyerPassword;
  List<Preference> get preferences => List.unmodifiable(_preferences);

  void updateBuyerName(String name) {
    _buyerName = name;
    notifyListeners();
  }

  void updateBuyerPhoneNumber(String name) {
    _buyerPhoneNumber = name;
    notifyListeners();
  }

  void updateBuyerEmail(String name) {
    _buyerEmail = name;
    notifyListeners();
  }

  void updateBuyerPassword(String name) {
    _buyerPassword = name;
    notifyListeners();
  }

  void togglePreference(Preference preferenceToToggle) {
    for (final pref in _preferences) {
      if (pref == preferenceToToggle) {
        pref.toggleSelected();
        break;
      }
    }
    notifyListeners();
  }
}

// --- UI Screens ---

class AuthSelectionScreen extends StatelessWidget {
  const AuthSelectionScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Авторизация'),
        centerTitle: true,
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              ElevatedButton(
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute<void>(
                        builder: (context) => const SellerInputScreen()),
                  );
                },
                style: ElevatedButton.styleFrom(
                  minimumSize: const Size.fromHeight(50), // Full width button
                ),
                child: const Text(
                  'Продавец',
                  style: TextStyle(fontSize: 18),
                ),
              ),
              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute<void>(
                        builder: (context) => const BuyerInputScreen()),
                  );
                },
                style: ElevatedButton.styleFrom(
                  minimumSize: const Size.fromHeight(50), // Full width button
                ),
                child: const Text(
                  'Покупатель',
                  style: TextStyle(fontSize: 18),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class SellerInputScreen extends StatefulWidget {
  const SellerInputScreen({super.key});

  @override
  State<SellerInputScreen> createState() => _SellerInputScreenState();
}

class _SellerInputScreenState extends State<SellerInputScreen> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _sellerNameController;
  late final TextEditingController _companyNameController;
  late final TextEditingController _phoneNumber;
  late final TextEditingController _userPassword;

  final phoneMaskFormatter = MaskTextInputFormatter(
    mask: '+7 (###) ###-##-##',
    filter: {"#": RegExp(r'[0-9]')},
    type: MaskAutoCompletionType.lazy,
  );

  @override
  void initState() {
    super.initState();
    _sellerNameController = TextEditingController();
    _companyNameController = TextEditingController();
    _phoneNumber = TextEditingController();
    _userPassword = TextEditingController();
  }

  @override
  void dispose() {
    _sellerNameController.dispose();
    _companyNameController.dispose();
    _phoneNumber.dispose();
    _userPassword.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Данные продавца'),
        centerTitle: true,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: <Widget>[
              TextFormField(
                controller: _sellerNameController,
                decoration: const InputDecoration(
                  labelText: 'Имя продавца',
                  border: OutlineInputBorder(),
                ),
                validator: (String? value) {
                  if (value == null || value.isEmpty) {
                    return 'Пожалуйста, введите имя продавца';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              TextFormField(
                controller: _companyNameController,
                decoration: const InputDecoration(
                  labelText: 'Название компании',
                  border: OutlineInputBorder(),
                ),
                validator: (String? value) {
                  if (value == null || value.isEmpty) {
                    return 'Пожалуйста, введите название компании';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              TextFormField(
                controller: _phoneNumber,
                inputFormatters: [phoneMaskFormatter],
                keyboardType: TextInputType.phone,
                decoration: const InputDecoration(
                  labelText: 'Номер телефона',
                  border: OutlineInputBorder(),
                  hintText: '+7 (___) ___-__-__',
                ),
                validator: (String? value) {
                  if (value == null || value.isEmpty) {
                    return 'Пожалуйста, введите номер телефона';
                  }
                  if (value.length < 18) {
                    return 'Номер введен не полностью';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
              TextFormField(
                controller: _userPassword,
                obscureText: true,
                decoration: const InputDecoration(
                  labelText: 'Придумайте пароль',
                  border: OutlineInputBorder(),
                ),
                validator: (String? value) {
                  if (value == null || value.isEmpty) {
                    return 'Пожалуйста, придумайте пароль';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 30),
              ElevatedButton(
                onPressed: () {
                  if (_formKey.currentState!.validate()) {
                    context
                        .read<SellerData>()
                        .updateSellerName(_sellerNameController.text);
                    context
                        .read<SellerData>()
                        .updateCompanyName(_companyNameController.text);
                    context
                        .read<SellerData>()
                        .updatePhoneNumber(_phoneNumber.text);
                    context
                        .read<SellerData>()
                        .updateUserPassword(_userPassword.text);

                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Данные продавца сохранены!')),
                    );
                    Navigator.pop(context); // Go back to selection screen
                  }
                },
                style: ElevatedButton.styleFrom(
                  minimumSize: const Size.fromHeight(50),
                ),
                child: const Text(
                  'Зарегистрироваться',
                  style: TextStyle(fontSize: 18),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class BuyerInputScreen extends StatefulWidget {
  const BuyerInputScreen({super.key});

  @override
  State<BuyerInputScreen> createState() => _BuyerInputScreenState();
}

class _BuyerInputScreenState extends State<BuyerInputScreen> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _buyerName;
  late final TextEditingController _buyerPhoneNumber;
  late final TextEditingController _buyerEmail;
  late final TextEditingController _buyerPassword;

  final phoneMaskFormatter = MaskTextInputFormatter(
    mask: '+7 (###) ###-##-##',
    filter: {"#": RegExp(r'[0-9]')},
    type: MaskAutoCompletionType.lazy,
  );

  @override
  void initState() {
    super.initState();
    _buyerName = TextEditingController();
    _buyerPhoneNumber = TextEditingController();
    _buyerEmail = TextEditingController();
    _buyerPassword = TextEditingController();
  }

  @override
  void dispose() {
    _buyerName.dispose();
    _buyerPhoneNumber.dispose();
    _buyerEmail.dispose();
    _buyerPassword.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Данные покупателя'),
        centerTitle: true,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: <Widget>[
              TextFormField(
                controller: _buyerName,
                decoration: const InputDecoration(
                  labelText: 'Ваше имя',
                  border: OutlineInputBorder(),
                ),
                validator: (String? value) {
                  if (value == null || value.isEmpty) {
                    return 'Пожалуйста, введите ваше имя';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
TextFormField(
                controller: _buyerPhoneNumber,
                inputFormatters: [phoneMaskFormatter],
                keyboardType: TextInputType.phone,
                decoration: const InputDecoration(
                  labelText: 'Номер телефона',
                  border: OutlineInputBorder(),
                  hintText: '+7 (___) ___-__-__',
                ),
                validator: (String? value) {
                  if (value == null || value.isEmpty) {
                    return 'Пожалуйста, введите номер телефона';
                  }
                  if (value.length < 18) {
                    return 'Номер введен не полностью';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
               TextFormField(
                controller: _buyerEmail,
                decoration: const InputDecoration(
                  labelText: 'Ваша почта',
                  border: OutlineInputBorder(),
                ),
                validator: (String? value) {
                  if (value == null || value.isEmpty) {
                    return 'Пожалуйста, введите вашу почту';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 20),
                TextFormField(
                controller: _buyerPassword,
                obscureText: true,
                decoration: const InputDecoration(
                  labelText: 'Придумайте пароль',
                  border: OutlineInputBorder(),
                ),
                validator: (String? value) {
                  if (value == null || value.isEmpty) {
                    return 'Пожалуйста, придумайте пароль';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 30),
              ElevatedButton(
                onPressed: () {
                  if (_formKey.currentState!.validate()) {
                    context
                        .read<BuyerData>()
                        .updateBuyerName(_buyerName.text);
                    context
                        .read<BuyerData>()
                        .updateBuyerPhoneNumber(_buyerPhoneNumber.text);
                    context
                        .read<BuyerData>()
                        .updateBuyerEmail(_buyerEmail.text);
                    context
                        .read<BuyerData>()
                        .updateBuyerPassword(_buyerPassword.text);

                    Navigator.push(
                      context,
                      MaterialPageRoute<void>(
                          builder: (context) => const BuyerPreferencesScreen()),
                    );
                  }
                },
                style: ElevatedButton.styleFrom(
                  minimumSize: const Size.fromHeight(50),
                ),
                child: const Text(
                  'Регистрация',
                  style: TextStyle(fontSize: 18),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class BuyerPreferencesScreen extends StatelessWidget {
  const BuyerPreferencesScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Выберите предпочтения'),
        centerTitle: true,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: <Widget>[
            Text(
              'Здравствуйте, ${context.watch<BuyerData>().buyerName}! Выберите, что вам интересно:',
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 20),
            Consumer<BuyerData>(
              builder: (context, buyerData, child) {
                return Expanded(
                  child: SingleChildScrollView(
                    child: Wrap(
                      spacing: 8.0, // gap between adjacent chips
                      runSpacing: 8.0, // gap between lines
                      children: buyerData.preferences
                          .map<Widget>((Preference preference) {
                        return ChoiceChip(
                          label: Text(preference.name),
                          selected: preference.isSelected,
                          onSelected: (bool selected) {
                            buyerData.togglePreference(preference);
                          },
                          selectedColor: Theme.of(context).primaryColor.withOpacity(0.2),
                          checkmarkColor: Theme.of(context).primaryColor,
                          labelStyle: TextStyle(
                            color: preference.isSelected
                                ? Theme.of(context).primaryColor
                                : Colors.black87,
                          ),
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(20.0),
                            side: BorderSide(
                              color: preference.isSelected
                                  ? Theme.of(context).primaryColor
                                  : Colors.grey,
                            ),
                          ),
                          backgroundColor: Colors.grey[200],
                        );
                      }).toList(),
                    ),
                  ),
                );
              },
            ),
            const SizedBox(height: 20),
            ElevatedButton(
              onPressed: () {
                // Here you might navigate to a main app screen or show a success message
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Предпочтения сохранены!')),
                );
                Navigator.popUntil(context, (route) => route.isFirst); // Go back to the root
              },
              style: ElevatedButton.styleFrom(
                minimumSize: const Size.fromHeight(50),
              ),
              child: const Text(
                'Готово',
                style: TextStyle(fontSize: 18),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

