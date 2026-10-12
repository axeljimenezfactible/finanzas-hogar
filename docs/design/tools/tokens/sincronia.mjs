// C2-09: compila dos veces a una carpeta temporal y compara byte a byte con los archivos versionados de docs/contratos/.
// Uso: npm run check:sync   (codigo 0 = identicos y deterministas)
import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, rmSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import StyleDictionary from 'style-dictionary';
import config from './config.mjs';

const aqui = (r) => fileURLToPath(new URL(r, import.meta.url));
const tmp = aqui('../../../../datos/tmp/c2-sync/');
const contratos = aqui('../../../contratos/');
const archivos = ['tokens.css', 'tailwind.theme.css'];
const sha = (f) => createHash('sha256').update(readFileSync(f)).digest('hex');
let mal = 0;
const corridas = [];
for (const n of [1, 2]) {
  rmSync(tmp, { recursive: true, force: true });
  mkdirSync(tmp, { recursive: true });
  const cfg = structuredClone({ ...config, hooks: undefined });
  cfg.hooks = config.hooks;
  cfg.platforms.web.buildPath = tmp;
  const sd = new StyleDictionary(cfg);
  await sd.buildAllPlatforms();
  corridas.push(archivos.map((a) => sha(tmp + a)));
}
archivos.forEach((a, i) => {
  const versionado = sha(contratos + a);
  const igual = corridas[0][i] === corridas[1][i] && corridas[0][i] === versionado;
  console.log(`${igual ? 'OK ' : 'DIF'} ${a} sha256 ${versionado}`);
  if (!igual) mal++;
});
rmSync(tmp, { recursive: true, force: true });
process.exit(mal ? 1 : 0);
