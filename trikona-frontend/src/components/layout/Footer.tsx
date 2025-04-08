import React from 'react';

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-gray-50 border-t border-gray-200">
      <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="text-gray-500 text-sm">
            &copy; {currentYear} Trikona. All rights reserved.
          </div>
          <div className="mt-4 md:mt-0">
            <a
              href="#"
              className="text-gray-500 hover:text-indigo-600 text-sm mx-3"
            >
              Privacy Policy
            </a>
            <a
              href="#"
              className="text-gray-500 hover:text-indigo-600 text-sm mx-3"
            >
              Terms of Service
            </a>
            <a
              href="#"
              className="text-gray-500 hover:text-indigo-600 text-sm mx-3"
            >
              Contact Us
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer; 