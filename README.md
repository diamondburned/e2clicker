<div align="center">
    <h1>e2clicker</h1>
    <img src="https://github.com/diamondburned/e2clicker/assets/8463786/39f7fda5-cea8-4eb3-8bd3-cbff9f8e0cfa" alt="screenshot" width="300" />
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
