import React from 'react'
import ReactDOM from 'react-dom/client'
import App from "./App";
import "./index.css";
import { BrowserRouter } from 'react-router-dom';
import {QueryClient , QueryClientProvider} from '@tanstack/react-query'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <BrowserRouter>
      <QueryClientProvider client={new QueryClient()}>
        <App />
      </QueryClientProvider>
    </BrowserRouter>
  </React.StrictMode>,
)
