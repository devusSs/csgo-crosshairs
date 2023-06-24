import React, { useEffect, useState } from 'react';
import { getMe, getCrosshairsUser } from '../api/requests';
import { AxiosError } from 'axios';
import { errorResponse, crosshairSuccessResponse } from '../api/types';


const UserProfile: React.FC = () => {
  const [email, setEmail] = useState('');
  const [age, setAge] = useState('');
  const [bio, setBio] = useState('');
  const [avatar, setAvatar] = useState('');
  const [countedCrosshairs, setCountedCrosshairs] = useState(0);

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  const handleAgeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setAge(e.target.value);
  };

  const handleBioChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setBio(e.target.value);
  };
 
  
  useEffect(() => {
    async function getCountedCrosshairs() {
      const response: any = await getCrosshairsUser()
      setCountedCrosshairs(response.data.data.crosshairs.length)
    }

    async function getUserAvatar(){
      const response: any = await getMe()

      setAvatar(response.data.data.profile_picture_link)
    }
    getUserAvatar()
    getCountedCrosshairs()
    },[])


  return (
    <div className="max-w-md mx-auto mt-8 p-4 bg-gray-100 rounded shadow">
      <div className="flex items-center mb-4">
        <img
          src={avatar}
          alt="Avatar"
          className="w-10 h-10 rounded-full mr-3"
        />
        <h1 className="text-2xl"> Profile </h1>
      </div>
      <div className="flex items-center mb-4">
        <span className="font-bold text-gray-800">Crosshairs saved:</span>
        <div className="w-10 h-1 flex items-center justify-center mr-3">
          {countedCrosshairs}
        </div>
      </div>
      <form>
        <div className="mb-2">
          <label htmlFor="email" className="block mb-1 font-bold text-gray-800">
            Email:
          </label>
          <input
            type="email"
            id="email"
            value={email}
            onChange={handleEmailChange}
            className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
          />
        </div>
        <div className="mb-2">
          <label htmlFor="age" className="block mb-1 font-bold text-gray-800">
            Age:
          </label>
          <input
            type="text"
            id="age"
            value={age}
            onChange={handleAgeChange}
            className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
          />
        </div>
        <div className="mb-2">
          <label htmlFor="bio" className="block mb-1 font-bold text-gray-800">
            Bio:
          </label>
          <textarea
            id="bio"
            value={bio}
            onChange={handleBioChange}
            className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:border-blue-500"
          ></textarea>
        </div>
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Save
        </button>
      </form>
    </div>
  );
};

export default UserProfile;
