<div align="center">
    <h1>e2clicker</h1>
    <img src="./assets/screenshots/dashboard-light.png" alt="screenshot" width="720" />
    <br />
    <br />
</div>

This service was written because this one was struggling to keep track of when to apply its estrogen patches,
so it built a service to keep track of that for it. Being a pure Go service that relies on [Gotify][gotify]
to send out push notifications to its phone, no Java code was needed.

[gotify]: https://gotify.net

## Architecture

- Web frontend that allows recording HRT history and displaying statistics
- Configuration:
  - Recommended time intervals for HRT
- Gotify integration: push notifications when it's time to do HRT

## Developing

e2clicker is designed so that you only need Nix to bootstrap the development
process.

- First, run `nix develop` to enter the Nix shell. You may also choose to use
  Direnv with `use flake` for an easier time.
- Then, run `just dev` to start the development server. This command will
  automatically start a Zellij session consisting of:
  - A NixOS virtual machine, which contains all the dependencies needed
    for e2clicker's backend to operate.
  - A Vite development server, which automatically reverse proxies API calls
    to the backend virtual machine while allowing hot reloading.
- To stop the development server:
  - First shut off the virtual machine gracefully by running `shutdown now` on
    the top pane.
  - Then, stop the Vite server by focusing on the bottom pane and pressing
    `Ctrl+C`.
  - Finally, you may exit out of all the Zellij panes. This will close Zellij
    entirely.

## Building

e2clicker exposes the backend and frontend packages as 2 separate Nix flake outputs:

```
git+file:///home/diamond/Projects/diamondburned/e2clicker
└───packages
    └───<ARCHITECTURE>
        ├───e2clicker-backend: package 'e2clicker-backend'
        └───e2clicker-frontend: package 'e2clicker-frontend'
```

You can build these 2 packages without needing to clone the repository by
running the following:

```sh
nix build "github:diamondburned/e2clicker#e2clicker-backend"
nix build "github:diamondburned/e2clicker#e2clicker-frontend"
```

## Deploying

e2clicker exposes a NixOS module that can be used to deploy it onto a server.

It also exposes a PostgreSQL module that initializes the e2clicker database for
it to use in the same machine. This module is optional; you can choose to use
any PostgreSQL database as needed.
