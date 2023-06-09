import React from 'react'
import useAuth from './useAuth'
import { useNavigate } from 'react-router-dom'
import { logoutUser } from '../api/requests'
import { AxiosError } from 'axios'
import { errorResponse , userSuccessResponse  } from '../api/types'

const useLogout = () => {
    const {auth, setAuth}: any = useAuth()
    const navigate = useNavigate();
    
    const logout = async () => {
        if (!auth?.role) {
            navigate('/users/login')
        }

        const response: any = await logoutUser()
        
        if (response instanceof AxiosError) {
            const errResponse = response?.response?.data as errorResponse;
            // TODO: toast error 
        } else {
            const sucResponse = response?.data as userSuccessResponse;
            localStorage.removeItem('role')
            setAuth({})
            navigate('/home')
        }
    }
    return logout
}

export default useLogout