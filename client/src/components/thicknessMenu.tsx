import React , {useState} from 'react'

interface Props {
    thickness : number
    setThickness : number
}

function thicknessMenu(props: Props) {

  const [thickness, setThickness] = useState('1')
    
    return (
        <div className='flex items-center justify-center'>
            <input type="range" name="ch_thickness" min={0} max={3} step={0.5} value={thickness} onChange={(e)=>setThickness(e.target.value)}/>
            <input className='items-end w-10' type="number" min={0} max={3} step={0.5} value={thickness} onChange={(e)=>setThickness(e.target.value)}/>
        </div>
    )
}

export default thicknessMenu