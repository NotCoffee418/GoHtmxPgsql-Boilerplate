package config

import "time"

const WebsiteTitle = "GoHtmxPgsql-Boilerplate"

// Redirects to HTTPS when schema is HTTP
// Reverse proxy should handle this, but it's here if needed.
const HttpsRedirect = false

// Internal port to listen on
const ListenPort = 8080

// Timeout in seconds for server to wait for response
const TimeoutSeconds = 15

// Should go compile minified global.css using npx postcss?
// Can be done manually or on compile instead with
// Note that this is also required for tailwindcss to work.
// `npx postcss ./static/css/global.css -o ./static/css/global.min.css`
const DoMinifyCss = true

// Database Connection Pooling settings
const DbMaxOpenConns int = 10
const DbMaxIdleConns int = 5
const DbConnMaxLifetime time.Duration = time.Hour
