import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useLinkedinAuthMutation } from '../redux/services/authApi';
import { useDispatch } from 'react-redux';
import { setCredentials } from '../redux/slices/authSlice';

const LinkedInCallback: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const dispatch = useDispatch();
  const [linkedinAuth] = useLinkedinAuthMutation();

  useEffect(() => {
    const handleCallback = async () => {
      const code = searchParams.get('code');
      const state = searchParams.get('state');
      const storedState = localStorage.getItem('linkedin_state');

      if (!code || !state || state !== storedState) {
        navigate('/auth');
        return;
      }

      try {
        const result = await linkedinAuth({ code }).unwrap();
        dispatch(setCredentials({ token: result.token || '', user: result.user || null }));
        navigate('/');
      } catch (error) {
        console.error('LinkedIn authentication failed:', error);
        navigate('/auth');
      }
    };

    handleCallback();
  }, [searchParams, navigate, linkedinAuth, dispatch]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center">
        <h2 className="text-xl font-semibold text-gray-900">Processing LinkedIn login...</h2>
        <p className="mt-2 text-gray-600">Please wait while we complete your authentication.</p>
      </div>
    </div>
  );
};

export default LinkedInCallback;