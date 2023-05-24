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

  useEffect((): any => {
    let isMounted = true;

      async function getMyUserInformation(){

        const response: any = await getMe()
        
        if (response instanceof AxiosError) {
            const errResponse = response?.response?.data as errorResponse;
            setAuth({}) 
            localStorage.removeItem('role')
        } else {
            const sucResponse = response?.data as userSuccessResponse;
            setAuth({'role' : sucResponse.data.role})
            localStorage.setItem('role', sucResponse.data.role)
        }
        isMounted && setLoading(false)
    }
    getMyUserInformation()
    !auth?.role ? setLoading(true) : setLoading(false)
    
    return () => (isMounted = false);
  }, [])
  return(
    <>
      {loading ? ( <div>Loading...</div> ) : (<Outlet />)}
    </>
  )
}

export default PersistLogin