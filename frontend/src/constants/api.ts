// API Configuration
export const API_CONFIG = {
  // Development - через API Gateway
  DEV_BASE_URL: 'http://localhost:8080',
  
  // Direct service URLs (для тестирования)
  DEV_AUTH_URL: 'http://localhost:8081',
  DEV_SCENARIO_URL: 'http://localhost:8082',
  
  // Production
  PROD_BASE_URL: 'https://api.lunary.app',
  
  // Timeouts
  TIMEOUT: 15000,
  
  // Retry configuration
  MAX_RETRIES: 3,
  RETRY_DELAY: 1000,
};

// API Endpoints
export const ENDPOINTS = {
  // Auth endpoints
  AUTH: {
    REGISTRATION: '/api/v1/registration',
    LOGIN: '/api/v1/login',
    REFRESH: '/api/v1/refresh',
    LOGOUT: '/api/v1/logout',
  },
  
  // Scenario endpoints (соответствуют backend API)
  SCENARIOS: {
    SEARCH: '/search',
    FEEDBACK: '/feedback',
  },
  
  // Health endpoints
  HEALTH: {
    LIVE: '/health/live',
    READY: '/health/ready',
    DETAILED: '/health',
  },
} as const;

// HTTP Status Codes
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  INTERNAL_SERVER_ERROR: 500,
  SERVICE_UNAVAILABLE: 503,
} as const;

// Error Codes (from backend)
export const ERROR_CODES = {
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  AUTH_ERROR: 'AUTH_ERROR',
  NOT_FOUND: 'NOT_FOUND',
  CONFLICT_ERROR: 'CONFLICT_ERROR',
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  NETWORK_ERROR: 'NETWORK_ERROR',
} as const;
