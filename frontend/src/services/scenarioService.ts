import {apiService} from './api';
import {ENDPOINTS, API_CONFIG} from '@constants/api';
import type {
  Scenario,
  Step,
  Feedback,
  FeedbackRequest,
  ApiResponse,
} from '../types/scenario';

interface ScenarioSearchRequest {
  query: string;
}

class ScenarioService {
  // Search scenarios by query
  public async searchScenarios(query: string): Promise<ApiResponse<Scenario>> {
    const searchData: ScenarioSearchRequest = { query };
    return apiService.post<Scenario>(ENDPOINTS.SCENARIOS.SEARCH, searchData);
  }

  // Submit feedback for scenario
  public async submitFeedback(
    feedback: FeedbackRequest,
  ): Promise<ApiResponse<{status: string; message: string}>> {
    return apiService.post<{status: string; message: string}>(
      ENDPOINTS.SCENARIOS.FEEDBACK, 
      feedback
    );
  }

  // Direct API calls to scenario service (для тестирования)
  public async directSearchScenarios(query: string): Promise<ApiResponse<Scenario>> {
    const searchData: ScenarioSearchRequest = { query };
    return apiService.post<Scenario>(`${API_CONFIG.DEV_SCENARIO_URL}/search`, searchData);
  }

  public async directSubmitFeedback(
    feedback: FeedbackRequest,
  ): Promise<ApiResponse<{status: string; message: string}>> {
    return apiService.post<{status: string; message: string}>(
      `${API_CONFIG.DEV_SCENARIO_URL}/feedback`, 
      feedback
    );
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
