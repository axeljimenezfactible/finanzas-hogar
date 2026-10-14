import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './estilos/app.css'
import { App } from './app/App'
async function arrancar() {
  if (import.meta.env.MODE === 'mocks') {
    await (await import('./mocks/navegador')).iniciarMocks()
  } else if (import.meta.env.PROD && 'serviceWorker' in navigator) {
    const { registerSW } = await import('virtual:pwa-register')
    registerSW({ immediate: true })
  }
  createRoot(document.getElementById('raiz')!).render(<StrictMode><App /></StrictMode>)
}
void arrancar()
