import React from 'react'
import { resetPassword } from '../api/requests';
import { AxiosError } from 'axios';
import { useState } from 'react';
import { useSearchParams, useParams, useNavigate } from 'react-router-dom';
import { authorizeRequestPasswordChange } from '../api/requests';
import {Email , errorResponse , userSuccessResponse} from '../api/types'


export default function ResetPassword() {
    
    const [email, setEmail] = useState<Email>({
        e_mail: ""
    });
    
    const [searchParams, setSearchParams] = useSearchParams();
   
    const navigate = useNavigate();
    const emailParams = searchParams.get('email')
   
    const resetPasswordCode = searchParams.get('code')
   
    var resetPasswordURL = `https://api.dropawp.com/api/users/resetPass?email=${emailParams}&code=${resetPasswordCode}`
    
    async function authorizePasswordReset(verifyURL: string){
        const response: any= await authorizeRequestPasswordChange(verifyURL);
        
        if (response instanceof AxiosError) {
            const errResponse = response?.response?.data as errorResponse;
            //toast error
            navigate('/users/reset-password')
        } else {
            const sucResponse = response?.data as userSuccessResponse;
            navigate('/users/newpassword')
        }
    }

    if (resetPasswordCode != null && resetPasswordCode != '' && resetPasswordCode != undefined && emailParams != null && emailParams != '' && emailParams != undefined){
        authorizePasswordReset(resetPasswordURL)
    }

    
    const setNewValue = (e_mail: string) => 
        setEmail(prevState => ({ ...prevState, e_mail}))
    
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const response: any = await resetPassword(email);
        setEmail({} as Email)
        if (response instanceof AxiosError) {
            const errResponse = response?.response?.data as errorResponse;

        } else {
            const sucResponse = response?.data as userSuccessResponse; 
        }
    }

    return (
        <div className="relative flex flex-col justify-center min-h-screen overflow-hidden">
            <div className="w-full p-6 m-auto bg-white rounded-md shadow-md lg:max-w-xl">
                <h1 className="text-3xl font-semibold text-center text-purple-700 underline">
                   Forgot Password
                </h1>
                <form className="mt-6" onSubmit={handleSubmit}>
                    <div className="mb-2">
                        <label
                            htmlFor="email"
                            className="block text-sm font-semibold text-gray-800"
                        >
                            Email
                        </label>
                        <input
                            type="email" placeholder='Enter your email' onChange={evt => {setNewValue(evt.target.value)}}
                            value={email.e_mail || ''}
                            className="block w-full px-4 py-2 mt-2 text-purple-700 bg-white border rounded-md focus:border-purple-400 focus:ring-purple-300 focus:outline-none focus:ring focus:ring-opacity-40"
                        />
                    </div>
                    <div className="mt-6">
                        <button className="w-full px-4 py-2 tracking-wide text-white transition-colors duration-200 transform bg-purple-700 rounded-md hover:bg-purple-600 focus:outline-none focus:bg-purple-600">
                            Reset Password
                        </button>
                    </div>
                </form>
    </div>
</div>
);
}
