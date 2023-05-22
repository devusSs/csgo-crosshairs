import React , { useRef, useState , useEffect } from 'react'
import {createAnUser} from '../api/posts'
import { User } from '../api/types'


export default function Register() {
    const [user, setUser] = useState<User>({
        e_mail: "",
        password: ""
    });
    
    const setNewValue = (e_mail: string , password: string) => 
        setUser(prevState => ({ ...prevState, [e_mail]: password }))
        console.log(user);
  
    
        

    return (
        <div className="relative flex flex-col justify-center min-h-screen overflow-hidden">
            <div className="w-full p-6 m-auto bg-white rounded-md shadow-md lg:max-w-xl">
                <h1 className="text-3xl font-semibold text-center text-purple-700 underline">
                   Sign Up
                </h1>
                <form className="mt-6">
                    <div className="mb-2">
                        <label
                            htmlFor="email"
                            className="block text-sm font-semibold text-gray-800"
                        >
                            Email
                        </label>
                        <input
                            type="email" placeholder='Enter your email' onChange={evt => {setNewValue('e_mail', evt.target.value)}}
                            value={user.e_mail}
                            className="block w-full px-4 py-2 mt-2 text-purple-700 bg-white border rounded-md focus:border-purple-400 focus:ring-purple-300 focus:outline-none focus:ring focus:ring-opacity-40"
                        />
                    </div>
                    <div className="mb-2">
                        <label
                            htmlFor="password" 
                            className="block text-sm font-semibold text-gray-800"
                        >
                            Password
                        </label>
                        <input
                            type="password" placeholder='Enter your password' onChange={evt => {setNewValue('password', evt.target.value)}}
                            value={user.password}
                            className="block w-full px-4 py-2 mt-2 text-purple-700 bg-white border rounded-md focus:border-purple-400 focus:ring-purple-300 focus:outline-none focus:ring focus:ring-opacity-40"
                        />
                    </div>
        
                    <div className="mt-6">
                        <button onClick={() => createAnUser(user)} className="w-full px-4 py-2 tracking-wide text-white transition-colors duration-200 transform bg-purple-700 rounded-md hover:bg-purple-600 focus:outline-none focus:bg-purple-600">
                            Register
                        </button>
                    </div>
                </form>
    </div>
</div>
);
}
