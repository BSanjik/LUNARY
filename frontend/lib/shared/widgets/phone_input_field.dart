import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class PhoneInputField extends StatefulWidget {
  final TextEditingController controller;
  final String? Function(String?)? validator;
  final String labelText;
  final String hintText;

  const PhoneInputField({
    super.key,
    required this.controller,
    this.validator,
    this.labelText = 'Номер телефона',
    this.hintText = '+7 700 123 45 67',
  });

  @override
  State<PhoneInputField> createState() => _PhoneInputFieldState();
}

class _PhoneInputFieldState extends State<PhoneInputField> {
  @override
  void initState() {
    super.initState();

    // Устанавливаем +7 если поле пустое
    if (widget.controller.text.isEmpty) {
      widget.controller.text = '+7 ';
      widget.controller.selection = TextSelection.fromPosition(
        TextPosition(offset: widget.controller.text.length),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      controller: widget.controller,
      keyboardType: TextInputType.phone,
      inputFormatters: [
        PhoneNumberFormatter(),
      ],
      decoration: InputDecoration(
        labelText: widget.labelText,
        hintText: widget.hintText,
        prefixIcon: const Icon(Icons.phone),
        border: const OutlineInputBorder(),
      ),
      validator: widget.validator,
    );
  }
}

class PhoneNumberFormatter extends TextInputFormatter {
  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    String newText = newValue.text;

    // Не даем удалить +7
    if (!newText.startsWith('+7')) {
      newText = '+7 ' + newText.replaceAll(RegExp(r'[^\d]'), '');
    }

    // Удаляем все нецифровые символы кроме +7
    String digitsOnly = newText.substring(2).replaceAll(RegExp(r'[^\d]'), '');

    // Ограничиваем до 10 цифр после +7
    if (digitsOnly.length > 10) {
      digitsOnly = digitsOnly.substring(0, 10);
    }

    // Форматируем: +7 XXX XXX XX XX
    String formatted = '+7';

    if (digitsOnly.isNotEmpty) {
      formatted += ' ';

      if (digitsOnly.length >= 1) {
        formatted += digitsOnly.substring(
            0, digitsOnly.length > 3 ? 3 : digitsOnly.length);
      }

      if (digitsOnly.length >= 4) {
        formatted +=
            ' ${digitsOnly.substring(3, digitsOnly.length > 6 ? 6 : digitsOnly.length)}';
      }

      if (digitsOnly.length >= 7) {
        formatted +=
            ' ${digitsOnly.substring(6, digitsOnly.length > 8 ? 8 : digitsOnly.length)}';
      }

      if (digitsOnly.length >= 9) {
        formatted += ' ${digitsOnly.substring(8)}';
      }
    }

    // Определяем новую позицию курсора
    int cursorPosition = formatted.length;

    return TextEditingValue(
      text: formatted,
      selection: TextSelection.collapsed(offset: cursorPosition),
    );
  }
}
