import React , { useState } from 'react'
import { createUser, verifyMail } from '../api/requests'
import { User, errorResponse, userSuccessResponse } from '../api/types'
import { AxiosError } from 'axios';
import { useSearchParams , useNavigate} from 'react-router-dom';



export default function Register() {
    const navigate = useNavigate();
    const [searchParams, setSearchParams] = useSearchParams();
    const verificationCode = searchParams.get('code')
    const BASE_URL = import.meta.env.REACT_APP_BASE_URL;

    var verifyURL = `${BASE_URL}/api/users/verifyMail?code=${verificationCode}`

    async function callVerifyMail(verifyURL: string){
        const response: any= await verifyMail(verifyURL);
        if (response instanceof AxiosError) {
            const errResponse = response?.response?.data as errorResponse;
            navigate('/users/register')
        } else {
            const sucResponse = response?.data as userSuccessResponse;
            navigate('/users/login')
        }
    }

    if (verificationCode != null && verificationCode != '' && verificationCode != undefined){
        callVerifyMail(verifyURL)
    }

    const [user, setUser] = useState<User>({
        e_mail: "",
        password: ""
    });
    
    const setNewValue = (e_mail: string , password: string) => 
        setUser(prevState => ({ ...prevState, [e_mail]: password }))

    
    const handleSubmit = async (e: any) => {
        console.log(user)

        e.preventDefault();
        setUser({} as User)
        
        const response: any = await createUser(user);
        if (response instanceof AxiosError) {
            const errResponse = response?.response?.data as errorResponse;
            // toast error
        } else {
            const sucResponse = response?.data as userSuccessResponse;
            // toast du nutten verify
        }
        console.log(user)
        console.log(response)
    }    

    return (
        <div className="relative flex flex-col justify-center min-h-screen overflow-hidden">
            <div className="w-full p-6 m-auto bg-white rounded-md shadow-md lg:max-w-xl">
                <h1 className="text-3xl font-semibold text-center text-purple-700 underline">
                   Sign Up
                </h1>
                <form onSubmit={handleSubmit} className="mt-6">
                    <div className="mb-2">
                        <label
                            htmlFor="email"
                            className="block text-sm font-semibold text-gray-800"
                        >
                            Email
                        </label>
                        <input
                            type="email" placeholder='Enter your email' onChange={evt => {setNewValue('e_mail', evt.target.value)}}
                            value={user.e_mail || ''}
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
                            value={user.password ||  ''} 
                            className="block w-full px-4 py-2 mt-2 text-purple-700 bg-white border rounded-md focus:border-purple-400 focus:ring-purple-300 focus:outline-none focus:ring focus:ring-opacity-40"
                        />
                    </div>
        
                    <div className="mt-6">
                        <button className="w-full px-4 py-2 tracking-wide text-white transition-colors duration-200 transform bg-purple-700 rounded-md hover:bg-purple-600 focus:outline-none focus:bg-purple-600">
                            Register
                        </button>
                    </div>
                </form>
    </div>
</div>
);
}
