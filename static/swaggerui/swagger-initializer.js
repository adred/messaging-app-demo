window.onload = () => {
  window.ui = SwaggerUIBundle({
    url: "/docs/openapi.yaml", // Confirm this path works in your browser by going directly to http://localhost:8080/docs/openapi.yaml
    dom_id: '#swagger-ui',
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    layout: "StandaloneLayout"
  });
};
