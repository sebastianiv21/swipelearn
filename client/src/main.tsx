import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import QueryProvider from './providers/QueryProvider'
import { Router } from './router/Router'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryProvider>
      <Router />
    </QueryProvider>
  </StrictMode>,
)
