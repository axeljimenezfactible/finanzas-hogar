// Aplica el tema guardado antes de pintar (evita el destello). «sistema» = sin atributo: manda prefers-color-scheme.
try { var t = localStorage.getItem('fh-tema'); if (t === 'claro' || t === 'oscuro') document.documentElement.setAttribute('data-theme', t === 'claro' ? 'light' : 'dark') } catch (e) {}
