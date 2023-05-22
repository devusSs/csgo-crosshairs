import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { getCrosshairPosts } from '../api/posts'


function SavedCrosshairs() {
  const { status, error, data:code,} = useQuery({
    queryKey: ['code'], 
    queryFn: getCrosshairPosts
  })

  if (status === "loading") {
    return <span>Loading...</span>
  }
  if (status === "error") {
    return <span>Error: {JSON.stringify(error)}</span>
  }
  

  return (
    <div>
    <div className="container mx-auto">
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
        <tr className="px-6 py-3 whitespace-nowrap">
          <td>{}</td>
          <td>Test</td>
          <td>Test</td>
          <td>Test</td>
        </tr>
      </tbody>
    </table>
  </div>
    </div>
  )
}

export default SavedCrosshairs