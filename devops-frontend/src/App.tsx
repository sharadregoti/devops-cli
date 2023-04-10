import Home from './pages/Home'
import { ConfigProvider, theme } from 'antd'
import './App.css'

function App() {

  return (
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
      }}
    >
      <Home />
    </ConfigProvider>
  )
}

export default App
