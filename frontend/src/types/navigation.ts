import type {StackScreenProps} from '@react-navigation/stack';
import type {BottomTabScreenProps} from '@react-navigation/bottom-tabs';

// Auth Stack
export type AuthStackParamList = {
  Login: undefined;
  Register: undefined;
  ForgotPassword: undefined;
};

// Main Tab Navigator
export type MainTabParamList = {
  Home: undefined;
  Scenarios: undefined;
  Profile: undefined;
  Settings: undefined;
};

// Scenario Stack
export type ScenarioStackParamList = {
  ScenarioList: undefined;
  ScenarioDetail: {
    scenarioId: number;
  };
  ScenarioSteps: {
    scenarioId: number;
  };
  Feedback: {
    scenarioId: number;
  };
};

// Root Stack
export type RootStackParamList = {
  Auth: undefined;
  Main: undefined;
  Scenario: {
    screen: keyof ScenarioStackParamList;
    params?: ScenarioStackParamList[keyof ScenarioStackParamList];
  };
};

// Screen Props Types
export type AuthScreenProps<T extends keyof AuthStackParamList> = 
  StackScreenProps<AuthStackParamList, T>;

export type MainTabScreenProps<T extends keyof MainTabParamList> = 
  BottomTabScreenProps<MainTabParamList, T>;

export type ScenarioScreenProps<T extends keyof ScenarioStackParamList> = 
  StackScreenProps<ScenarioStackParamList, T>;

export type RootScreenProps<T extends keyof RootStackParamList> = 
  StackScreenProps<RootStackParamList, T>;
