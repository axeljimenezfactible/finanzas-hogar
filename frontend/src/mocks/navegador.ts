import { setupWorker } from 'msw/browser'
import { manejadores } from './manejadores'
export async function iniciarMocks() {
  await setupWorker(...manejadores).start({ serviceWorker: { url: '/mockServiceWorker.js' }, quiet: true })
}
