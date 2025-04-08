// components/AuthForm.tsx
import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useLoginMutation, useRegisterMutation, useRequestVerificationLinkMutation } from '../../redux/services/authApi';
import { AppDispatch, RootState } from '../../redux/store';
import { setCredentials } from '../../redux/slices/authSlice';
import { showToast } from '../../redux/slices/toastSlice';
import { Link, useNavigate } from 'react-router-dom';
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';
import { FaLinkedin } from 'react-icons/fa';

interface FormData {
  firstName?: string;
  lastName?: string;
  email: string;
  password: string;
  confirmPassword?: string;
  linkedinUrl?: string;
}

const AuthForm: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const { token } = useSelector((state: RootState) => state.auth);
  const [isLogin, setIsLogin] = useState<boolean>(true);
  const [isInitialRegistration, setIsInitialRegistration] = useState<boolean>(false);
  const [isEmailSent, setIsEmailSent] = useState<boolean>(false);
  const [email, setEmail] = useState<string>('');
  const [formData, setFormData] = useState<FormData>({
    firstName: '',
    lastName: '',
    email: '',
    password: '',
    confirmPassword: '',
    linkedinUrl: '',
  });
  const [showPassword, setShowPassword] = useState<boolean>(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState<boolean>(false);
  const [isLinkedInLoading, setIsLinkedInLoading] = useState<boolean>(false);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

  const [register, { isLoading: isRegistering }] = useRegisterMutation();
  const [login, { isLoading: isLoggingIn }] = useLoginMutation();
  const [requestVerificationLink] = useRequestVerificationLinkMutation();

  // Redirect to dashboard if user is already logged in
  useEffect(() => {
    if (token) {
      navigate('/');
    }
  }, [token, navigate]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (isInitialRegistration) {
      setEmail(e.target.value);
    } else {
      setFormData({ ...formData, [e.target.name]: e.target.value });
    }
  };

  const isStrongPassword = (password: string): boolean => {
    const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[\W_]).{8,}$/;
    return regex.test(password);
  };

  const handleInitialRegistrationSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    
    if (!email) {
      dispatch(showToast({ message: 'Please enter your email address', type: 'error' }));
      return;
    }
    
    setIsSubmitting(true);
    try {
      await requestVerificationLink({ email }).unwrap();
      setIsEmailSent(true);
      dispatch(showToast({ 
        message: 'Verification link sent! Please check your email.', 
        type: 'success' 
      }));
    } catch (err: any) {
      dispatch(showToast({ 
        message: err?.data?.error || 'Failed to send verification link', 
        type: 'error' 
      }));
      setIsSubmitting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!isLogin) {
      if (!isStrongPassword(formData.password)) {
        dispatch(
          showToast({
            message:
              'Password must be at least 8 characters and include uppercase, lowercase, a digit, and a special character.',
            type: 'error',
          })
        );
        return;
      }
      if (formData.password !== formData.confirmPassword) {
        dispatch(showToast({ message: 'Passwords do not match.', type: 'error' }));
        return;
      }
      try {
        await register(formData as any).unwrap();
        dispatch(showToast({ 
          message: 'Registration successful! Please check your email for verification link.', 
          type: 'success' 
        }));
        setIsLogin(true);
      } catch (err: any) {
        dispatch(showToast({ message: err?.data?.error || 'Registration failed', type: 'error' }));
      }
    } else {
      try {
        const result = await login({
          email: formData.email,
          password: formData.password,
        }).unwrap();
        
        // Only set credentials if we have both token and user
        if (result.token && result.user) {
          console.log('Result : ', result);
          dispatch(setCredentials({ token: result.token, user: result.user }));
          dispatch(showToast({ message: 'Login successful!', type: 'success' }));
          // Redirect to the landing page (dashboard)
          navigate('/');
        } else {
          dispatch(showToast({ message: 'Login failed: Invalid response from server', type: 'error' }));
        }
      } catch (err: any) {
        dispatch(showToast({ message: err?.data?.error || 'Login failed', type: 'error' }));
      }
    }
  };

  const toggleMode = () => {
    if (isLogin) {
      setIsInitialRegistration(true);
      setIsLogin(false);
    } else {
      setIsInitialRegistration(false);
      setIsLogin(true);
      setIsEmailSent(false);
      setEmail('');
    }
  };

  const handleLinkedInLogin = async () => {
    setIsLinkedInLoading(true);
    try {
      // LinkedIn OAuth configuration
      const clientId = import.meta.env.VITE_LINKEDIN_CLIENT_ID;
      const redirectUri = import.meta.env.VITE_LINKEDIN_REDIRECT_URI_DEV || 'http://localhost:5173/auth/linkedin/callback';
      const scope = 'r_liteprofile r_emailaddress';
      const state = Math.random().toString(36).substring(7);

      // Store state in localStorage to verify when we return
      localStorage.setItem('linkedin_state', state);

      // Construct LinkedIn OAuth URL
      const params = new URLSearchParams({
        response_type: 'code',
        client_id: clientId,
        redirect_uri: redirectUri,
        scope: scope,
        state: state,
      });

      const authUrl = `https://www.linkedin.com/oauth/v2/authorization?${params.toString()}`;
      
      // Redirect to LinkedIn
      window.location.href = authUrl;
    } catch (error) {
      dispatch(showToast({ 
        message: 'Failed to initialize LinkedIn login', 
        type: 'error' 
      }));
      setIsLinkedInLoading(false);
    }
  };

  // Render initial registration form
  if (isInitialRegistration) {
    return (
      <div className="w-full max-w-md bg-white rounded-xl shadow-lg p-8">
        <h2 className="text-2xl font-semibold text-gray-800 text-center mb-6">
          Start Your Registration
        </h2>
        
        {!isEmailSent ? (
          <>
            <p className="text-gray-600 text-center mb-6">
              Enter your email address to begin the registration process. We'll send you a verification link to continue.
            </p>
            
            <form onSubmit={handleInitialRegistrationSubmit}>
              <div className="mb-6">
                <label htmlFor="email" className="block text-gray-700 mb-1">
                  Email Address
                </label>
                <input
                  id="email"
                  type="email"
                  value={email}
                  onChange={handleChange}
                  required
                  className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
                  placeholder="your.email@example.com"
                />
              </div>
              
              <button
                type="submit"
                disabled={isSubmitting}
                className={`w-full py-2 px-4 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition-colors ${
                  isSubmitting ? 'opacity-70 cursor-not-allowed' : ''
                }`}
              >
                {isSubmitting ? 'Sending...' : 'Send Verification Link'}
              </button>
            </form>
          </>
        ) : (
          <div className="text-center">
            <div className="mb-4 text-green-500">
              <svg className="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h3 className="text-xl font-medium text-gray-900 mb-2">Check Your Email</h3>
            <p className="text-gray-600 mb-6">
              We've sent a verification link to <span className="font-medium">{email}</span>. 
              Please click the link to continue with your registration.
            </p>
            <p className="text-sm text-gray-500 mb-4">
              Didn't receive the email? Check your spam folder or try again.
            </p>
            <button
              onClick={() => setIsEmailSent(false)}
              className="text-indigo-600 hover:underline focus:outline-none"
            >
              Try a different email
            </button>
          </div>
        )}
        
        <div className="mt-6 text-center">
          <p className="text-sm text-gray-600">
            Already have an account?{' '}
            <button onClick={() => setIsLogin(true)} className="text-indigo-600 hover:underline">
              Login
            </button>
          </p>
        </div>
      </div>
    );
  }

  // Render login/registration form
  return (
    <div className="w-full max-w-md bg-white rounded-xl shadow-lg p-8">
      <h2 className="text-2xl font-semibold text-gray-800 text-center mb-6">
        {isLogin ? 'Login' : 'Register'}
      </h2>
      <form onSubmit={handleSubmit}>
        {!isLogin && (
          <>
            <div className="mb-4">
              <label htmlFor="firstName" className="block text-gray-700 mb-1">
                First Name
              </label>
              <input
                id="firstName"
                type="text"
                name="firstName"
                value={formData.firstName}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
              />
            </div>
            <div className="mb-4">
              <label htmlFor="lastName" className="block text-gray-700 mb-1">
                Last Name
              </label>
              <input
                id="lastName"
                type="text"
                name="lastName"
                value={formData.lastName}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
              />
            </div>
            <div className="mb-4">
              <label htmlFor="linkedinUrl" className="block text-gray-700 mb-1">
                LinkedIn Profile URL
              </label>
              <input
                id="linkedinUrl"
                type="url"
                name="linkedinUrl"
                value={formData.linkedinUrl}
                onChange={handleChange}
                required
                placeholder="https://www.linkedin.com/in/your-profile"
                className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
              />
              <p className="mt-1 text-sm text-gray-500">
                This will be used to verify your identity
              </p>
            </div>
          </>
        )}
        <div className="mb-4">
          <label htmlFor="email" className="block text-gray-700 mb-1">
            Email
          </label>
          <input
            id="email"
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
          />
        </div>
        <div className="relative mb-4">
          <label htmlFor="password" className="block text-gray-700 mb-1">
            Password
          </label>
          <div className="relative">
            <input
              id="password"
              type={showPassword ? 'text' : 'password'}
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
              className="w-full px-3 py-2 pr-10 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
            />
            <button
              type="button"
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
              onClick={() => setShowPassword(!showPassword)}
            >
              {showPassword ? (
                <EyeSlashIcon className="h-5 w-5" />
              ) : (
                <EyeIcon className="h-5 w-5" />
              )}
            </button>
          </div>
        </div>
        {!isLogin && (
          <div className="relative mb-6">
            <label htmlFor="confirmPassword" className="block text-gray-700 mb-1">
              Confirm Password
            </label>
            <div className="relative">
              <input
                id="confirmPassword"
                type={showConfirmPassword ? 'text' : 'password'}
                name="confirmPassword"
                value={formData.confirmPassword || ''}
                onChange={handleChange}
                required
                className="w-full px-3 py-2 pr-10 border border-gray-300 rounded focus:outline-none focus:ring focus:border-blue-300"
              />
              <button
                type="button"
                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              >
                {showConfirmPassword ? (
                  <EyeSlashIcon className="h-5 w-5" />
                ) : (
                  <EyeIcon className="h-5 w-5" />
                )}
              </button>
            </div>
          </div>
        )}
        <button
          type="submit"
          className="w-full py-2 px-4 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition-colors"
        >
          {isLogin
            ? isLoggingIn
              ? 'Logging in...'
              : 'Login'
            : isRegistering
            ? 'Registering...'
            : 'Register'}
        </button>
      </form>
      
      {isLogin && (
        <div className="mt-6">
          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-300"></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-2 bg-white text-gray-500">Or continue with</span>
            </div>
          </div>
          
          <div className="mt-6">
            <button
              type="button"
              className="w-full flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
              onClick={handleLinkedInLogin}
              disabled={isLinkedInLoading}
            >
              <FaLinkedin className="h-5 w-5 text-[#0077B5] mr-2" />
              <span>{isLinkedInLoading ? 'Connecting...' : 'Sign in with LinkedIn'}</span>
            </button>
          </div>
        </div>
      )}
      
      <div className="mt-4 text-center">
        {isLogin && (
          <div className="mb-2">
            <Link to="/forgot-password" className="text-sm text-blue-600 hover:underline">
              Forgot Password?
            </Link>
          </div>
        )}
        <button onClick={toggleMode} className="text-indigo-600 hover:underline focus:outline-none">
          {isLogin ? "Don't have an account? Register" : 'Already have an account? Login'}
        </button>
      </div>
    </div>
  );
};

export default AuthForm;
