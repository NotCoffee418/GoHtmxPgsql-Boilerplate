package config

import (
	"compress/gzip"
	"time"
)

const WebsiteTitle = "GoWebsite-Boilerplate"

// ListenPort Internal port to listen on
const ListenPort = 8080

// TimeoutSeconds for server to wait for response
const TimeoutSeconds = 15

// DoMinifyCss indicates if the app should compile minified global.css using npx postcss?
// Can be done manually or on compile instead with
// Note that this is also required for tailwindcss to work.
// uses `npx postcss`
const DoMinifyCss = true

// DbMaxOpenConns Database Connection Pooling settings
const DbMaxOpenConns int = 10
const DbMaxIdleConns int = 5
const DbConnMaxLifetime time.Duration = time.Hour

// WsReadBufferSize WebSocket settings
const WsReadBufferSize = 1024
const WsWriteBufferSize = 1024

// Set level of gzip compression
const GzipCompressionLevel = gzip.BestSpeed
