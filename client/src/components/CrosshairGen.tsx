import React from 'react'
import ColorMenu from '../utils/colorMenu'
import CrosshairStylingMenu from '../utils/crosshairStylingMenu'
import StyleMenu from '../utils/styleMenu'

function CrosshairGen() {
  return (
    <div>
        <div className='flex items-center justify-center h-screen'>
        <StyleMenu/>
        <ColorMenu/>
        <CrosshairStylingMenu/>
      </div>
    </div>
  )
}

export default CrosshairGen