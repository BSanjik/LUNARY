import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {useSelector} from 'react-redux';
import Icon from 'react-native-vector-icons/MaterialIcons';

import {theme} from '@constants/theme';
import type {RootState} from '@store/index';
import type {MainTabScreenProps} from '@types/navigation';

export const HomeScreen: React.FC<MainTabScreenProps<'Home'>> = ({navigation}) => {
  const {user} = useSelector((state: RootState) => state.auth);

  const quickActions = [
    {
      title: 'Сценарии',
      description: 'Изучить доступные сценарии',
      icon: 'list',
      color: theme.colors.primary,
      onPress: () => navigation.navigate('Scenarios'),
    },
    {
      title: 'Профиль',
      description: 'Настроить личные данные',
      icon: 'person',
      color: theme.colors.secondary,
      onPress: () => navigation.navigate('Profile'),
    },
    {
      title: 'Настройки',
      description: 'Конфигурация приложения',
      icon: 'settings',
      color: theme.colors.info,
      onPress: () => navigation.navigate('Settings'),
    },
  ];

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView style={styles.scrollView}>
        <View style={styles.header}>
          <Text style={styles.welcome}>Добро пожаловать!</Text>
          <Text style={styles.userInfo}>
            {user?.phone || 'Пользователь'}
          </Text>
        </View>

        <View style={styles.statsContainer}>
          <View style={styles.statCard}>
            <Icon name="trending-up" size={32} color={theme.colors.success} />
            <Text style={styles.statNumber}>12</Text>
            <Text style={styles.statLabel}>Завершено</Text>
          </View>
          
          <View style={styles.statCard}>
            <Icon name="schedule" size={32} color={theme.colors.warning} />
            <Text style={styles.statNumber}>3</Text>
            <Text style={styles.statLabel}>В процессе</Text>
          </View>
          
          <View style={styles.statCard}>
            <Icon name="star" size={32} color={theme.colors.secondary} />
            <Text style={styles.statNumber}>4.8</Text>
            <Text style={styles.statLabel}>Рейтинг</Text>
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Быстрые действия</Text>
          <View style={styles.actionsContainer}>
            {quickActions.map((action, index) => (
              <TouchableOpacity
                key={index}
                style={styles.actionCard}
                onPress={action.onPress}>
                <View style={[styles.actionIcon, {backgroundColor: action.color}]}>
                  <Icon name={action.icon} size={24} color={theme.colors.white} />
                </View>
                <View style={styles.actionContent}>
                  <Text style={styles.actionTitle}>{action.title}</Text>
                  <Text style={styles.actionDescription}>{action.description}</Text>
                </View>
                <Icon name="chevron-right" size={24} color={theme.colors.gray400} />
              </TouchableOpacity>
            ))}
          </View>
        </View>

        <View style={styles.section}>
          <Text style={styles.sectionTitle}>Последняя активность</Text>
          <View style={styles.activityContainer}>
            <View style={styles.activityItem}>
              <View style={styles.activityIcon}>
                <Icon name="check-circle" size={20} color={theme.colors.success} />
              </View>
              <View style={styles.activityContent}>
                <Text style={styles.activityTitle}>Сценарий "Базовая настройка" завершен</Text>
                <Text style={styles.activityTime}>2 часа назад</Text>
              </View>
            </View>
            
            <View style={styles.activityItem}>
              <View style={styles.activityIcon}>
                <Icon name="play-circle-filled" size={20} color={theme.colors.info} />
              </View>
              <View style={styles.activityContent}>
                <Text style={styles.activityTitle}>Начат сценарий "Продвинутая конфигурация"</Text>
                <Text style={styles.activityTime}>Вчера</Text>
              </View>
            </View>
            
            <View style={styles.activityItem}>
              <View style={styles.activityIcon}>
                <Icon name="rate-review" size={20} color={theme.colors.warning} />
              </View>
              <View style={styles.activityContent}>
                <Text style={styles.activityTitle}>Оставлен отзыв для сценария</Text>
                <Text style={styles.activityTime}>3 дня назад</Text>
              </View>
            </View>
          </View>
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
    padding: theme.spacing.lg,
    paddingTop: theme.spacing.xl,
  },
  welcome: {
    fontSize: theme.fontSize.xxl,
    fontWeight: theme.fontWeight.bold,
    color: theme.colors.white,
    marginBottom: theme.spacing.sm,
  },
  userInfo: {
    fontSize: theme.fontSize.lg,
    color: theme.colors.primaryLight,
  },
  statsContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: theme.spacing.lg,
    marginTop: -theme.spacing.lg,
  },
  statCard: {
    backgroundColor: theme.colors.white,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.md,
    alignItems: 'center',
    flex: 1,
    marginHorizontal: theme.spacing.sm,
    ...theme.shadows.md,
  },
  statNumber: {
    fontSize: theme.fontSize.xl,
    fontWeight: theme.fontWeight.bold,
    color: theme.colors.text,
    marginTop: theme.spacing.sm,
  },
  statLabel: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textSecondary,
    marginTop: theme.spacing.xs,
  },
  section: {
    padding: theme.spacing.lg,
  },
  sectionTitle: {
    fontSize: theme.fontSize.lg,
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.text,
    marginBottom: theme.spacing.md,
  },
  actionsContainer: {
    gap: theme.spacing.md,
  },
  actionCard: {
    backgroundColor: theme.colors.white,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.md,
    flexDirection: 'row',
    alignItems: 'center',
    ...theme.shadows.sm,
  },
  actionIcon: {
    width: 48,
    height: 48,
    borderRadius: theme.borderRadius.lg,
    alignItems: 'center',
    justifyContent: 'center',
    marginRight: theme.spacing.md,
  },
  actionContent: {
    flex: 1,
  },
  actionTitle: {
    fontSize: theme.fontSize.md,
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.text,
    marginBottom: theme.spacing.xs,
  },
  actionDescription: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textSecondary,
  },
  activityContainer: {
    backgroundColor: theme.colors.white,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.md,
    ...theme.shadows.sm,
  },
  activityItem: {
    flexDirection: 'row',
    alignItems: 'flex-start',
    paddingVertical: theme.spacing.md,
    borderBottomWidth: 1,
    borderBottomColor: theme.colors.gray100,
  },
  activityIcon: {
    marginRight: theme.spacing.md,
    marginTop: theme.spacing.xs,
  },
  activityContent: {
    flex: 1,
  },
  activityTitle: {
    fontSize: theme.fontSize.md,
    color: theme.colors.text,
    marginBottom: theme.spacing.xs,
  },
  activityTime: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textMuted,
  },
});
