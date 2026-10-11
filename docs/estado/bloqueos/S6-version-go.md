# Bloqueo S6-version-go
Contexto: Al crear `backend/go.mod` (fase S6 parte a). La decisión I aprueba Go 1.27 y el criterio S6-01 espera `go version` = go1.27.x y `go 1.27` en go.mod. La máquina de casa tiene Go 1.26.5 (`go version go1.26.5-X:nodwarf5 linux/amd64`) y la base de paquetes sincronizada de pacman solo ofrece `go 2:1.26.5-1`. Instalar Go 1.27 con pacman obligaría a sincronizar repositorios y actualizar el sistema completo (`pacman -Syu`); actualizar un solo paquete sobre una base vieja es una actualización parcial, que Arch no admite.
Qué falta: ¿Cómo se obtiene Go 1.27 en la máquina de casa?
Opciones:
- A. Axel corre `sudo pacman -Syu` (actualización completa del sistema, con reinicio si toca el kernel) y luego se cambia go.mod a `go 1.27`. Consecuencia: cumple S6-01 con el Go del sistema; es un cambio grande para la máquina.
- B. Dejar que Go descargue su propia cadena 1.27 (`GOTOOLCHAIN=auto` con `go 1.27.x` en go.mod) en la caché de módulos del usuario. Consecuencia: no toca el sistema, pero la primera compilación necesita internet y la cadena vive fuera de pacman.
- C. Seguir con Go 1.26 hasta la siguiente actualización normal del sistema. Consecuencia: S6-01 no se puede declarar cumplido; el código actual no usa nada exclusivo de 1.27.
Impacto: Solo la versión declarada; `go build ./...` y `go vet ./...` pasan con 1.26.5.
Cómo seguí: `backend/go.mod` declara `go 1.26` (compatible con el Go instalado, como pide la tarea) y no instalé nada con pacman.
