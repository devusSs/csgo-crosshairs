import { Route, Routes } from 'react-router-dom'
import CrosshairGen from './components/CrosshairGen.tsx'
import Navbar from './components/Navbar.tsx'
import Home from './components/Home.tsx'
import Demo from './components/Demo.tsx'
import Login from './components/Login.tsx'
import SignUp from './components/Register.tsx'

const App = () => {

  return (
    <div className=''>
        <Navbar/> 
        <Routes>
          <Route path='/home' element={<Home/>} />
          <Route path='/generator' element={<CrosshairGen/>} />
          <Route path='/demo' element={<Demo/>} />
          <Route path='/login' element={<Login/>} />
          <Route path='/register' element={<SignUp/>} />
        </Routes> 
    </div>
  )
}

export default App
