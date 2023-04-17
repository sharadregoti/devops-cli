import Home from './pages/Home'
import PluginSelector from './pages/PluginSelector'
import { ConfigProvider, theme } from 'antd'
import './App.css'
import { BrowserRouter as Router, Route, Routes, useParams } from 'react-router-dom';
import { Provider } from 'react-redux';
import store from './redux/store';

function App() {
  return (
    <ConfigProvider
      theme={{
        algorithm: theme.defaultAlgorithm,
      }}
    >
      <Provider store={store}>
        <Router>
          <Routes>
            <Route path="/" element={<PluginSelector />} />
            <Route path="/plugin/:pluginName/session/:sessionId/:authId/:contextId" element={<Home />} />
            {/* <Route path="/project/:projectId" element={<Project />} />
          <Route path="/project" element={<ProjectSelector />} /> */}
          </Routes>
        </Router>
      </Provider>
    </ConfigProvider>
  )
}

export default App
