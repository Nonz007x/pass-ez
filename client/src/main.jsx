import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import ClearSessionOnRefresh from './components/ClearSessionOnRefresh/ClearSessionOnRefresh.jsx'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <ClearSessionOnRefresh />
    <App />
  </React.StrictMode>,
)
