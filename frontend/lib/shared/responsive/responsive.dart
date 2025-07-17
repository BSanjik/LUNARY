import 'package:flutter/material.dart';
import 'breakpoints.dart';

class Responsive {
  static double width(BuildContext context) {
    return MediaQuery.of(context).size.width;
  }

  static double height(BuildContext context) {
    return MediaQuery.of(context).size.height;
  }

  static bool isMobile(BuildContext context) {
    return Breakpoints.isMobile(width(context));
  }

  static bool isTablet(BuildContext context) {
    return Breakpoints.isTablet(width(context));
  }

  static bool isDesktop(BuildContext context) {
    return Breakpoints.isDesktop(width(context));
  }

  static DeviceType deviceType(BuildContext context) {
    return Breakpoints.getDeviceType(width(context));
  }

  // Адаптивные значения
  static T value<T>(
    BuildContext context, {
    required T mobile,
    T? tablet,
    T? desktop,
  }) {
    final deviceType = Responsive.deviceType(context);

    switch (deviceType) {
      case DeviceType.mobile:
        return mobile;
      case DeviceType.tablet:
        return tablet ?? mobile;
      case DeviceType.desktop:
        return desktop ?? tablet ?? mobile;
    }
  }

  // Адаптивные отступы
  static EdgeInsets padding(BuildContext context) {
    return EdgeInsets.all(
      value(
        context,
        mobile: 16.0,
        tablet: 24.0,
        desktop: 32.0,
      ),
    );
  }

  // Адаптивные размеры форм
  static double formWidth(BuildContext context) {
    return value(
      context,
      mobile: width(context) * 0.9,
      tablet: 400.0,
      desktop: 450.0,
    );
  }

  // Адаптивный размер шрифта
  static double fontSize(
    BuildContext context, {
    required double mobile,
    double? tablet,
    double? desktop,
  }) {
    return value(
      context,
      mobile: mobile,
      tablet: tablet,
      desktop: desktop,
    );
  }
}
