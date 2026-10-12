// Vista previa de Storybook: mismos estilos que la app (tokens C2 + Tailwind con solo las escalas de C2).
import { useLayoutEffect, type ReactNode } from 'react'
import type { Preview } from '@storybook/react-vite'
import '../src/estilos/app.css'

type Tema = 'sistema' | 'claro' | 'oscuro'

/** Pone `data-theme` en <html>, igual que src/app/tema.ts: «sistema» lo quita y manda prefers-color-scheme. */
function ContenedorTema({ tema, children }: { tema: Tema; children: ReactNode }) {
  useLayoutEffect(() => {
    const raiz = document.documentElement
    if (tema === 'sistema') raiz.removeAttribute('data-theme')
    else raiz.setAttribute('data-theme', tema === 'claro' ? 'light' : 'dark')
  }, [tema])
  return <>{children}</>
}

const preview: Preview = {
  initialGlobals: { tema: 'sistema' },
  globalTypes: {
    tema: {
      description: 'Tema de color (DESIGN.md §2)',
      toolbar: {
        title: 'Tema',
        icon: 'paintbrush',
        dynamicTitle: true,
        items: [
          { value: 'sistema', title: 'Según el sistema' },
          { value: 'claro', title: 'Claro' },
          { value: 'oscuro', title: 'Oscuro' },
        ],
      },
    },
  },
  decorators: [
    (Historia, ctx) => (
      <ContenedorTema tema={(ctx.globals.tema as Tema) ?? 'sistema'}>
        <Historia />
      </ContenedorTema>
    ),
  ],
  parameters: {
    layout: 'padded',
    // El fondo lo pone el tema (ui.surface.page); el selector de fondos de Storybook no se usa.
    backgrounds: { disabled: true },
    controls: { expanded: false, sort: 'requiredFirst' },
    options: { storySort: { order: ['Componentes'] } },
    // Puntos de quiebre de DESIGN.md §9.4 (360, 412, 768, 1024, 1440).
    viewport: {
      options: {
        m360: { name: 'Móvil 360', styles: { width: '360px', height: '800px' }, type: 'mobile' },
        m412: { name: 'Móvil 412', styles: { width: '412px', height: '915px' }, type: 'mobile' },
        t768: { name: 'Tableta 768', styles: { width: '768px', height: '1024px' }, type: 'tablet' },
        e1024: { name: 'Escritorio 1024', styles: { width: '1024px', height: '768px' }, type: 'desktop' },
        e1440: { name: 'Escritorio 1440', styles: { width: '1440px', height: '900px' }, type: 'desktop' },
      },
    },
  },
}
export default preview
