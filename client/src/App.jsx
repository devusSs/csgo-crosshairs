import { Route, Routes } from 'react-router-dom'
import CrosshairGen from './components/CrosshairGen'
import Navbar from './components/Navbar'
import Home from './components/Home'
import Demo from './components/Demo'
import Login from './components/Login'
import SignUp from './components/Register'
import ResetPassword from './components/ResetPassword'
import SavedCrosshairs from './components/SavedCrosshairs'
import PersistLogin from './components/PersistLogin'
import UserProfile from './components/UserProfile'
import Footer from './components/Footer'

const App = () => {
 
  
  return (
    <div className=''>
        <Navbar/> 
        <Routes>
          <Route element={<PersistLogin/>}>
              <Route path='/home' element={<Home/>} />  
              <Route path='/crosshairs/generator' element={<CrosshairGen/>} />
              <Route path='/crosshairs/saved' element={<SavedCrosshairs/>} />
              <Route path='/crosshairs/demo' element={<Demo/>} />
              <Route path='/users/profile' element={<UserProfile/>} />
              <Route path='/users/login' element={<Login/>} />
              <Route path='/users/register' element={<SignUp/>} />
              <Route path='/users/register?code=' element={<SignUp/>} />
              <Route path='/users/reset-password' element={<ResetPassword/>} />
              <Route path='/users/reset-password?email=&code=' element={<ResetPassword/>} />
              <Route path='/users/new-password' element={<ResetPassword/>} />

              <Route path='*' element={<Home/>} />
          </Route>
        </Routes> 
        <Footer/>
    </div>
  )
}

export default App
