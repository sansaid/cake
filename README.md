# cake

`cake` is a CLI system that turns your Raspberry Pi into a tiny container running machine.

Running `cake` on your Raspberry Pi is all you'll need to run any container on it. Once you've got `cake` running on your sweet Raspi, you won't need to login again to update your container. You just need to push your changes to your Docker Hub repo and the container on your Raspi will automatically run!

To install `cake` on your Raspi, just run:

```
aptitude install cake
```

## Project Structure

`cake` is minimal and only contains 2 components (which is currently maintained by 1 developer). So `cake` currently exists as a monorepo, with each component maintained under its own directory.

* `cakectl` - this is where the `cakectl` lives
* `caked` - this is where the `cake` daemon lives

Both components are required for `cake` to run.

## Components

### cakectl

`cakectl` is simply a way to interface with the `caked` daemon and tell it what to do from a friendly CLI.

### caked

`caked` is a daemon that runs the core part of `cake`'s logic. `caked` and `cakectl` talk to each other using gRPC.

## Usage

Deploy your container

```sh
cakectl start --image "sansaid:debotbot/image"
```
