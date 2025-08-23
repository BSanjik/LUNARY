export interface Scenario {
  id: number;
  title: string;
  description: string;
  category: string;
  duration: string;
  participants: string;
  budget: string;
  checklist: string[];
  created_at: string;
}

export interface Step {
  id: number;
  scenario_id: number;
  step_number: number;
  title: string;
  description: string;
  action_type: string;
  is_completed: boolean;
}

export interface Feedback {
  id: number;
  scenario_id: number;
  user_id?: number;
  rating: number;
  comment?: string;
  created_at: string;
}

export interface FeedbackRequest {
  scenario_id: number;
  rating: number;
  comment?: string;
  user_login: string;
}

export interface ScenarioState {
  scenarios: Scenario[];
  currentScenario: Scenario | null;
  steps: Step[];
  isLoading: boolean;
  error: string | null;
}

export interface ApiError {
  code: string;
  message: string;
  details?: string;
}

export interface ApiResponse<T = any> {
  data?: T;
  error?: ApiError;
}