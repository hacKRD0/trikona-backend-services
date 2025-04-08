import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../redux/store';

const Dashboard: React.FC = () => {
  const { token } = useSelector((state: RootState) => state.auth);

  if (!token) {
    return null; // This will be handled by the router's protected route
  }

  return (
    <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div className="px-4 py-6 sm:px-0">
        <h1 className="text-2xl font-semibold text-gray-900">Dashboard</h1>
        <div className="mt-6 bg-white shadow rounded-lg p-6">
          <p className="text-gray-600">
            Welcome to your dashboard! This is where you'll see your personalized content.
          </p>
        </div>
      </div>
    </div>
  );
};

export default Dashboard; 