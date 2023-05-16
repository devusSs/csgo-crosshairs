import StyleMenu from './components/styleMenu.tsx'
import ColorMenu from './components/colorMenu.tsx'
import GapMenu from './components/gapMenu.tsx'
import ThicknessMenu from './components/thicknessMenu.tsx'

const App = () => {

  return (
    <div>  
      <StyleMenu/>
      <ColorMenu/>
      <GapMenu/>
      <ThicknessMenu/>
    </div>
  )
}

export default App
