// slices/toastSlice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { toast, ToastPosition } from 'react-toastify';

interface ToastState {
  message: string;
  type: string;
  position: ToastPosition;
  duration: number;
}

const initialState: ToastState = {
  message: '',
  type: 'default',
  position: 'top-right' as ToastPosition,
  duration: 1000, // Duration in milliseconds for toast messages
};

export const toastSlice = createSlice({
  name: 'toast',
  initialState,
  reducers: {
    showToast: (state, action: PayloadAction<{ message: string; type: string }>) => {
      const toastOptions = {
        duration: state.duration,
        position: state.position,
      };

      // Display the toast using react-toastify based on the type
      switch (action.payload.type) {
        case 'success':
          toast.success(action.payload.message, toastOptions);
          break;
        case 'error':
          toast.error(action.payload.message, toastOptions);
          break;
        case 'loading':
          toast.loading(action.payload.message, toastOptions);
          break;
        default:
          toast(action.payload.message, toastOptions);
          break;
      }

      // Update the slice state (optional, if you need to track the current toast)
      state.message = action.payload.message;
      state.type = action.payload.type || 'default';
    },
    // Optional: you might also add a hideToast action if you want to clear the toast state.
    hideToast: (state) => {
      state.message = '';
      state.type = 'default';
    },
  },
});

export const { showToast, hideToast } = toastSlice.actions;
export default toastSlice.reducer;
