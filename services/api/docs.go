package api

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

const docsTmpl = `
<!doctype html>
<head>
	<title>e2clicker API Reference</title>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
</head>

<body>
	<script id="api-reference" data-url="/api/openapi.json"></script>
	<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
`

func mountDocs(mux *http.ServeMux, swagger *openapi3.T) {
	mux.HandleFunc("GET /api/openapi.json", respondSwagger(swagger))
	mux.HandleFunc("GET /api/docs", respondDocs)
}

func respondSwagger(swaggerAPI *openapi3.T) http.HandlerFunc {
	b, err := swaggerAPI.MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("cannot marshal swagger API: %w", err))
	}

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func respondDocs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(docsTmpl))
}
