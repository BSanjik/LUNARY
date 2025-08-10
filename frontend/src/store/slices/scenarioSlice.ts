import {createSlice, createAsyncThunk, PayloadAction} from '@reduxjs/toolkit';
import {scenarioService} from '@services/scenarioService';
import type {
  ScenarioState,
  Scenario,
  Step,
  Feedback,
  FeedbackRequest,
  ApiError,
} from '@types/scenario';

// Initial state
const initialState: ScenarioState = {
  scenarios: [],
  currentScenario: null,
  steps: [],
  isLoading: false,
  error: null,
};

// Async thunks
export const fetchScenarios = createAsyncThunk(
  'scenario/fetchScenarios',
  async (_, {rejectWithValue}) => {
    const response = await scenarioService.getScenarios();
    
    if (response.error) {
      return rejectWithValue(response.error);
    }
    
    return response.data || [];
  },
);

export const fetchScenario = createAsyncThunk(
  'scenario/fetchScenario',
  async (id: number, {rejectWithValue}) => {
    const response = await scenarioService.getScenario(id);
    
    if (response.error) {
      return rejectWithValue(response.error);
    }
    
    return response.data;
  },
);

export const fetchScenarioSteps = createAsyncThunk(
  'scenario/fetchScenarioSteps',
  async (scenarioId: number, {rejectWithValue}) => {
    const response = await scenarioService.getScenarioSteps(scenarioId);
    
    if (response.error) {
      return rejectWithValue(response.error);
    }
    
    return response.data || [];
  },
);

export const submitScenarioFeedback = createAsyncThunk(
  'scenario/submitFeedback',
  async (
    {scenarioId, feedback}: {scenarioId: number; feedback: FeedbackRequest},
    {rejectWithValue},
  ) => {
    const response = await scenarioService.submitFeedback(scenarioId, feedback);
    
    if (response.error) {
      return rejectWithValue(response.error);
    }
    
    return response.data;
  },
);

// Scenario slice
const scenarioSlice = createSlice({
  name: 'scenario',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setCurrentScenario: (state, action: PayloadAction<Scenario | null>) => {
      state.currentScenario = action.payload;
    },
    updateStepCompletion: (
      state,
      action: PayloadAction<{stepId: number; isCompleted: boolean}>,
    ) => {
      const {stepId, isCompleted} = action.payload;
      const step = state.steps.find(s => s.id === stepId);
      if (step) {
        step.is_completed = isCompleted;
      }
    },
    clearSteps: (state) => {
      state.steps = [];
    },
  },
  extraReducers: (builder) => {
    // Fetch scenarios
    builder
      .addCase(fetchScenarios.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchScenarios.fulfilled, (state, action) => {
        state.isLoading = false;
        state.scenarios = action.payload;
        state.error = null;
      })
      .addCase(fetchScenarios.rejected, (state, action) => {
        state.isLoading = false;
        state.error = (action.payload as ApiError)?.message || 'Failed to fetch scenarios';
      });

    // Fetch scenario
    builder
      .addCase(fetchScenario.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchScenario.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentScenario = action.payload;
        state.error = null;
      })
      .addCase(fetchScenario.rejected, (state, action) => {
        state.isLoading = false;
        state.error = (action.payload as ApiError)?.message || 'Failed to fetch scenario';
      });

    // Fetch scenario steps
    builder
      .addCase(fetchScenarioSteps.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchScenarioSteps.fulfilled, (state, action) => {
        state.isLoading = false;
        state.steps = action.payload;
        state.error = null;
      })
      .addCase(fetchScenarioSteps.rejected, (state, action) => {
        state.isLoading = false;
        state.error = (action.payload as ApiError)?.message || 'Failed to fetch scenario steps';
      });

    // Submit feedback
    builder
      .addCase(submitScenarioFeedback.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(submitScenarioFeedback.fulfilled, (state) => {
        state.isLoading = false;
        state.error = null;
      })
      .addCase(submitScenarioFeedback.rejected, (state, action) => {
        state.isLoading = false;
        state.error = (action.payload as ApiError)?.message || 'Failed to submit feedback';
      });
  },
});

export const {
  clearError,
  setCurrentScenario,
  updateStepCompletion,
  clearSteps,
} = scenarioSlice.actions;
export default scenarioSlice.reducer;
