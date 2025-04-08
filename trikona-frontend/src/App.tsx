// App.tsx
import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { useSelector } from 'react-redux';
import { RootState } from './redux/store';
import Layout from './components/layout/Layout';
import AuthPage from './pages/AuthPage';
import LinkedInCallback from './pages/LinkedInCallback';
import Dashboard from './pages/Dashboard';
import Profile from './pages/Profile';
import RoleSelection from './pages/RoleSelection';
import RegistrationForm from './pages/RegistrationForm';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { token } = useSelector((state: RootState) => state.auth);
  
  if (!token) {
    return <Navigate to="/auth" replace />;
  }

  return <>{children}</>;
};

const App: React.FC = () => {
  return (
    <Routes>
      <Route path="/auth" element={<AuthPage />} />
      <Route path="/linkedin/callback" element={<LinkedInCallback />} />
      <Route path="/auth/role-selection" element={<RoleSelection />} />
      <Route path="/auth/register" element={<RegistrationForm />} />
      
      <Route path="/" element={<Layout />}>
        <Route index element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        } />
        <Route path="profile" element={
          <ProtectedRoute>
            <Profile />
          </ProtectedRoute>
        } />
      </Route>
    </Routes>
  );
};

export default App;
