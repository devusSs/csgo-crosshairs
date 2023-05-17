import React , {ChangeEventHandler, useState} from 'react'
import {RgbaColor, RgbaColorPicker} from 'react-colorful'
import {getHexColor, getRgbaCSS, getRgbaObject} from '../utils/colors'

interface Props {}

function colorMenu(props: Props) {
  
  const defaultColor: RgbaColor = {r: 255, g: 255, b: 255, a: 1}

  const [selectedColor, setSelectedColor] = useState<RgbaColor>({...defaultColor})

  const [colorInput, setColorInput] = useState(getHexColor(selectedColor))

  const [changing, setChanging] = useState(false) 
  const [showColorPicker, setShowColorPicker] = useState(false)

  const handleOnColorChange = (value: RgbaColor) => {
    setColorInput(getHexColor(value))
    setSelectedColor(value)
  }

  const handleOnColorInputChange: ChangeEventHandler<HTMLInputElement> =({target}) => { 
    try{
      const {value} = target
      setColorInput(value)
      const rgbaColor = getRgbaObject(value);
    }
    catch(error){}
  }

  const hideColorPicker = () => setShowColorPicker(false)


  return (
    <div onBlur={!changing ? hideColorPicker : undefined} className='flex items-center justify-center w-56'>
      <div className='flex relative'>
        <button onClick={() => setShowColorPicker(!showColorPicker)} className='w-10 h-10 rounded-sm border-2 relative'>
          <div className='absolute -inset-0 z-10 bg-red-500 opacity-100' style={{background:getRgbaCSS(selectedColor)}}/>

          <div className='absolute -inset-0 z-0'>
            <img src="./empty.png" alt="" className='w-full h-full' />
          </div>
        </button>

        <div className='ml-5'>
          <p>Color</p>
          <input type="text" 
          className='bg-transparent border-b w-24 outline-none' 
          value={colorInput} onChange={handleOnColorInputChange} 
          onBlur={() => setColorInput(getHexColor(selectedColor))}
          />
        </div>

        {showColorPicker && <div className='absolute top-full left-0 mt-5'>
          <RgbaColorPicker 
          color={selectedColor} 
          onChange={handleOnColorChange} 
          onMouseDownCapture={() => setChanging(true)} 
          onMouseUp={() => setChanging(false)}
          onMouseLeave={() => setChanging(false)}
          /> 
        </div>}
      </div>
    </div>
  )
}

export default colorMenu