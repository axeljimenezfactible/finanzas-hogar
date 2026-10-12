// Formatos propios de Style Dictionary 5.6.1 para C2 (tokens.css y tailwind.theme.css).
// Nada de valores escritos aqui: todo sale de C2-tokens.json. Salida determinista (sin fechas).

const EXT = 'finanzas';
const REM = (n) => `${Number((n / 16).toFixed(4))}rem`;
const GRUPOS_REM = ['font.size', 'font.line-height', 'space'];

// ---------- utilidades ----------
const limpia = (path) => (['@', '$root'].includes(path.at(-1)) ? path.slice(0, -1) : path);
const esTema = (t) => ['ui', 'overlay'].includes(t.path[0]) && ['light', 'dark'].includes(t.path.at(-1));
const ext = (t) => t.$extensions?.[EXT] ?? {};

export function cssName(path) {
  const p = limpia(path);
  const [a, b] = p;
  if (a === 'font' && b === 'family') return ['font', ...p.slice(2)].join('-');
  if (a === 'font' && b === 'line-height') return p.slice(1).join('-').replace(/^/, 'line-height-').replace('line-height-line-height-', 'line-height-');
  if (a === 'motion' && b === 'duration') return ['duration', ...p.slice(2)].join('-');
  if (a === 'motion' && b === 'easing') return ['ease', ...p.slice(2)].join('-');
  return p.join('-');
}
const themeName = (path) => {
  const p = path.slice(0, -1);
  return (p[0] === 'ui' ? p : ['ui', ...p]).join('-');
};
// {color.neutral.50} -> var(--color-neutral-50); {ui.text.secondary.light} -> var(--ui-text-secondary)
function refs(str) {
  return str.replace(/\{([^}]+)\}/g, (_, r) => {
    const p = r.split('.');
    if (p[0] === 'ui' || p[0] === 'overlay') return `var(--${themeName(p)})`;
    return `var(--${cssName(p)})`;
  });
}
const num = (n) => String(Number(n));

function dim(v, path) {
  if (typeof v === 'string' && v.startsWith('{')) return refs(v);
  const n = parseFloat(v);
  if (n === 0) return '0';
  const dotted = limpia(path).join('.');
  if (GRUPOS_REM.some((g) => dotted.startsWith(g + '.'))) return REM(n);
  return `${num(n)}px`;
}
function fontFamily(arr) {
  return arr.map((f) => (/\s/.test(f) || f === 'Inter' ? `"${f}"` : f)).join(', ');
}
function pct(a) { return `${Number((a * 100).toFixed(2))}%`; }
function mezcla(colorVar, alpha) { return `color-mix(in srgb, ${colorVar} ${pct(alpha)}, transparent)`; }

function sombra(t) {
  const e = ext(t);
  if (e.css) return e.css;
  return t.original.$value.map((c, i) => {
    const col = mezcla(refs(c.color), e.alpha[i]);
    return `${[c.offsetX, c.offsetY, c.blur, c.spread].map((d) => (parseFloat(d) === 0 ? '0' : d)).join(' ')} ${col}`;
  }).join(',\n    ');
}

export function valorCss(t) {
  const v = t.original.$value;
  switch (t.$type) {
    case 'color': {
      const e = ext(t);
      if (e.alpha !== undefined) return mezcla(refs(v), e.alpha);
      return refs(v);
    }
    case 'dimension': return dim(v, t.path);
    case 'fontFamily': return fontFamily(v);
    case 'fontWeight': return num(v);
    case 'number': return ext(t).unit ? `${num(v)}${ext(t).unit}` : num(v);
    case 'duration': return v;
    case 'cubicBezier': return `cubic-bezier(${v.map(num).join(', ')})`;
    case 'shadow': return sombra(t);
    case 'strokeStyle': return typeof v === 'string' ? 'none' : v.dashArray.map((d) => parseFloat(d)).join(' ');
    default: throw new Error(`Tipo DTCG sin formato: ${t.$type} (${t.path.join('.')})`);
  }
}

// ---------- tokens.css ----------
const ENCABEZADO = (nombre) => `/* ${nombre}
 * C2 v1.0 · Finanzas Hogar · GENERADO por Style Dictionary 5.6.1 desde docs/contratos/C2-tokens.json.
 * No se edita a mano: cambia C2-tokens.json y compila (docs/design/tools/tokens/README.md).
 */
`;
const TITULOS = {
  color: 'Color · primitivos (iguales en los dos temas; C2-colores.json)',
  font: 'Tipografía', space: 'Espaciado (retícula de 4 pt)', radius: 'Radios', border: 'Bordes', focus: 'Foco',
  shadow: 'Sombras (valor del tema claro; en oscuro, --ui-shadow-* vale none y se usa borde)', size: 'Medidas de control, icono y componente',
  layout: 'Medidas de estructura', icon: 'Iconos', motion: 'Movimiento', breakpoint: 'Puntos de quiebre (referencia; no sirven dentro de @media)',
  viz: 'Series de gráficas: trazo (stroke-dasharray; none = continuo)',
};

export function formatoCss({ dictionary }) {
  const todos = dictionary.allTokens;
  const tema = todos.filter(esTema);
  const fijos = todos.filter((t) => !esTema(t));
  const sans = todos.find((t) => t.path.join('.') === 'font.family.sans');
  const caras = ext(sans).faces.map((f) => `@font-face {
  font-family: "Inter";
  src: url("${f.src}") format("${f.format}");
  font-weight: ${f.weight};
  font-style: ${f.style};
  font-display: swap;
}
`).join('');

  let r = ENCABEZADO('tokens.css');
  r += `\n/* Fuentes autoalojadas (Inter 4.1, SIL OFL 1.1: ${ext(sans).license}). Rutas relativas a este archivo; la PWA las copia a /fonts/ (WF-3). */\n${caras}`;

  r += `\n:root {\n  color-scheme: light;\n  font-optical-sizing: auto;\n`;
  let grupo = '';
  for (const t of fijos) {
    if (t.path[0] !== grupo) { grupo = t.path[0]; r += `\n  /* ${TITULOS[grupo] ?? grupo} */\n`; }
    r += `  --${cssName(t.path)}: ${valorCss(t)};\n`;
    if (t.path.join('.') === 'font.family.sans') r += `  --font-numeric-tabular: ${ext(t).numeric};\n`;
  }
  const claros = tema.filter((t) => t.path.at(-1) === 'light');
  const oscuros = tema.filter((t) => t.path.at(-1) === 'dark');
  const sombras = fijos.filter((t) => t.$type === 'shadow' && t.path[1] !== '0');
  r += `\n  /* Variables temáticas --ui-* · tema claro (predeterminado) */\n`;
  for (const t of claros) r += `  --${themeName(t.path)}: ${valorCss(t)};\n`;
  // literal (no var(--shadow-N)): Tailwind define --shadow-N: var(--ui-shadow-N) y una referencia mutua sería circular
  for (const t of sombras) r += `  --ui-shadow-${t.path[1]}: ${valorCss(t)};\n`;
  r += `}\n`;

  const bloqueOscuro = (sangria) => {
    const s = ' '.repeat(sangria);
    let b = `${s}color-scheme: dark;\n`;
    for (const t of oscuros) b += `${s}--${themeName(t.path)}: ${valorCss(t)};\n`;
    for (const t of sombras) b += `${s}--ui-shadow-${t.path[1]}: ${ext(t).dark};\n`;
    return b;
  };
  r += `\n/* Tema oscuro según el sistema, salvo que el interruptor fije «claro» */\n@media (prefers-color-scheme: dark) {\n  :root:not([data-theme="light"]) {\n${bloqueOscuro(4)}  }\n}\n`;
  r += `\n/* Tema oscuro fijado por el interruptor (Sistema · Claro · Oscuro) */\n:root[data-theme="dark"] {\n${bloqueOscuro(2)}}\n`;
  r += `\n/* Dinero, porcentajes y cantidades: cifras tabulares (DESIGN §3.3) */\n.money {\n  font-variant-numeric: var(--font-numeric-tabular);\n}\n`;
  r += `\n/* Movimiento reducido: se anulan las transiciones no esenciales (DESIGN §6.3) */
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    transition-duration: 0.01ms !important;
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    scroll-behavior: auto !important;
  }
  .icon-spinner { animation: spin 2s cubic-bezier(0, 0, 1, 1) infinite !important; }
}
`;
  return r;
}

// ---------- tailwind.theme.css ----------
export function formatoTailwind({ dictionary }) {
  const todos = dictionary.allTokens;
  const get = (ruta) => todos.find((t) => t.path.join('.') === ruta);
  const claros = todos.filter((t) => esTema(t) && t.path.at(-1) === 'light');
  let r = ENCABEZADO('tailwind.theme.css');
  r += `
/* Uso: en la hoja de entrada del frontend, en este orden: @import "tailwindcss"; @import ".../docs/contratos/tailwind.theme.css"; (este archivo trae tokens.css).
 * Solo se exponen como utilidades los semánticos --ui-*, nunca los primitivos (DESIGN §10.4): un componente no combina colores a mano.
 * Los componentes no usan la variante dark:; el tema lo resuelven las variables --ui-*. */
@import "./tokens.css";

/* Medidas de 44 px para puntero grueso (Android y toda pantalla táctil) */
@custom-variant coarse (@media (any-pointer: coarse));

@theme inline {
  /* Se quitan las escalas por defecto de Tailwind: solo existen las de C2 */
  --color-*: initial;
  --font-*: initial;
  --text-*: initial;
  --font-weight-*: initial;
  --radius-*: initial;
  --shadow-*: initial;
  --ease-*: initial;
  --breakpoint-*: initial;

  /* Color: semánticos temáticos (cada uno resuelve claro u oscuro por var(--ui-*)) */
  --color-transparent: transparent;
  --color-current: currentColor;
`;
  let grupo = '';
  for (const t of claros) {
    const n = themeName(t.path).slice(3); // sin «ui-»
    const g = (t.path[0] === 'ui' ? t.path[1] : t.path[0]).replace(/-surface$/, '');
    if (g !== grupo) { grupo = g; r += `\n  /* ${g} */\n`; }
    r += `  --color-${n}: var(--${themeName(t.path)});\n`;
  }
  r += `\n  /* Tipografía */\n  --font-sans: ${fontFamily(get('font.family.sans').original.$value)};\n`;
  for (const n of ['caption', 'small', 'body', 'title-3', 'title-2', 'title-1']) {
    r += `  --text-${n}: ${dim(get(`font.size.${n}`).original.$value, ['font', 'size', n])};\n  --text-${n}--line-height: ${dim(get(`font.line-height.${n}`).original.$value, ['font', 'line-height', n])};\n`;
  }
  for (const n of ['regular', 'medium', 'semibold', 'bold']) r += `  --font-weight-${n}: ${num(get(`font.weight.${n}`).original.$value)};\n`;
  r += `\n  /* Espaciado: p-1 … p-12 coinciden con space.1 … space.12 */\n  --spacing: ${dim(get('space.1').original.$value, ['space', '1'])};\n`;
  r += `\n  /* Radios */\n`;
  for (const n of ['sm', 'md', 'lg', 'full']) r += `  --radius-${n}: ${dim(get(`radius.${n}`).original.$value, ['radius', n])};\n`;
  r += `\n  /* Sombras (claro; en oscuro valen none y la superficie lleva borde) */\n  --shadow-0: none;\n  --shadow-1: var(--ui-shadow-1);\n  --shadow-2: var(--ui-shadow-2);\n`;
  r += `\n  /* Curvas (las duraciones usan duration-150, duration-200 y duration-250) */\n`;
  for (const n of ['standard', 'emphasized']) r += `  --ease-${n}: ${valorCss(get(`motion.easing.${n}`))};\n`;
  r += `\n  /* Puntos de quiebre: 360 es la base sin prefijo */\n`;
  for (const n of ['phone', 'tablet', 'desktop', 'wide']) r += `  --breakpoint-${n}: ${REM(parseFloat(get(`breakpoint.${n}`).original.$value))};\n`;
  r += `}\n`;
  return r;
}
