// Configuración de Style Dictionary 5.6.1 para C2: lee C2-tokens.json y emite tokens.css y tailwind.theme.css en docs/contratos/.
// Uso: npx style-dictionary build --config config.mjs   (o: npm run build)
import { fileURLToPath } from 'node:url';
import { formatoCss, formatoTailwind } from './formatos.mjs';

const aqui = (ruta) => fileURLToPath(new URL(ruta, import.meta.url));

export default {
  source: [aqui('../../../contratos/C2-tokens.json')],
  hooks: {
    formats: { 'finanzas/css': formatoCss, 'finanzas/tailwind': formatoTailwind },
    // El nombre completo es la ruta: sin esto Style Dictionary usa solo el último segmento y marca falsas colisiones (light, dark, 1...).
    transforms: { 'finanzas/nombre': { type: 'name', transform: (t) => t.path.join('.') } },
  },
  log: { warnings: 'error', verbosity: 'default', errors: { brokenReferences: 'throw' } },
  platforms: {
    web: {
      transforms: ['finanzas/nombre'],
      buildPath: aqui('../../../contratos/') ,
      files: [
        { destination: 'tokens.css', format: 'finanzas/css' },
        { destination: 'tailwind.theme.css', format: 'finanzas/tailwind' },
      ],
    },
  },
};
