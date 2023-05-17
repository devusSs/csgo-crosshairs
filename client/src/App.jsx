import StyleMenu from './components/styleMenu.tsx'
import ColorMenu from './components/colorMenu.tsx'
import CrosshairStylingMenu from './components/crosshairStylingMenu.tsx'
import Navbar from './components/Navbar.tsx'

const App = () => {

  return (
    <div className=''>
      <Navbar/>  
      <div className='flex items-center justify-center h-screen'>
        <StyleMenu/>
        <ColorMenu/>
        <CrosshairStylingMenu/>
      </div>
    </div>
  )
}

export default App
