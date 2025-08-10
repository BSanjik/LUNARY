import axios, {
  AxiosInstance,
  AxiosRequestConfig,
  AxiosResponse,
  AxiosError,
} from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';
import {API_CONFIG, HTTP_STATUS, ERROR_CODES} from '@constants/api';
import type {ApiResponse, ApiError} from '@types/auth';

class ApiService {
  private instance: AxiosInstance;
  private token: string | null = null;

  constructor() {
    this.instance = axios.create({
      baseURL: __DEV__ ? API_CONFIG.DEV_BASE_URL : API_CONFIG.PROD_BASE_URL,
      timeout: API_CONFIG.TIMEOUT,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
    this.loadToken();
  }

  private setupInterceptors() {
    // Request interceptor
    this.instance.interceptors.request.use(
      (config) => {
        if (this.token) {
          config.headers.Authorization = `Bearer ${this.token}`;
        }
        return config;
      },
      (error) => Promise.reject(error),
    );

    // Response interceptor
    this.instance.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        if (error.response?.status === HTTP_STATUS.UNAUTHORIZED) {
          await this.clearToken();
          // Redirect to login screen
        }
        return Promise.reject(this.handleError(error));
      },
    );
  }

  private async loadToken() {
    try {
      const token = await AsyncStorage.getItem('authToken');
      if (token) {
        this.token = token;
      }
    } catch (error) {
      console.error('Error loading token:', error);
    }
  }

  public async setToken(token: string) {
    this.token = token;
    try {
      await AsyncStorage.setItem('authToken', token);
    } catch (error) {
      console.error('Error saving token:', error);
    }
  }

  public async clearToken() {
    this.token = null;
    try {
      await AsyncStorage.removeItem('authToken');
    } catch (error) {
      console.error('Error clearing token:', error);
    }
  }

  private handleError(error: AxiosError): ApiError {
    if (!error.response) {
      return {
        code: ERROR_CODES.NETWORK_ERROR,
        message: 'Network error. Please check your connection.',
      };
    }

    const {data, status} = error.response;
    
    // Handle structured errors from backend
    if (data && typeof data === 'object' && 'code' in data) {
      return data as ApiError;
    }

    // Handle HTTP status codes
    switch (status) {
      case HTTP_STATUS.BAD_REQUEST:
        return {
          code: ERROR_CODES.VALIDATION_ERROR,
          message: 'Invalid request data',
        };
      case HTTP_STATUS.UNAUTHORIZED:
        return {
          code: ERROR_CODES.AUTH_ERROR,
          message: 'Authentication required',
        };
      case HTTP_STATUS.NOT_FOUND:
        return {
          code: ERROR_CODES.NOT_FOUND,
          message: 'Resource not found',
        };
      case HTTP_STATUS.INTERNAL_SERVER_ERROR:
        return {
          code: ERROR_CODES.INTERNAL_ERROR,
          message: 'Internal server error',
        };
      default:
        return {
          code: ERROR_CODES.INTERNAL_ERROR,
          message: 'Something went wrong',
        };
    }
  }

  // Generic request methods
  public async get<T>(
    url: string,
    config?: AxiosRequestConfig,
  ): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.instance.get(url, config);
      return {data: response.data};
    } catch (error) {
      return {error: error as ApiError};
    }
  }

  public async post<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.instance.post(
        url,
        data,
        config,
      );
      return {data: response.data};
    } catch (error) {
      return {error: error as ApiError};
    }
  }

  public async put<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig,
  ): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.instance.put(
        url,
        data,
        config,
      );
      return {data: response.data};
    } catch (error) {
      return {error: error as ApiError};
    }
  }

  public async delete<T>(
    url: string,
    config?: AxiosRequestConfig,
  ): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<T> = await this.instance.delete(url, config);
      return {data: response.data};
    } catch (error) {
      return {error: error as ApiError};
    }
  }
}

export const apiService = new ApiService();
