import React from 'react'
import ColorMenu from '../utils/colorMenu'
import CrosshairStylingMenu from '../utils/crosshairStylingMenu'
import StyleMenu from '../utils/styleMenu'
import GenerateSave from '../utils/GenerateSave'
import { useEffect } from 'react'
import useAuth  from '../hooks/useAuth'
import { useNavigate } from 'react-router-dom'

function CrosshairGen() {
 const navigate = useNavigate();
  const { auth }: any = useAuth();
 
 useEffect((): any => {
    let isMounted = true;
    if (auth?.role !== 'user' && auth?.role !== 'admin') {
      navigate('/login')
    }
    return () => (isMounted = false);
 }, [auth?.role])
  return (
    <div>
        <div className='flex items-center justify-center h-screen'>
          <StyleMenu/>
          <ColorMenu/>
          <CrosshairStylingMenu/>
          <GenerateSave />
      </div>
    </div>
  )
}

export default CrosshairGen