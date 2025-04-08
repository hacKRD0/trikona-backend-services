// store.ts
import { combineReducers, configureStore, Middleware, isRejectedWithValue } from '@reduxjs/toolkit';
import { setupListeners } from '@reduxjs/toolkit/query';
import { persistReducer, persistStore } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import authReducer, { logout } from './slices/authSlice';
import toastReducer, { showToast } from './slices/toastSlice';
import { authApi } from './services/authApi';

// Create your root reducer by combining all slice reducers.
const rootReducer = combineReducers({
  auth: authReducer,
  toast: toastReducer,
  [authApi.reducerPath]: authApi.reducer,
});

// Configure redux-persist
const persistConfig = {
  key: 'root',
  storage,
  // Optionally, blacklist slices that you don't want to persist.
  blacklist: [authApi.reducerPath, 'toast'],
};

const persistedReducer = persistReducer(persistConfig, rootReducer);

// Custom middleware to catch RTK Query errors (e.g., unauthorized errors)
interface ErrorPayload {
  status: number;
}

const rtkQueryErrorLogger: Middleware = (store) => (next) => (action) => {
  if (isRejectedWithValue(action)) {
    const payload = action.payload as ErrorPayload | undefined;
    if (payload && (payload.status === 401 || payload.status === 403)) {
      store.dispatch(logout());
      store.dispatch(
        showToast({
          message: 'Session Expired! Please login again.',
          type: 'error',
        })
      );
    }
  }
  return next(action);
};

export const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: false, // Disabling serializableCheck for redux-persist.
    }).concat(authApi.middleware, rtkQueryErrorLogger),
});

// Enable RTK Query's automatic re-fetching, etc.
setupListeners(store.dispatch);

// Create the persistor which will be used by PersistGate.
export const persistor = persistStore(store);

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
