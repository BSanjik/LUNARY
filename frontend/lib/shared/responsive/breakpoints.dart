class Breakpoints {
  // Точки перелома экрана
  static const double mobile = 600;
  static const double tablet = 900;
  static const double desktop = 1200;

  // Проверка типа устройства
  static bool isMobile(double width) => width < mobile;
  static bool isTablet(double width) => width >= mobile && width < desktop;
  static bool isDesktop(double width) => width >= desktop;

  // Получение типа устройства
  static DeviceType getDeviceType(double width) {
    if (isMobile(width)) return DeviceType.mobile;
    if (isTablet(width)) return DeviceType.tablet;
    return DeviceType.desktop;
  }
}

enum DeviceType {
  mobile,
  tablet,
  desktop,
}
