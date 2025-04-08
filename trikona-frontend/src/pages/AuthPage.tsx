// pages/AuthPage.tsx
import React from 'react';
import {AuthForm} from '../components/index';

const AuthPage: React.FC = () => {
  return (
    <div className="min-h-screen flex flex-col md:flex-row">
      {/* Left side: Marketing/illustration */}
      <div className="md:w-1/2 bg-gradient-to-br from-blue-500 to-indigo-600 flex flex-col justify-center items-center p-8 text-white">
        <h1 className="text-4xl font-bold mb-4">Welcome to Trikona</h1>
        <p className="text-lg mb-8 text-center">
          Secure, modern authentication made simple.
        </p>
        <img
          src="/auth-illustration.svg"
          alt="Authentication illustration"
          className="w-64 h-64"
        />
      </div>
      {/* Right side: Auth form */}
      <div className="md:w-1/2 flex flex-col justify-center items-center p-8 bg-gray-50">
        <AuthForm />
      </div>
    </div>
  );
};

export default AuthPage;
