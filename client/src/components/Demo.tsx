import React from 'react'
import { useEffect } from 'react'
import useAuth  from '../hooks/useAuth'
import { useNavigate } from 'react-router-dom'


function Demo() {
  const navigate = useNavigate();
  const { auth }: any = useAuth();
 
 useEffect((): any => {
    let isMounted = true;
    if (auth?.role !== 'user' && auth?.role !== 'admin') {
      navigate('/users/login')
    }
    return () => (isMounted = false);
 }, [auth?.role])
  return (
    <div className='text-xl'>Demo</div>
  )
}

export default Demo