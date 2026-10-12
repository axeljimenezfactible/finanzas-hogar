// Storybook 10 de Finanzas Hogar (WF-3a · front-variantes). Muestra cada componente de DESIGN.md §7
// con sus variantes, estados y tamaños, en claro y oscuro. Sin addons, sin red, sin telemetría.
import type { StorybookConfig } from '@storybook/react-vite'
import type { PluginOption } from 'vite'

// Quita recursivamente los plugins de la PWA: Storybook no necesita service worker ni manifest.
const sinPwa = (plugins: PluginOption[] = []): PluginOption[] =>
  plugins.flatMap((p): PluginOption[] => {
    if (Array.isArray(p)) return [sinPwa(p)]
    if (p && typeof p === 'object' && 'name' in p && typeof p.name === 'string' && p.name.startsWith('vite-plugin-pwa')) return []
    return [p]
  })

const config: StorybookConfig = {
  framework: { name: '@storybook/react-vite', options: {} },
  stories: ['../src/componentes/**/*.stories.tsx'],
  addons: [],
  // Inter autoalojada: «npm run tokens» la copia a public/fonts.
  staticDirs: [{ from: '../public/fonts', to: '/fonts' }],
  core: { disableTelemetry: true, disableWhatsNewNotifications: true, enableCrashReports: false },
  typescript: { check: false, reactDocgen: false },
  viteFinal: async (cfg) => ({ ...cfg, plugins: sinPwa(cfg.plugins as PluginOption[]) }),
}
export default config
