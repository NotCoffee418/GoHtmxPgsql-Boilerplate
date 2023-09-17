package config

const WebsiteTitle = ""

// Redirects to HTTPS when schema is HTTP
// Reverse proxy should handle this, but it's here if needed.
const HttpsRedirect = false

// Internal port to listen on
const ListenPort = 8080

// Timeout in seconds for server to wait for response
const TimeoutSeconds = 15
