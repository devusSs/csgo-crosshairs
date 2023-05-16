import StyleMenu from './components/styleMenu.tsx'
import ColorMenu from './components/colorMenu.tsx'
import CrosshairStylingMenu from './components/crosshairStylingMenu.tsx'

const App = () => {

  return (
    <div className='border-2 items-center flex w-1/2 h-1/2'>  
      <StyleMenu/>
      <ColorMenu/>
      <CrosshairStylingMenu/>
    </div>
  )
}

export default App
