import React, {useState} from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Switch,
  Alert,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import Icon from 'react-native-vector-icons/MaterialIcons';

import {theme} from '@constants/theme';
import type {MainTabScreenProps} from '@types/navigation';

export const SettingsScreen: React.FC<MainTabScreenProps<'Settings'>> = () => {
  const [notificationsEnabled, setNotificationsEnabled] = useState(true);
  const [darkModeEnabled, setDarkModeEnabled] = useState(false);
  const [autoSyncEnabled, setAutoSyncEnabled] = useState(true);
  const [analyticsEnabled, setAnalyticsEnabled] = useState(false);

  const settingSections = [
    {
      title: 'Уведомления',
      items: [
        {
          type: 'switch' as const,
          title: 'Включить уведомления',
          description: 'Получать push-уведомления о новых сценариях',
          value: notificationsEnabled,
          onValueChange: setNotificationsEnabled,
          icon: 'notifications',
        },
        {
          type: 'button' as const,
          title: 'Время уведомлений',
          description: 'Настроить время получения уведомлений',
          icon: 'schedule',
          onPress: () => {
            Alert.alert('Функция', 'Настройки времени будут добавлены в следующих версиях');
          },
        },
      ],
    },
    {
      title: 'Внешний вид',
      items: [
        {
          type: 'switch' as const,
          title: 'Темная тема',
          description: 'Использовать темное оформление',
          value: darkModeEnabled,
          onValueChange: (value: boolean) => {
            setDarkModeEnabled(value);
            Alert.alert('Функция', 'Темная тема будет добавлена в следующих версиях');
          },
          icon: 'dark-mode',
        },
        {
          type: 'button' as const,
          title: 'Язык',
          description: 'Русский',
          icon: 'language',
          onPress: () => {
            Alert.alert('Функция', 'Выбор языка будет добавлен в следующих версиях');
          },
        },
      ],
    },
    {
      title: 'Синхронизация',
      items: [
        {
          type: 'switch' as const,
          title: 'Автосинхронизация',
          description: 'Автоматически синхронизировать данные',
          value: autoSyncEnabled,
          onValueChange: setAutoSyncEnabled,
          icon: 'sync',
        },
        {
          type: 'button' as const,
          title: 'Принудительная синхронизация',
          description: 'Синхронизировать данные сейчас',
          icon: 'cloud-sync',
          onPress: () => {
            Alert.alert('Синхронизация', 'Данные успешно синхронизированы');
          },
        },
      ],
    },
    {
      title: 'Конфиденциальность',
      items: [
        {
          type: 'switch' as const,
          title: 'Аналитика',
          description: 'Помочь улучшить приложение, отправляя анонимные данные',
          value: analyticsEnabled,
          onValueChange: setAnalyticsEnabled,
          icon: 'analytics',
        },
        {
          type: 'button' as const,
          title: 'Политика конфиденциальности',
          description: 'Прочитать полную политику конфиденциальности',
          icon: 'privacy-tip',
          onPress: () => {
            Alert.alert('Функция', 'Политика конфиденциальности будет добавлена в следующих версиях');
          },
        },
        {
          type: 'button' as const,
          title: 'Очистить кэш',
          description: 'Очистить временные файлы и кэш',
          icon: 'cleaning-services',
          onPress: () => {
            Alert.alert(
              'Очистка кэша',
              'Вы уверены, что хотите очистить кэш приложения?',
              [
                {text: 'Отмена', style: 'cancel'},
                {
                  text: 'Очистить',
                  onPress: () => Alert.alert('Успешно', 'Кэш очищен'),
                },
              ],
            );
          },
        },
      ],
    },
    {
      title: 'О приложении',
      items: [
        {
          type: 'button' as const,
          title: 'Версия приложения',
          description: '1.0.0',
          icon: 'info',
          onPress: () => {
            Alert.alert(
              'О приложении',
              'LUNARY v1.0.0\n\nПриложение для работы со сценариями и объявлениями.\n\n© 2024 LUNARY',
            );
          },
        },
        {
          type: 'button' as const,
          title: 'Обратная связь',
          description: 'Отправить отзыв или предложение',
          icon: 'feedback',
          onPress: () => {
            Alert.alert('Обратная связь', 'feedback@lunary.app\n\nМы ценим ваше мнение!');
          },
        },
        {
          type: 'button' as const,
          title: 'Оценить приложение',
          description: 'Оставить отзыв в App Store / Google Play',
          icon: 'star-rate',
          onPress: () => {
            Alert.alert('Функция', 'Ссылки на магазины приложений будут добавлены в следующих версиях');
          },
        },
      ],
    },
  ];

  const renderSettingItem = (item: any, index: number) => {
    const isLast = index === settingSections.length - 1;

    if (item.type === 'switch') {
      return (
        <View key={index} style={[styles.settingItem, isLast && styles.lastItem]}>
          <View style={styles.settingIcon}>
            <Icon name={item.icon} size={20} color={theme.colors.primary} />
          </View>
          <View style={styles.settingContent}>
            <Text style={styles.settingTitle}>{item.title}</Text>
            <Text style={styles.settingDescription}>{item.description}</Text>
          </View>
          <Switch
            value={item.value}
            onValueChange={item.onValueChange}
            trackColor={{
              false: theme.colors.gray300,
              true: theme.colors.primaryLight,
            }}
            thumbColor={item.value ? theme.colors.primary : theme.colors.gray400}
          />
        </View>
      );
    }

    return (
      <TouchableOpacity
        key={index}
        style={[styles.settingItem, isLast && styles.lastItem]}
        onPress={item.onPress}>
        <View style={styles.settingIcon}>
          <Icon name={item.icon} size={20} color={theme.colors.primary} />
        </View>
        <View style={styles.settingContent}>
          <Text style={styles.settingTitle}>{item.title}</Text>
          <Text style={styles.settingDescription}>{item.description}</Text>
        </View>
        <Icon name="chevron-right" size={20} color={theme.colors.gray400} />
      </TouchableOpacity>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <Text style={styles.headerTitle}>Настройки</Text>
          <Text style={styles.headerSubtitle}>
            Настройте приложение под себя
          </Text>
        </View>

        {settingSections.map((section, sectionIndex) => (
          <View key={sectionIndex} style={styles.section}>
            <Text style={styles.sectionTitle}>{section.title}</Text>
            <View style={styles.sectionContent}>
              {section.items.map((item, itemIndex) => 
                renderSettingItem(item, itemIndex)
              )}
            </View>
          </View>
        ))}

        <View style={styles.footer}>
          <Text style={styles.footerText}>
            Если у вас есть вопросы или предложения, свяжитесь с нами по адресу support@lunary.app
          </Text>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: theme.colors.backgroundSecondary,
  },
  scrollView: {
    flex: 1,
  },
  header: {
    backgroundColor: theme.colors.white,
    padding: theme.spacing.lg,
    borderBottomWidth: 1,
    borderBottomColor: theme.colors.gray200,
  },
  headerTitle: {
    fontSize: theme.fontSize.xxl,
    fontWeight: theme.fontWeight.bold,
    color: theme.colors.text,
    marginBottom: theme.spacing.xs,
  },
  headerSubtitle: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textSecondary,
  },
  section: {
    marginTop: theme.spacing.lg,
    paddingHorizontal: theme.spacing.lg,
  },
  sectionTitle: {
    fontSize: theme.fontSize.lg,
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.text,
    marginBottom: theme.spacing.md,
  },
  sectionContent: {
    backgroundColor: theme.colors.white,
    borderRadius: theme.borderRadius.lg,
    ...theme.shadows.sm,
  },
  settingItem: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: theme.colors.gray100,
  },
  lastItem: {
    borderBottomWidth: 0,
  },
  settingIcon: {
    marginRight: theme.spacing.md,
  },
  settingContent: {
    flex: 1,
  },
  settingTitle: {
    fontSize: theme.fontSize.md,
    fontWeight: theme.fontWeight.medium,
    color: theme.colors.text,
    marginBottom: theme.spacing.xs,
  },
  settingDescription: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textSecondary,
    lineHeight: 18,
  },
  footer: {
    padding: theme.spacing.xl,
    marginTop: theme.spacing.lg,
  },
  footerText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textMuted,
    textAlign: 'center',
    lineHeight: 20,
  },
});
