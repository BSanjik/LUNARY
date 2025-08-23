// Простой тест API для проверки подключения к backend
import { authService } from '../services/authService';
import { scenarioService } from '../services/scenarioService';

export class ApiTest {
  // Тест подключения к Auth Service
  static async testAuth() {
    console.log('🔐 Тестируем Auth Service...');
    
    try {
      // Тест регистрации
      const registerResult = await authService.register({
        phone: '+77001234567',
        password: 'TestPassword123',
        email: 'test@example.com'
      });

      if (registerResult.error) {
        console.log('❌ Registration error:', registerResult.error);
      } else {
        console.log('✅ Registration successful:', registerResult.data);
      }

      // Тест логина
      const loginResult = await authService.login({
        phone: '+77001234567',
        password: 'TestPassword123'
      });

      if (loginResult.error) {
        console.log('❌ Login error:', loginResult.error);
      } else {
        console.log('✅ Login successful:', loginResult.data);
        return loginResult.data;
      }
    } catch (error) {
      console.error('❌ Auth test failed:', error);
    }
  }

  // Тест подключения к Scenario Service
  static async testScenarios() {
    console.log('📝 Тестируем Scenario Service...');
    
    try {
      // Прямой вызов к scenario service
      const searchResult = await scenarioService.directSearchScenarios('семейный ужин');

      if (searchResult.error) {
        console.log('❌ Scenario search error:', searchResult.error);
      } else {
        console.log('✅ Scenario search successful:', searchResult.data);
        return searchResult.data;
      }

      // Тест отзыва
      const feedbackResult = await scenarioService.directSubmitFeedback({
        scenario_id: 1,
        rating: 5,
        comment: 'Тест из мобильного приложения',
        user_login: 'test@example.com'
      });

      if (feedbackResult.error) {
        console.log('❌ Feedback error:', feedbackResult.error);
      } else {
        console.log('✅ Feedback successful:', feedbackResult.data);
      }
    } catch (error) {
      console.error('❌ Scenario test failed:', error);
    }
  }

  // Полный тест всех API
  static async runFullTest() {
    console.log('🚀 Запускаем полный тест API...');
    
    await this.testAuth();
    await this.testScenarios();
    
    console.log('✅ Тест завершен!');
  }
}

// Health check для всех сервисов
export const checkBackendHealth = async () => {
  const services = [
    { name: 'API Gateway', url: 'http://localhost:8080/health' },
    { name: 'Auth Service', url: 'http://localhost:8081/health' },
    { name: 'Scenario Service', url: 'http://localhost:8082/health' }
  ];

  console.log('🏥 Проверяем здоровье сервисов...');

  for (const service of services) {
    try {
      const response = await fetch(service.url);
      if (response.ok) {
        console.log(`✅ ${service.name}: OK`);
      } else {
        console.log(`❌ ${service.name}: Error ${response.status}`);
      }
    } catch (error) {
      console.log(`❌ ${service.name}: Connection failed`);
    }
  }
};
