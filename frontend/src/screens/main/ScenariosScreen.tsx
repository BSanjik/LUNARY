import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  RefreshControl,
} from 'react-native';
import {SafeAreaView} from 'react-native-safe-area-context';
import {useDispatch, useSelector} from 'react-redux';
import Icon from 'react-native-vector-icons/MaterialIcons';

import {fetchScenarios} from '@store/slices/scenarioSlice';
import {scenarioService} from '@services/scenarioService';
import {theme} from '@constants/theme';
import type {AppDispatch, RootState} from '@store/index';
import type {MainTabScreenProps} from '@types/navigation';
import type {Scenario} from '@types/scenario';

export const ScenariosScreen: React.FC<MainTabScreenProps<'Scenarios'>> = () => {
  const dispatch = useDispatch<AppDispatch>();
  const {scenarios, isLoading} = useSelector((state: RootState) => state.scenario);

  React.useEffect(() => {
    dispatch(fetchScenarios());
  }, [dispatch]);

  const onRefresh = () => {
    dispatch(fetchScenarios());
  };

  const renderScenarioItem = ({item}: {item: Scenario}) => (
    <TouchableOpacity style={styles.scenarioCard}>
      <View style={styles.scenarioHeader}>
        <View style={styles.scenarioInfo}>
          <Text style={styles.scenarioTitle}>{item.title}</Text>
          <Text style={styles.scenarioDescription} numberOfLines={2}>
            {item.description}
          </Text>
        </View>
        <View style={styles.scenarioMeta}>
          <View
            style={[
              styles.difficultyBadge,
              {backgroundColor: scenarioService.getDifficultyColor(item.difficulty)},
            ]}>
            <Text style={styles.difficultyText}>
              {scenarioService.formatDifficulty(item.difficulty)}
            </Text>
          </View>
        </View>
      </View>

      <View style={styles.scenarioFooter}>
        <View style={styles.timeContainer}>
          <Icon name="schedule" size={16} color={theme.colors.textMuted} />
          <Text style={styles.timeText}>
            {scenarioService.formatEstimatedTime(item.estimated_time)}
          </Text>
        </View>
        
        <View style={styles.statusContainer}>
          {item.is_active ? (
            <View style={styles.activeStatus}>
              <Icon name="check-circle" size={16} color={theme.colors.success} />
              <Text style={styles.activeText}>Доступен</Text>
            </View>
          ) : (
            <View style={styles.inactiveStatus}>
              <Icon name="pause-circle-filled" size={16} color={theme.colors.gray400} />
              <Text style={styles.inactiveText}>Недоступен</Text>
            </View>
          )}
        </View>

        <TouchableOpacity style={styles.startButton}>
          <Text style={styles.startButtonText}>Начать</Text>
          <Icon name="play-arrow" size={16} color={theme.colors.white} />
        </TouchableOpacity>
      </View>
    </TouchableOpacity>
  );

  const renderEmptyState = () => (
    <View style={styles.emptyState}>
      <Icon name="list-alt" size={64} color={theme.colors.gray300} />
      <Text style={styles.emptyTitle}>Нет доступных сценариев</Text>
      <Text style={styles.emptyDescription}>
        Сценарии появятся здесь, когда они будут добавлены
      </Text>
      <TouchableOpacity style={styles.refreshButton} onPress={onRefresh}>
        <Icon name="refresh" size={20} color={theme.colors.primary} />
        <Text style={styles.refreshText}>Обновить</Text>
      </TouchableOpacity>
    </View>
  );

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Сценарии</Text>
        <Text style={styles.headerSubtitle}>
          Доступно {scenarios.filter(s => s.is_active).length} сценариев
        </Text>
      </View>

      <FlatList
        data={scenarios}
        renderItem={renderScenarioItem}
        keyExtractor={(item) => item.id.toString()}
        contentContainerStyle={styles.listContent}
        refreshControl={
          <RefreshControl
            refreshing={isLoading}
            onRefresh={onRefresh}
            colors={[theme.colors.primary]}
            tintColor={theme.colors.primary}
          />
        }
        ListEmptyComponent={renderEmptyState}
        showsVerticalScrollIndicator={false}
      />
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: theme.colors.backgroundSecondary,
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
  listContent: {
    padding: theme.spacing.lg,
    paddingBottom: theme.spacing.xxl,
  },
  scenarioCard: {
    backgroundColor: theme.colors.white,
    borderRadius: theme.borderRadius.lg,
    padding: theme.spacing.lg,
    marginBottom: theme.spacing.md,
    ...theme.shadows.md,
  },
  scenarioHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: theme.spacing.md,
  },
  scenarioInfo: {
    flex: 1,
    marginRight: theme.spacing.md,
  },
  scenarioTitle: {
    fontSize: theme.fontSize.lg,
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.text,
    marginBottom: theme.spacing.sm,
  },
  scenarioDescription: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textSecondary,
    lineHeight: 20,
  },
  scenarioMeta: {
    alignItems: 'flex-end',
  },
  difficultyBadge: {
    paddingHorizontal: theme.spacing.sm,
    paddingVertical: theme.spacing.xs,
    borderRadius: theme.borderRadius.full,
  },
  difficultyText: {
    fontSize: theme.fontSize.sm,
    fontWeight: theme.fontWeight.medium,
    color: theme.colors.white,
  },
  scenarioFooter: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingTop: theme.spacing.md,
    borderTopWidth: 1,
    borderTopColor: theme.colors.gray100,
  },
  timeContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  timeText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.textMuted,
    marginLeft: theme.spacing.xs,
  },
  statusContainer: {
    flex: 1,
    alignItems: 'center',
  },
  activeStatus: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  activeText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.success,
    marginLeft: theme.spacing.xs,
    fontWeight: theme.fontWeight.medium,
  },
  inactiveStatus: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  inactiveText: {
    fontSize: theme.fontSize.sm,
    color: theme.colors.gray400,
    marginLeft: theme.spacing.xs,
    fontWeight: theme.fontWeight.medium,
  },
  startButton: {
    backgroundColor: theme.colors.primary,
    flexDirection: 'row',
    alignItems: 'center',
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
    borderRadius: theme.borderRadius.md,
  },
  startButtonText: {
    fontSize: theme.fontSize.sm,
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.white,
    marginRight: theme.spacing.xs,
  },
  emptyState: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: theme.spacing.xxl,
  },
  emptyTitle: {
    fontSize: theme.fontSize.lg,
    fontWeight: theme.fontWeight.semibold,
    color: theme.colors.text,
    marginTop: theme.spacing.lg,
    marginBottom: theme.spacing.sm,
  },
  emptyDescription: {
    fontSize: theme.fontSize.md,
    color: theme.colors.textSecondary,
    textAlign: 'center',
    marginBottom: theme.spacing.xl,
    paddingHorizontal: theme.spacing.xl,
  },
  refreshButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: theme.colors.primaryLight,
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: theme.spacing.md,
    borderRadius: theme.borderRadius.md,
  },
  refreshText: {
    fontSize: theme.fontSize.md,
    fontWeight: theme.fontWeight.medium,
    color: theme.colors.primary,
    marginLeft: theme.spacing.sm,
  },
});
