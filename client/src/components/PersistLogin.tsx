import React, { useEffect , useState } from 'react'
import useAuth  from '../hooks/useAuth'
import { getMe } from '../api/requests'
import { errorResponse, userSuccessResponse } from '../api/types'
import { AxiosError } from 'axios'
import { Outlet } from 'react-router-dom'
import { useNavigate } from 'react-router-dom'

const PersistLogin = () => {
  const [loading, setLoading] = useState(true);
  const { auth , setAuth }: any = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
      
      async function getMyUserInformation(){
        const response = await getMe()
        if (response instanceof AxiosError) {
            const errResponse = response?.response?.data as errorResponse;
        } else {
            const sucResponse = response.data as userSuccessResponse;
            setAuth({'role' : sucResponse.data.role})
            setLoading(false)
        }
    }
    getMyUserInformation()
    !auth?.role ? setLoading(true) : setLoading(false)
  })

  return (
    <>
        {loading ? navigate('/login') : (<Outlet/>)}
    </>
  )
}

export default PersistLogin