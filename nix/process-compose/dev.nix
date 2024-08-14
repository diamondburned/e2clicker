rec {
  backend = {
    port = 8001;
    command = "just dev-backend";
    healthPath = "/debug/health";
    environment = {
      DATABASE_URI = "postgresql://e2clicker@localhost";
      HTTP_ADDRESS = "localhost:${toString backend.port}";
    };
  };
  frontend = {
    port = 8000;
    command = "just dev-frontend --port ${toString frontend.port}";
    healthPath = "/";
    environment = {
      BACKEND_HTTP_ADDRESS = "http://localhost:${toString backend.port}";
    };
  };
  container-postgresql = {

  };
}
