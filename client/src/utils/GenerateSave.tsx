import React from 'react'
import { encode } from '../data/encode'

function generateSave() {
  const crosshair = {
    cl_crosshairgap: 1,
    cl_crosshair_outlinethickness: 1,
    cl_crosshaircolor_r: 1,
    cl_crosshaircolor_g: 1,
    cl_crosshaircolor_b: 1,
    cl_crosshairalpha: 1,
    cl_crosshair_dynamic_splitdist: 1,
    cl_fixedcrosshairgap: 1,
    cl_crosshaircolor: 1,
    cl_crosshair_drawoutline: true,
    cl_crosshair_dynamic_splitalpha_innermod: 1,
    cl_crosshair_dynamic_splitalpha_outermod: 1,
    cl_crosshair_dynamic_maxdist_splitratio: 1,
    cl_crosshairthickness: 1,
    cl_crosshairstyle: 1,
    cl_crosshairdot: true,
    cl_crosshairgap_useweaponvalue: true,
    cl_crosshairusealpha: true,
    cl_crosshair_t: true,
    cl_crosshairsize: 4,
  }
  
  return (
    <div>
      <ul>
        <li>
          <label htmlFor="sharecode">Sharecode:</label> <br/>
          <input className="bg-blue-400 w-auto" type="text" id="sharecode" value={encode(crosshair)} />
        </li>
        <li>
          <button className='border-2 mt-5 w-20 bg-green-400 opacity-100 rounded-lg'>Save</button>
        </li>
      </ul>
    </div>
  )
}

export default generateSave