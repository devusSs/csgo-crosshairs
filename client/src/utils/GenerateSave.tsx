import React from 'react'

function generateSave() {
  return (
    <div>
      <ul>
        <li>
          <button className='border-2 mt-5 w-20 bg-cyan-400 opacity-100 rounded-lg'>Generate</button>
        </li>
        <li>
          <button className='border-2 mt-5 w-20 bg-green-400 opacity-100 rounded-lg'>Save</button>
        </li>
      </ul>
    </div>
  )
}

export default generateSave