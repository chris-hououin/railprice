import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import 'leaflet/dist/leaflet.css'
import './index.css'
import App from './App.jsx'
import {HashRouter, Navigate, Route, Routes} from "react-router-dom";

const queryClient = new QueryClient()

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <HashRouter>
        <Routes>
          <Route path="" element={<Navigate to="/orig" />} />
          <Route path="/orig/" element={<Navigate to="/orig/7022" />} />
          <Route path="/orig/:station" element={<App direction="orig" />} />
          <Route path="/dest/" element={<Navigate to="/dest/7022" />} />
          <Route path="/dest/:station" element={<App direction="dest" />} />
        </Routes>
      </HashRouter>
    </QueryClientProvider>
  </StrictMode>,
)
