import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch, RootState } from '../../redux/store';
import { logout } from '../../redux/slices/authSlice';
import { showToast } from '../../redux/slices/toastSlice';
import { UserCircleIcon } from '@heroicons/react/24/outline';

const Header: React.FC = () => {
  const navigate = useNavigate();
  const dispatch = useDispatch<AppDispatch>();
  const [isProfileMenuOpen, setIsProfileMenuOpen] = useState(false);
  const { token } = useSelector((state: RootState) => state.auth);

  const handleLogout = () => {
    dispatch(logout());
    dispatch(showToast({ message: 'Logged out successfully', type: 'success' }));
    navigate('/');
  };

  const toggleProfileMenu = () => {
    setIsProfileMenuOpen(!isProfileMenuOpen);
  };

  return (
    <header className="bg-white shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <Link to="/" className="text-xl font-bold text-indigo-600">
              Trikona
            </Link>
          </div>

          {token ? (
            <div className="relative">
              <button
                onClick={toggleProfileMenu}
                className="flex items-center text-gray-700 hover:text-indigo-600 focus:outline-none"
              >
                <UserCircleIcon className="h-8 w-8" />
              </button>

              {isProfileMenuOpen && (
                <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-10">
                  <Link
                    to="/profile"
                    className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    onClick={() => setIsProfileMenuOpen(false)}
                  >
                    Your Profile
                  </Link>
                  <button
                    onClick={() => {
                      handleLogout();
                      setIsProfileMenuOpen(false);
                    }}
                    className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    Sign out
                  </button>
                </div>
              )}
            </div>
          ) : (
            <div className="flex items-center space-x-4">
              <Link
                to="/"
                className="text-gray-700 hover:text-indigo-600"
              >
                Login
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header; 