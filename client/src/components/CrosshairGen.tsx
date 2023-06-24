import React from 'react'
import { useEffect } from 'react'
import useAuth  from '../hooks/useAuth'
import { useNavigate } from 'react-router-dom'
import {encode} from '../data/encode'
import {ChangeEventHandler, useState, ChangeEvent} from 'react'
import {RgbaColor, RgbaColorPicker} from 'react-colorful'
import {getHexColor, getRgbaCSS, getRgbaObject} from '../utils/colors'
import Switch from 'react-switch'


function CrosshairGen() {
  const [style, setStyle] = useState('4');
  
  const isChecked= (value: string) => value === style;
  
  const onSelect = (event: ChangeEvent<HTMLInputElement>) => {
    setStyle(event.currentTarget.value);
  }
  
  const [gap, setGap] = useState('0') 
  const [size, setSize] = useState('1') 
  
  const [thickness, setThickness] = useState('1') 
  
  const [outline, setOutline] = useState('0')
  const toggleOutline = () => setOutline(outline === '0' ? '1' : '0') 
  
  const [Dot, setDot] = useState('0')
  const toggleDot = () => setDot(Dot === '0' ? '1' : '0')
  
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

 const navigate = useNavigate();
  const { auth }: any = useAuth();

  const crosshair = {
    cl_crosshairgap: parseInt(gap),
    cl_crosshair_outlinethickness: 1,
    cl_crosshaircolor_r: selectedColor.r,
    cl_crosshaircolor_g: selectedColor.g,
    cl_crosshaircolor_b: selectedColor.b,
    cl_crosshairalpha: Math.round(selectedColor.a),
    cl_crosshair_dynamic_splitdist: 7,
    cl_fixedcrosshairgap: 3,
    cl_crosshaircolor: 5,
    cl_crosshair_drawoutline: parseInt(outline),
    cl_crosshair_dynamic_splitalpha_innermod: 1,
    cl_crosshair_dynamic_splitalpha_outermod: 0,
    cl_crosshair_dynamic_maxdist_splitratio: 0,
    cl_crosshairthickness: parseInt(thickness),
    cl_crosshairstyle: parseInt(style),
    cl_crosshairdot: parseInt(Dot),
    cl_crosshairgap_useweaponvalue: 0,
    cl_crosshairusealpha: 1,
    cl_crosshair_t: 0,
    cl_crosshairsize: parseInt(size),
  } 

 useEffect((): any => {
    let isMounted = true;
    if (auth?.role !== 'user' && auth?.role !== 'admin') {
      navigate('/users/login')
    }
    return () => (isMounted = false);
 }, [auth?.role])

/* HTML render */ 

  return (
    <div>
      <div className='flex items-center justify-center h-screen'>

{/* Crosshairstyle component */}

        <fieldset className='w-44'>
          <legend>Style:</legend>
            <div className='field'>
              <label> 
                <input type="radio" name="style" value="0" checked={isChecked('0')} onChange={onSelect}/>
                <span>Default</span>
              </label>
            </div>
            <div className='field'>
              <label> 
                <input type="radio" name="style" value="1" checked={isChecked('1')} onChange={onSelect}/>
                <span>Default Static</span>
              </label>
            </div>
            <div className='field'>
              <label> 
                <input type="radio" name="style" value="2" checked={isChecked('2')} onChange={onSelect}/>
                <span>Classic</span>
              </label>
            </div>
            <div className='field'>
              <label> 
                <input type="radio" name="style" value="3" checked={isChecked('3')} onChange={onSelect}/>
                <span>Classic Dynamic</span>
              </label>
            </div>
            <div className='field'>
              <label> 
                <input type="radio" name="style" value="4" checked={isChecked('4')} onChange={onSelect}/>
                <span>Classic Static</span>
              </label>
            </div>
            <div className='field'>
              <label> 
                <input type="radio" name="style" value="5" checked={isChecked('5')} onChange={onSelect}/>
                <span>Classic Hybrid</span>
              </label>
            </div>
          </fieldset>

{/* Size / Gap / Thickness / Outline / Dot settings */}

        <div>
          <ul>
            <li>
              <label htmlFor="ch_size"> Size:
                <br />
                <input type="range" name="ch_size" min={1} max={5} step={0.5} value={size} onChange={(e)=>setSize(e.target.value)}/>
                <input className='w-10' type="number" min={1} max={5} step={0.5} value={size} onChange={(e)=>setSize(e.target.value)}/>
              </label>
            </li>

            <li>
              <label htmlFor="ch_gap"> Gap:
                <br />
                <input type="range" name="ch_gap" min={-5} max={5} step={1} value={gap} onChange={(e)=>setGap(e.target.value)}/>
                <input className='w-10' type="number" min={0} max={3} step={1} value={gap} onChange={(e)=>setGap(e.target.value)}/>
              </label>
            </li>

            <li>
                  <label> Thickness:
                      <br />
                      <input type="range" name="ch_thickness" min={0} max={3} step={0.5} value={thickness} onChange={(e)=>setThickness(e.target.value)}/>
                      <input className='w-10' type="number" min={0} max={3} step={0.5} value={thickness} onChange={(e)=>setThickness(e.target.value)}/>
                  </label>
            </li>

            <li>
                  <p>Outline:</p>
                  <Switch onChange={toggleOutline} checked={outline === '1'} />
            </li>

            <li>
                  <p>Dot:</p>
                  <Switch onChange={toggleDot} checked={Dot === '1'} />
            </li>
          </ul>

        </div>
            


{/* Colormenu */}

      <div onBlur={!changing ? hideColorPicker : undefined} className='flex items-center justify-center w-56'>
        <div className='flex relative'>
          <button onClick={() => setShowColorPicker(!showColorPicker)} className='w-10 h-10 rounded-sm border-2 relative'>
            <div className='absolute -inset-0 z-10 bg-red-500 opacity-100' style={{background:getRgbaCSS(selectedColor)}}/>
            <div className='absolute -inset-0 z-0'>
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

{/* Generate and Save option for user */}

        <ul>
          <li>
            <label htmlFor="sharecode">Sharecode:</label> <br/>
            <input className="mb-6 bg-gray-100 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:focus:ring-blue-500 dark:focus:border-blue-500" type="text" id="sharecode" value={encode(crosshair)} readOnly/>
          </li>
          <li>
            <button className='border-2 mt-5 w-20 bg-green-400 opacity-100 rounded-lg'>Save</button>
          </li>
        </ul>
      </div>
    </div>
  )
}

export default CrosshairGen
