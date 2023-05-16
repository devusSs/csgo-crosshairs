import React , {useState} from 'react'

interface Props {
    gap : number
    setGap : number
}

function gapMenu(props: Props) {
    
    const [gap, setGap] = useState('0')
    
    

    return (
        <div className='flex items-center justify-center'>
            <input type="range" name="ch_gap" min={-5} max={5} step={1} value={gap} onChange={(e)=>setGap(e.target.value)}/>
            <input className='items-end w-10' type="number" min={0} max={3} step={1} value={gap} onChange={(e)=>setGap(e.target.value)}/>
        </div>
    )
}

export default gapMenu