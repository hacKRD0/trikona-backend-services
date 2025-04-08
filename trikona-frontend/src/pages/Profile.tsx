import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../redux/store';

const Profile: React.FC = () => {
  const { token, user } = useSelector((state: RootState) => state.auth);

  if (!token) {
    return null; // This will be handled by the router's protected route
  } else {
    console.log('User : ', user);
  }

  return (
    <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div className="px-4 py-6 sm:px-0">
        <h1 className="text-2xl font-semibold text-gray-900">Your Profile</h1>
        <div className="mt-6 bg-white shadow rounded-lg p-6">
          <div className="space-y-6">
            <div>
              <h2 className="text-lg font-medium text-gray-900">Personal Information</h2>
              <p className="mt-1 text-sm text-gray-500">
                This is your profile page. You can view your information here.
              </p>
            </div>
            
            <div className="border-t border-gray-200 pt-6">
              <h3 className="text-md font-medium text-gray-900">Account Information</h3>
              <div className="mt-4 grid grid-cols-1 gap-y-6 sm:grid-cols-6 gap-x-4">
                <div className="sm:col-span-3">
                  <div className="flex flex-col">
                    <span className="text-sm font-medium text-gray-500">First Name</span>
                    <span className="mt-1 text-sm text-gray-900">{user?.firstName || 'Not provided'}</span>
                  </div>
                </div>
                
                <div className="sm:col-span-3">
                  <div className="flex flex-col">
                    <span className="text-sm font-medium text-gray-500">Last Name</span>
                    <span className="mt-1 text-sm text-gray-900">{user?.lastName || 'Not provided'}</span>
                  </div>
                </div>
                
                <div className="sm:col-span-3">
                  <div className="flex flex-col">
                    <span className="text-sm font-medium text-gray-500">Email</span>
                    <span className="mt-1 text-sm text-gray-900">{user?.email || 'Not provided'}</span>
                  </div>
                </div>
                
                <div className="sm:col-span-3">
                  <div className="flex flex-col">
                    <span className="text-sm font-medium text-gray-500">Role</span>
                    <span className="mt-1 text-sm text-gray-900">{user?.role || 'Not provided'}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Profile; 