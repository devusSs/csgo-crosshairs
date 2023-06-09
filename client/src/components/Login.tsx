import React , {useState} from 'react';
import {User , errorResponse , userSuccessResponse} from '../api/types'
import { Link , useNavigate } from 'react-router-dom';
import {loginUser} from '../api/requests'
import { AxiosError } from 'axios';
import useAuth from '../hooks/useAuth';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css'


export default function Login() {
    const navigate = useNavigate();
    
    const [user, setUser] = useState<User>({
        e_mail: "",
        password: ""
    });
    const {auth, setAuth} : any = useAuth();
    
    const toastSuccessMessage = () => {
        toast.success('Login successful!', { position: "top-right",
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        progress: undefined,
        theme: "colored", });
    }
    const toastErrorMessage = () => {
        toast.error('Login failed!', { position: "top-right",
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
        progress: undefined,
        theme: "colored",});
    }  

    const handleToast = () => {
        if (auth) {
            toastSuccessMessage()
        } else {
            toastErrorMessage()
        }
    }

    const setNewValue = (e_mail: string , password: string) => 
        setUser(prevState => ({ ...prevState, [e_mail]: password }))

        const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const response: any = await loginUser(user);
        if (response instanceof AxiosError) {
            const errResponse = response?.response?.data as errorResponse;
            setUser({} as User)

        } else {
            const sucResponse = response?.data as userSuccessResponse;
            localStorage.setItem('role', sucResponse.data.role)
            setAuth(sucResponse.data)
            navigate('/home')
        }
    }
    

    return (
        <div>
            <ToastContainer 
                position="top-right"
                autoClose={5000}
                hideProgressBar={false}
                newestOnTop={false}
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                draggable
                pauseOnHover
                theme="colored"
            />
        <div className="relative flex flex-col justify-center min-h-screen overflow-hidden">
            <div className="w-full p-6 m-auto bg-white rounded-md shadow-md lg:max-w-xl">
                <h1 className="text-3xl font-semibold text-center text-purple-700 underline">
                   Sign in
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
                            value={user.password || ''}
                            className="block w-full px-4 py-2 mt-2 text-purple-700 bg-white border rounded-md focus:border-purple-400 focus:ring-purple-300 focus:outline-none focus:ring focus:ring-opacity-40"
                        />
                    </div>
                    <Link to="/users/reset-password" className="text-xs text-purple-600 hover:underline">
                            Forgot Password?
                    </Link>
                    <div className="mt-6">
                        <button onClick={handleToast} className="w-full px-4 py-2 tracking-wide text-white transition-colors duration-200 transform bg-purple-700 rounded-md hover:bg-purple-600 focus:outline-none focus:bg-purple-600">
                            Login
                        </button>
                        
                    </div>
                </form>

                <p className="mt-8 text-xs font-light text-center text-gray-700">
                    {" "}
                    Don't have an account?{" "}
                    <Link to="/users/register" className="font-medium text-purple-600 hover:underline">        
                            Sign up
                    </Link>
                </p>
            </div>
        </div>
        </div>
    );
}