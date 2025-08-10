import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  Alert,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {useDispatch, useSelector} from 'react-redux';
import Icon from 'react-native-vector-icons/MaterialIcons';

import {logoutUser} from '@store/slices/authSlice';
import {theme} from '@constants/theme';
import type {AppDispatch, RootState} from '@store/index';
import type {MainTabScreenProps} from '@types/navigation';

export const ProfileScreen: React.FC<MainTabScreenProps<'Profile'>> = () => {
  const dispatch = useDispatch<AppDispatch>();
  const {user, isLoading} = useSelector((state: RootState) => state.auth);

  const handleLogout = () => {
    Alert.alert(
      'Выход',
      'Вы уверены, что хотите выйти из аккаунта?',
      [
        {
          text: 'Отмена',
          style: 'cancel',
        },
        {
          text: 'Выйти',
          style: 'destructive',
          onPress: () => dispatch(logoutUser()),
        },
      ],
    );
  };

  const profileSections = [
    {
      title: 'Основная информация',
      items: [
        {
          label: 'Номер телефона',
          value: user?.phone || 'Не указан',
          icon: 'phone',
        },
        {
          label: 'Email',
          value: user?.email || 'Не указан',
          icon: 'email',
        },
        {
          label: 'Дата регистрации',
          value: user?.created_at 
            ? new Date(user.created_at).toLocaleDateString('ru-RU')
            : 'Не указана',
          icon: 'calendar-today',
        },
      ],
    },
    {
      title: 'Статистика',
      items: [
        {
          label: 'Завершенных сценариев',
          value: '12',
          icon: 'check-circle',
        },
        {
          label: 'Общее время',
          value: '24 часа',
          icon: 'schedule',
        },
        {
          label: 'Средний рейтинг',
          value: '4.8 ⭐',
          icon: 'star',
        },
      ],
    },
  ];

  const actionItems = [
    {
      title: 'Редактировать профиль',
      icon: 'edit',
      color: theme.colors.primary,
      onPress: () => {
        // TODO: Navigate to edit profile screen
        Alert.alert('Функция', 'Редактирование профиля будет добавлено в следующих версиях');
      },
    },
    {
      title: 'Уведомления',
      icon: 'notifications',
      color: theme.colors.info,
      onPress: () => {
        // TODO: Navigate to notifications settings
        Alert.alert('Функция', 'Настройки уведомлений будут добавлены в следующих версиях');
      },
    },
    {
      title: 'Конфиденциальность',
      icon: 'security',
      color: theme.colors.warning,
      onPress: () => {
        // TODO: Navigate to privacy settings
        Alert.alert('Функция', 'Настройки конфиденциальности будут добавлены в следующих версиях');
      },
    },
    {
      title: 'Поддержка',
      icon: 'help',
      color: theme.colors.secondary,
      onPress: () => {
        Alert.alert('Поддержка', 'support@lunary.app\n\nМы ответим в течение 24 часов');
      },
    },
  ];

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <View style={styles.avatar}>
            <Icon name="person" size={48} color={theme.colors.white} />
          </View>
          <Text style={styles.userName}>
            {user?.phone ? `ID: ${user.id}` : 'Пользователь'}
          </Text>
          <Text style={styles.userPhone}>{user?.phone}</Text>
        </View>

        {profileSections.map((section, sectionIndex) => (
          <View key={sectionIndex} style={styles.section}>
            <Text style={styles.sectionTitle}>{section.title}</Text>
            <View style={styles.sectionContent}>
              {section.items.map((item, itemIndex) => (
                <View key={itemIndex} style={styles.infoItem}>
                  <View style={styles.infoIcon}>
                    <Icon name={item.icon} size={20} color={theme.colors.primary} />
                  </View>
                  <View style={styles.infoContent}>
                    <Text style={styles.infoLabel}>{item.label}</Text>
                    <Text style={styles.infoValue}>{item.value}</Text>
                  </View>
                </View>
              ))}
            </View>
          </View>
        ))}

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Действия</Text>
          <View style={styles.sectionContent}>
            {actionItems.map((action, index) => (
              <TouchableOpacity
                key={index}
                style={styles.actionItem}
                onPress={action.onPress}>
                <View style={[styles.actionIcon, {backgroundColor: action.color}]}>
                  <Icon name={action.icon} size={20} color={theme.colors.white} />
                </View>
                <Text style={styles.actionTitle}>{action.title}</Text>
                <Icon name="chevron-right" size={20} color={theme.colors.gray400} />
              </TouchableOpacity>
            ))}
          </View>
        </View>

        <View style={styles.section}>
          <TouchableOpacity
            style={styles.logoutButton}
            onPress={handleLogout}
            disabled={isLoading}>
            <Icon 
              name="logout" 
              size={20} 
              color={theme.colors.error} 
            />
            <Text style={styles.logoutText}>
              {isLoading ? 'Выход...' : 'Выйти из аккаунта'}
            </Text>
          </TouchableOpacity>
        </View>

        <View style={styles.footer}>
          <Text style={styles.footerText}>LUNARY v1.0.0</Text>
          <Text style={styles.footerText}>© 2024 LUNARY. Все права защищены.</Text>
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
    backgroundColor: theme.colors.primary,
    alignItems: 'center',
    padding: theme.spacing.xl,
    paddingTop: theme.spacing.xxl,
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: theme.colors.primaryDark,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: theme.spacing.lg,
  },
  userName: {
    fontSize: theme.fontSize.xl,
    fontWeight: theme.fontWeight.bold,
    color: theme.colors.white,
    marginBottom: theme.spacing.sm,
  },
  userPhone: {
    fontSize: theme.fontSize.md,
    color: theme.colors.primaryLight,
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
  infoItem: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: theme.colors.gray100,
  },
  infoIcon: {
    marginRight: theme.spacing.md,
  },
  infoContent: {
    flex: 1,
  },
  infoLabel: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textMuted,
    marginBottom: theme.spacing.xs,
  },
  infoValue: {
    fontSize: theme.fontSize.md,
    fontWeight: theme.fontWeight.medium,
    color: theme.colors.text,
  },
  actionItem: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: theme.colors.gray100,
  },
  actionIcon: {
    width: 36,
    height: 36,
    borderRadius: theme.borderRadius.md,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: theme.spacing.md,
  },
  actionTitle: {
    flex: 1,
    fontSize: theme.fontSize.md,
    fontWeight: theme.fontWeight.medium,
    color: theme.colors.text,
  },
  logoutButton: {
    backgroundColor: theme.colors.white,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.md,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    ...theme.shadows.sm,
  },
  logoutText: {
    fontSize: theme.fontSize.md,
    fontWeight: theme.fontWeight.medium,
    color: theme.colors.error,
    marginLeft: theme.spacing.sm,
  },
  footer: {
    alignItems: 'center',
    padding: theme.spacing.xl,
    marginTop: theme.spacing.lg,
  },
  footerText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textMuted,
    textAlign: 'center',
    marginBottom: theme.spacing.xs,
  },
});
