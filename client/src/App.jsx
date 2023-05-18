import { Route, Routes } from 'react-router-dom'
import CrosshairGen from './components/CrosshairGen'
import Navbar from './components/Navbar'
import Home from './components/Home'
import Demo from './components/Demo'
import Login from './components/Login'
import SignUp from './components/Register'
import ResetPassword from './components/ResetPw'
import SavedCrosshairs from './components/SavedCrosshairs'

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
          <Route path='/reset' element={<ResetPassword/>} />
          <Route path='/saved' element={<SavedCrosshairs/>} />
        </Routes> 
    </div>
  )
}

export default App
