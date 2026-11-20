This is a simple [Javascript Service Worker](https://python.org/) server template running with [WinterJS](https://github.com/wasmerio/winterjs).

> This starter's full tutorial is available [here](https://docs.wasmer.io/edge/quickstart/js-wintercg).

## Usage

Modify the logic of your the Javascript worker in the `src/index.js` file.

You can run the JS Service Worker locally with (check out the [Wasmer install guide](https://docs.wasmer.io/install)):

```bash
wasmer run .
```

Open [http://localhost:8080](http://localhost:8080) with your browser to see the worker working!

## Tests

Unit tests use the built-in Node.js test runner, no dependencies to install. They load `src/index.js` with a shimmed `addEventListener` and send Telegram-like updates to the fetch handler.

```bash
npm test
```

Requires Node.js 24 or newer. The same command runs in CI on every push and pull request to `main` (`.github/workflows/test.yml`).

## Deploy on Wasmer Edge

The easiest way to deploy your Javascript Worker is to use the [Wasmer Edge](https://wasmer.io/products/edge).

Live example: https://wasmer-js-worker-starter.wasmer.app/

```bash
wasmer deploy
```
