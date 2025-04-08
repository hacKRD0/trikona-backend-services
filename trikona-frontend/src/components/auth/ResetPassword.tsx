// components/ResetPasswordForm.tsx
import React, { useState } from 'react';
import { useResetPasswordMutation } from '../../redux/services/authApi';
import { useDispatch } from 'react-redux';
import { AppDispatch } from '../../redux/store';
import { showToast } from '../../redux/slices/toastSlice';
import { useParams } from 'react-router-dom';

const ResetPasswordForm: React.FC = () => {
  const { token } = useParams<{ token: string }>();
  const dispatch = useDispatch<AppDispatch>();
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const [resetPassword, { isLoading }] = useResetPasswordMutation();

  const isStrongPassword = (password: string): boolean => {
    const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
    return regex.test(password);
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!isStrongPassword(newPassword)) {
      dispatch(
        showToast({
          message: 'Password must be at least 8 characters and include uppercase, lowercase, a digit, and a special character.',
          type: 'error',
        })
      );
      return;
    }
    if (newPassword !== confirmPassword) {
      dispatch(showToast({ message: 'Passwords do not match.', type: 'error' }));
      return;
    }
    try {
      const res = await resetPassword({ token: token!, newPassword, confirmPassword }).unwrap();
      dispatch(showToast({ message: res.message, type: 'success' }));
    } catch (err: any) {
      dispatch(showToast({ message: err?.data?.error || 'Password reset failed', type: 'error' }));
    }
  };

  return (
    <div className="w-full max-w-md bg-white rounded-xl shadow-lg p-8">
      <h2 className="text-2xl font-semibold text-gray-800 text-center mb-6">Reset Password</h2>
      <form onSubmit={handleSubmit}>
        <div className="relative mb-4">
          <label htmlFor="newPassword" className="block text-gray-700 mb-1">
            New Password
          </label>
          <input
            id="newPassword"
            type={showNewPassword ? 'text' : 'password'}
            name="newPassword"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
          />
          <button
            type="button"
            className="absolute inset-y-0 right-0 pr-3 flex items-center text-sm text-blue-600"
            onClick={() => setShowNewPassword(!showNewPassword)}
          >
            {showNewPassword ? 'Hide' : 'Show'}
          </button>
        </div>
        <div className="relative mb-6">
          <label htmlFor="confirmPassword" className="block text-gray-700 mb-1">
            Confirm Password
          </label>
          <input
            id="confirmPassword"
            type={showConfirmPassword ? 'text' : 'password'}
            name="confirmPassword"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
          />
          <button
            type="button"
            className="absolute inset-y-0 right-0 pr-3 flex items-center text-sm text-blue-600"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
          >
            {showConfirmPassword ? 'Hide' : 'Show'}
          </button>
        </div>
        <button
          type="submit"
          disabled={isLoading}
          className="w-full py-2 px-4 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition-colors"
        >
          {isLoading ? 'Resetting...' : 'Reset Password'}
        </button>
      </form>
    </div>
  );
};

export default ResetPasswordForm;
