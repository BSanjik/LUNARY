import {apiService} from './api';
import {ENDPOINTS} from '@constants/api';
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  ApiResponse,
} from '../types/auth';

class AuthService {
  // User registration
  public async register(
    data: RegisterRequest,
  ): Promise<ApiResponse<AuthResponse>> {
    return apiService.post<AuthResponse>(ENDPOINTS.AUTH.REGISTRATION, data);
  }

  // User login
  public async login(data: LoginRequest): Promise<ApiResponse<AuthResponse>> {
    return apiService.post<AuthResponse>(ENDPOINTS.AUTH.LOGIN, data);
  }

  // User logout
  public async logout(): Promise<ApiResponse<void>> {
    const response = await apiService.post<void>(ENDPOINTS.AUTH.LOGOUT);
    
    // Clear token regardless of response
    await apiService.clearToken();
    
    return response;
  }

  // Refresh token
  public async refreshToken(): Promise<ApiResponse<AuthResponse>> {
    return apiService.post<AuthResponse>(ENDPOINTS.AUTH.REFRESH);
  }

  // Validate phone number format
  public validatePhone(phone: string): boolean {
    // E.164 format validation
    const phoneRegex = /^\+[1-9]\d{1,14}$/;
    return phoneRegex.test(phone);
  }

  // Validate password strength
  public validatePassword(password: string): {
    isValid: boolean;
    errors: string[];
  } {
    const errors: string[] = [];
    
    if (password.length < 6) {
      errors.push('Password must be at least 6 characters');
    }
    
    if (!/[A-Z]/.test(password)) {
      errors.push('Password must contain at least one uppercase letter');
    }
    
    if (!/[a-z]/.test(password)) {
      errors.push('Password must contain at least one lowercase letter');
    }
    
    if (!/\d/.test(password)) {
      errors.push('Password must contain at least one number');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  // Format phone number for display
  public formatPhoneNumber(phone: string): string {
    // Remove all non-digit characters except +
    const cleaned = phone.replace(/[^\d+]/g, '');
    
    // Add + if not present
    if (!cleaned.startsWith('+')) {
      return `+${cleaned}`;
    }
    
    return cleaned;
  }
}

export const authService = new AuthService();
