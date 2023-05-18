import React from 'react'
import { Link , useMatch , useResolvedPath} from 'react-router-dom'

function Navbar() {
  return (
    <div>
      <nav className='m-0 bg-sky-950 text-white flex justify-between items-center gap-2 p-1'>
        <Link to="/"></Link>
        <ul className='p-0 m-0 list-none flex gap-3 items-stretch'>
          <li className=' hover:bg-slate-700 active:bg-slate-500'>
            <Link to="home">Home</Link>
          </li>
          <li className=' hover:bg-slate-700 active:bg-slate-500'>
            <Link to="generator">Crosshair Generator</Link>
          </li>
          <li className=' hover:bg-slate-700 active:bg-slate-500'>
            <Link to="demo">Demo Extractor</Link>
          </li>
          <li className=' hover:bg-slate-700 active:bg-slate-500'>
            <Link to="login">Login</Link>
          </li>
        </ul>
      </nav>
    </div>
  )
}

function customLink (to: string, children: string, ...props) {
  const resolvedPath = useResolvedPath(to)
  const isActive = useMatch({path: resolvedPath.pathname, end: true})

  return (
    <li className={isActive ? "active" : ""}>
      <Link to={to} {...props}>
        {children}
      </Link>
    </li>
  )
}

export default Navbar



