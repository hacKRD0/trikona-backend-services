import { StrictMode, Suspense } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import './index.css'
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';
import { store } from './redux/store';
import { ToastContainer } from 'react-toastify';
import { PersistGate } from "redux-persist/integration/react";
import { persistStore } from "redux-persist";
import { LoadingFallback } from './components/index';

const persistor = persistStore(store);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Suspense fallback={<LoadingFallback />}>
      <Provider store={store}>
        <PersistGate loading={null} persistor={persistor}>
          <BrowserRouter>
            <App />
            <ToastContainer />
          </BrowserRouter>
        </PersistGate>
      </Provider>
    </Suspense>
  </StrictMode>,
)
