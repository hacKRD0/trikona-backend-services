// components/ForgotPasswordForm.tsx
import React, { useState } from 'react';
import { useForgotPasswordMutation } from '../../redux/services/authApi';
import { useDispatch } from 'react-redux';
import { AppDispatch } from '../../redux/store';
import { showToast } from '../../redux/slices/toastSlice';

const ForgotPasswordForm: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const [email, setEmail] = useState('');
  const [forgotPassword, { isLoading }] = useForgotPasswordMutation();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    try {
      const res = await forgotPassword({ email }).unwrap();
      dispatch(showToast({ message: res.message, type: 'success' }));
    } catch (err: any) {
      dispatch(showToast({ message: err?.data?.error || 'Failed to send reset email', type: 'error' }));
    }
  };

  return (
    <div className="w-full max-w-md bg-white rounded-xl shadow-lg p-8">
      <h2 className="text-2xl font-semibold text-gray-800 text-center mb-6">Forgot Password</h2>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="forgot-email" className="block text-gray-700 mb-1">
            Enter your email address
          </label>
          <input
            id="forgot-email"
            type="email"
            name="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
          />
        </div>
        <button
          type="submit"
          disabled={isLoading}
          className="w-full py-2 px-4 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition-colors"
        >
          {isLoading ? 'Sending...' : 'Send Reset Link'}
        </button>
      </form>
    </div>
  );
};

export default ForgotPasswordForm;
