// slices/authSlice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { IUser } from '../services/authApi';

export type ToastType = 'error' | 'success' | 'info' | 'warning';

interface Toast {
  message: string;
  type: ToastType;
}

interface AuthState {
  token: string | null;
  user: IUser | null;
  loading: boolean;
  error: string | null;
  toast: Toast | null;
}

const initialState: AuthState = {
  token: localStorage.getItem('token'),
  user: null,
  loading: false,
  error: null,
  toast: null,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setCredentials: (
      state,
      action: PayloadAction<{ token: string; user: IUser | null }>
    ) => {
      const { token, user } = action.payload;
      state.token = token;
      state.user = user;
      localStorage.setItem('token', token);
    },
    setLoading: (state, action: PayloadAction<boolean>) => {
      state.loading = action.payload;
    },
    setError: (state, action: PayloadAction<string | null>) => {
      state.error = action.payload;
    },
    logout: (state) => {
      state.token = null;
      state.user = null;
      localStorage.removeItem('token');
    },
    // Action to show a toast notification with a custom message and type.
    showToast(state, action: PayloadAction<Toast>) {
      state.toast = action.payload;
    },
    // Action to clear/hide the current toast notification.
    hideToast(state) {
      state.toast = null;
    },
  },
});

export const { setCredentials, setLoading, setError, logout, showToast, hideToast } = authSlice.actions;
export default authSlice.reducer;
