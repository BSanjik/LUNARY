import {apiService} from './api';
import {ENDPOINTS} from '@constants/api';
import type {
  Scenario,
  Step,
  Feedback,
  FeedbackRequest,
  ApiResponse,
} from '@types/scenario';

class ScenarioService {
  // Get all scenarios
  public async getScenarios(): Promise<ApiResponse<Scenario[]>> {
    return apiService.get<Scenario[]>(ENDPOINTS.SCENARIOS.LIST);
  }

  // Get scenario by ID
  public async getScenario(id: number): Promise<ApiResponse<Scenario>> {
    const url = ENDPOINTS.SCENARIOS.DETAIL.replace(':id', id.toString());
    return apiService.get<Scenario>(url);
  }

  // Get scenario steps
  public async getScenarioSteps(scenarioId: number): Promise<ApiResponse<Step[]>> {
    const url = ENDPOINTS.SCENARIOS.STEPS.replace(':id', scenarioId.toString());
    return apiService.get<Step[]>(url);
  }

  // Submit feedback for scenario
  public async submitFeedback(
    scenarioId: number,
    feedback: FeedbackRequest,
  ): Promise<ApiResponse<Feedback>> {
    const url = ENDPOINTS.SCENARIOS.FEEDBACK.replace(':id', scenarioId.toString());
    return apiService.post<Feedback>(url, feedback);
  }

  // Validate feedback rating
  public validateRating(rating: number): boolean {
    return rating >= 1 && rating <= 5 && Number.isInteger(rating);
  }

  // Format difficulty level
  public formatDifficulty(difficulty: string): string {
    const difficultyMap: Record<string, string> = {
      easy: 'Легкий',
      medium: 'Средний',
      hard: 'Сложный',
    };
    
    return difficultyMap[difficulty] || difficulty;
  }

  // Format estimated time
  public formatEstimatedTime(minutes: number): string {
    if (minutes < 60) {
      return `${minutes} мин`;
    }
    
    const hours = Math.floor(minutes / 60);
    const remainingMinutes = minutes % 60;
    
    if (remainingMinutes === 0) {
      return `${hours} ч`;
    }
    
    return `${hours} ч ${remainingMinutes} мин`;
  }

  // Get difficulty color
  public getDifficultyColor(difficulty: string): string {
    const colorMap: Record<string, string> = {
      easy: '#10B981', // green
      medium: '#F59E0B', // yellow
      hard: '#EF4444', // red
    };
    
    return colorMap[difficulty] || '#6B7280'; // gray
  }

  // Calculate progress percentage
  public calculateProgress(steps: Step[]): number {
    if (steps.length === 0) return 0;
    
    const completedSteps = steps.filter(step => step.is_completed).length;
    return Math.round((completedSteps / steps.length) * 100);
  }
}

export const scenarioService = new ScenarioService();
