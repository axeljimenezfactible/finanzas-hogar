// Package db expone las migraciones goose de backend/db/migrations embebidas en el binario.
// go:embed no puede subir de carpeta, por eso este archivo vive junto a migrations/ (las migraciones no se tocan).
package db

import "embed"

// Migrations contiene migrations/*.sql (esquema 1.4: 00001 a 00004 y las que se agreguen).
//
//go:embed migrations/*.sql
var Migrations embed.FS
