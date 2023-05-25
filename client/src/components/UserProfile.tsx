import React from 'react';

// User Profile Page component
const UserProfile: React.FC = () => {
  return (
    <div className="max-w-md mx-auto bg-white shadow-lg rounded-lg overflow-hidden">
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-2">Manfred Meer</h1>
        <p className="text-gray-600 mb-2">Age: 21</p>
        <p className="text-gray-600 mb-2">Email: test@test.de</p>
        <p className="text-gray-600">Bio: super typ</p>
      </div>
    </div>
  );
};

export default UserProfile;