import 'package:flutter/material.dart';
import '../responsive.dart';
import '../responsive_wrapper.dart';

class AdaptiveForm extends StatelessWidget {
  final Widget child;
  final String? title;
  final String? subtitle;
  final Widget? header;

  const AdaptiveForm({
    super.key,
    required this.child,
    this.title,
    this.subtitle,
    this.header,
  });

  @override
  Widget build(BuildContext context) {
    return ResponsiveWrapper(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Отступ сверху
          SizedBox(
            height: Responsive.value(
              context,
              mobile: 40.0,
              tablet: 60.0,
              desktop: 80.0,
            ),
          ),

          // Заголовок
          if (header != null)
            header!
          else if (title != null) ...[
            // Логотип
            Icon(
              Icons.psychology,
              size: Responsive.value(
                context,
                mobile: 60.0,
                tablet: 80.0,
                desktop: 100.0,
              ),
              color: Colors.blue,
            ),
            const SizedBox(height: 16),

            Text(
              'LUNARY',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: Responsive.fontSize(
                  context,
                  mobile: 28.0,
                  tablet: 32.0,
                  desktop: 36.0,
                ),
                fontWeight: FontWeight.bold,
                color: Colors.blue,
              ),
            ),

            if (title != null) ...[
              const SizedBox(height: 8),
              Text(
                title!,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: Responsive.fontSize(
                    context,
                    mobile: 16.0,
                    tablet: 18.0,
                    desktop: 20.0,
                  ),
                  color: Colors.grey,
                ),
              ),
            ],
          ],

          SizedBox(
            height: Responsive.value(
              context,
              mobile: 32.0,
              tablet: 48.0,
              desktop: 64.0,
            ),
          ),

          // Основной контент
          child,

          // Отступ снизу
          SizedBox(
            height: Responsive.value(
              context,
              mobile: 40.0,
              tablet: 60.0,
              desktop: 80.0,
            ),
          ),
        ],
      ),
    );
  }
}
