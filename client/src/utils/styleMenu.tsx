import React , {useState, ChangeEvent}  from 'react'


function styleMenu() {
  const [checked, setChecked] = useState('style_4');

  const isChecked= (value: string) => value === checked;
  
  const onSelect = (event: ChangeEvent<HTMLInputElement>) => {
    setChecked(event.currentTarget.value);
  }

  return (

    <fieldset className='w-44'>
      <legend>Style:</legend>
        <div className='field'>
          <label> 
            <input type="radio" name="style" value="style_0" checked={isChecked('style_0')} onChange={onSelect}/>
            <span>Default</span>
          </label>
        </div>
        <div className='field'>
          <label> 
            <input type="radio" name="style" value="style_1" checked={isChecked('style_1')} onChange={onSelect}/>
            <span>Default Static</span>
          </label>
        </div>
        <div className='field'>
          <label> 
            <input type="radio" name="style" value="style_2" checked={isChecked('style_2')} onChange={onSelect}/>
            <span>Classic</span>
          </label>
        </div>
        <div className='field'>
          <label> 
            <input type="radio" name="style" value="style_3" checked={isChecked('style_3')} onChange={onSelect}/>
            <span>Classic Dynamic</span>
          </label>
        </div>
        <div className='field'>
          <label> 
            <input type="radio" name="style" value="style_4" checked={isChecked('style_4')} onChange={onSelect}/>
            <span>Classic Static</span>
          </label>
        </div>
        <div className='field'>
          <label> 
            <input type="radio" name="style" value="style_5" checked={isChecked('style_5')} onChange={onSelect}/>
            <span>Classic Hybrid</span>
          </label>
        </div>
    </fieldset>
  )
}

export default styleMenu
