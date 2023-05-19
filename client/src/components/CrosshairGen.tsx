import React from 'react'
import ColorMenu from '../utils/colorMenu'
import CrosshairStylingMenu from '../utils/crosshairStylingMenu'
import StyleMenu from '../utils/styleMenu'
import GenerateSave from '../utils/GenerateSave'

function CrosshairGen() {
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