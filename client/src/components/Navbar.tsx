import React, { useState, useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import { BiLogIn , BiLogOut } from 'react-icons/bi';
import useAuth from '../hooks/useAuth';

import useLogout from '../hooks/useLogout';

const Navbar: React.FC = () => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const {auth}: any = useAuth();
  const logout = useLogout();

  const toggleDropdown = () => {
    setIsOpen(!isOpen);
  };

  const callLogout = async() => {
    await logout()
  }

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    window.addEventListener('click', handleClickOutside);

    return () => {
      window.removeEventListener('click', handleClickOutside);
    };
  }, []);

  return (
    <nav className="bg-gray-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center">
            <Link to="/home" className="text-white font-semibold text-lg">
              dropawp.com
            </Link>
          </div>
          <div className="hidden md:block">
            <div className="ml-4 flex items-center">
            
              {auth?.role && (
                <div className="relative" ref={dropdownRef}>
                <button onClick={toggleDropdown} className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium focus:outline-none">
                  Crosshair Generator
                </button>
                {isOpen && (
                    <div className="origin-top-right absolute right-auto mt-2 w-36 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5">
                    <div className="py-1" role="menu" aria-orientation="vertical" aria-labelledby="options-menu">
                      <Link to="/crosshairs/generator" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" role="menuitem">
                        Generator
                      </Link>
                      <Link to="/crosshairs/saved" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" role="menuitem">
                        Saved Crosshairs
                      </Link>
                    </div>
                  </div>
                )}
                </div>
                )}

              {auth?.role && (
                <Link to="/crosshairs/demo" className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">
                Demo Extractor
                </Link>
              )}
              
              {!auth?.role && (
                <Link to="/users/login" className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">
                  <span>
                    Login
                  </span>
                  <BiLogIn className="inline-block ml-1" />
                </Link>
              )}

              {auth?.role && ( 
                // Add dropdown menu for user profile
                <button onClick={callLogout} className="text-gray-300 hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium">
                <span>
                  Logout
                </span>
                <BiLogOut className="inline-block ml-1" />
              </button>
              )}

              
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;