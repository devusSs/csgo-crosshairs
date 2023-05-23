import React , {useEffect, useState} from 'react'
import { getCrosshairsUser } from '../api/requests'
import { Crosshair } from '../api/types'
import { AxiosError } from 'axios';
import { errorResponse, crosshairSuccessResponse } from '../api/types';
import { useNavigate } from 'react-router-dom';


function SavedCrosshairs() {
  const [crosshairs, setCrosshairs] = useState([] as Crosshair[])
  const [status, setStatus] = useState('Loading...')
  const navigate = useNavigate();

  useEffect(() => {
    async function getCrosshairs(){
      const response = await getCrosshairsUser()

      if (response instanceof AxiosError) {
        const errResponse = response?.response?.data as errorResponse;
        setStatus(errResponse.error.error_message)
        navigate('/login')
        
      } else {
          const sucResponse = response.data as crosshairSuccessResponse;
          setCrosshairs(sucResponse.data.crosshairs)
      }

    }
    getCrosshairs()
  },[])
  
  return (
    <div>
      <div className="container mx-auto">
        {crosshairs.length ? (

        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Crosshair</th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Note</th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
              <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Options</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {crosshairs.map((crosshair) => (
            <tr key={crosshair.code} className="px-6 py-3 whitespace-nowrap">
              <td key={'code'} className="px-6 py-4 whitespace-nowrap">{crosshair.code}</td>
              <td key={'note'} className="px-6 py-4 whitespace-nowrap">{crosshair.note}</td>
              <td key={'added'} className="px-6 py-4 whitespace-nowrap">{crosshair.added.slice(0, 10)}</td>
            </tr>
            ))}
          </tbody>
        </table>
        ) : (<div>{status}</div>)}
      </div>
    </div>
  )
}

export default SavedCrosshairs