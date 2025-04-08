import React, { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useDispatch } from 'react-redux';
import { AppDispatch } from '../redux/store';
import { showToast } from '../redux/slices/toastSlice';
import { useUpdateUserRoleMutation } from '../redux/services/authApi';

const roles = [
  {
    id: 'student',
    title: 'Student',
    description: 'Access educational resources, courses, and learning materials.',
    icon: '🎓'
  },
  {
    id: 'professional',
    title: 'Professional',
    description: 'Enhance your skills and connect with industry experts.',
    icon: '💼'
  },
  {
    id: 'instructor',
    title: 'Instructor',
    description: 'Create and share educational content with learners.',
    icon: '👨‍🏫'
  },
  {
    id: 'corporate_admin',
    title: 'Corporate Administrator',
    description: 'Manage corporate learning programs and employee development.',
    icon: '🏢'
  },
  {
    id: 'college_admin',
    title: 'College Administrator',
    description: 'Oversee educational programs and student management.',
    icon: '🏛️'
  }
];

const RoleSelection: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  const [selectedRole, setSelectedRole] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  
  const token = searchParams.get('token');
  const [updateUserRole] = useUpdateUserRoleMutation();

  if (!token) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <div>
            <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
              Invalid Verification Link
            </h2>
            <p className="mt-2 text-center text-sm text-gray-600">
              This verification link is invalid or has expired.
            </p>
          </div>
        </div>
      </div>
    );
  }

  const handleRoleSelection = (roleId: string) => {
    setSelectedRole(roleId);
  };

  const handleSubmit = async () => {
    if (!selectedRole) {
      dispatch(showToast({ message: 'Please select a role', type: 'error' }));
      return;
    }

    setIsSubmitting(true);
    try {
      await updateUserRole({ token, role: selectedRole }).unwrap();
      dispatch(showToast({ 
        message: 'Role selection successful! Your account is being verified.', 
        type: 'success' 
      }));
      
      // Show a message about the verification process
      setTimeout(() => {
        navigate('/auth');
      }, 3000);
    } catch (error) {
      dispatch(showToast({ 
        message: 'Failed to update role. Please try again.', 
        type: 'error' 
      }));
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <div className="text-center mb-8">
          <h2 className="text-3xl font-extrabold text-gray-900">
            Select Your Role
          </h2>
          <p className="mt-2 text-lg text-gray-600">
            Please choose your role carefully as features vary significantly between roles.
          </p>
        </div>

        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <ul className="divide-y divide-gray-200">
            {roles.map((role) => (
              <li key={role.id}>
                <div 
                  className={`block hover:bg-gray-50 cursor-pointer ${
                    selectedRole === role.id ? 'bg-indigo-50' : ''
                  }`}
                  onClick={() => handleRoleSelection(role.id)}
                >
                  <div className="px-4 py-4 sm:px-6">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center">
                        <p className="text-2xl mr-4">{role.icon}</p>
                        <p className="text-sm font-medium text-indigo-600 truncate">
                          {role.title}
                        </p>
                      </div>
                      <div className="ml-2 flex-shrink-0 flex">
                        {selectedRole === role.id && (
                          <p className="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                            Selected
                          </p>
                        )}
                      </div>
                    </div>
                    <div className="mt-2 sm:flex sm:justify-between">
                      <div className="sm:flex">
                        <p className="flex items-center text-sm text-gray-500">
                          {role.description}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>

        <div className="mt-8 bg-yellow-50 border-l-4 border-yellow-400 p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-yellow-800">
                Important Notice
              </h3>
              <div className="mt-2 text-sm text-yellow-700">
                <p>
                  We take spam seriously and to ensure every single account is genuine on the site, 
                  we may not accept certain unverifiable accounts and reserve the right to de-platform 
                  them without any notice. Your registration details will be validated, which may take 
                  up to 8-24 hours as the process of verifying your LinkedIn ID is manual.
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="mt-8 flex justify-center">
          <button
            type="button"
            onClick={handleSubmit}
            disabled={!selectedRole || isSubmitting}
            className={`w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white ${
              !selectedRole || isSubmitting
                ? 'bg-indigo-300 cursor-not-allowed'
                : 'bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
            }`}
          >
            {isSubmitting ? 'Submitting...' : 'Confirm Role Selection'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default RoleSelection; 