import React , {useState} from 'react'
import Switch from 'react-switch'

interface Props {
}

function thicknessMenu(props: Props) {
    const [gap, setGap] = useState('0') 

    const [thickness, setThickness] = useState('1') 

    const [outline, setOutline] = useState('0')
    const toggleOutline = () => setOutline(outline === '0' ? '1' : '0') 

    const [Dot, setDot] = useState('0')
    const toggleDot = () => setDot(Dot === '0' ? '1' : '0')
  
    return (
        <div className='mx-10'>
            <div>
                <label htmlFor="ch_gap"> Gap:
                    <br />
                    <input type="range" name="ch_gap" min={-5} max={5} step={1} value={gap} onChange={(e)=>setGap(e.target.value)}/>
                    <input className='w-10' type="number" min={0} max={3} step={1} value={gap} onChange={(e)=>setGap(e.target.value)}/>
                </label>
            </div>
            <div>
                <label> Thickness:
                    <br />
                    <input type="range" name="ch_thickness" min={0} max={3} step={0.5} value={thickness} onChange={(e)=>setThickness(e.target.value)}/>
                    <input className='w-10' type="number" min={0} max={3} step={0.5} value={thickness} onChange={(e)=>setThickness(e.target.value)}/>
                </label>
            </div> 
            <div>
                <p>Outline:</p>
                <Switch onChange={toggleOutline} checked={outline === '1'} />
            </div>
            <div>
                <p>Dot:</p>
                <Switch onChange={toggleDot} checked={Dot === '1'} />
            </div>
        </div>
    )
}

export default thicknessMenu