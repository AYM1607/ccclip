# ccclip
End to end encrypted cloud clipboard.

# Installation

Download the prebuilt latest release binary for your OS/arch and ensure that it's in your PATH.

I just do `sudo mv ccclip /usr/local/bin` for either MacOS or Linux.

# Configuration

By default, ccclip looks for a ccclip directory in the [user home directory](https://pkg.go.dev/os#UserHomeDir) and it'll store the configuration file and device keys there.
You don't need to do anything if you want to stores your config in the default location; to override this behavior, use the `--config-dir` global flag.

If you're onboarding a new device (with an already created account) then you should pre-populate your config file (ccclip.yaml in your config dir) to look like this:
```json
{"Email":"{your-email}","DeviceId":""}
```
Yes, I know, the extension is yaml but the file has JSON. I'll change it, don't worry.

It's also a good idea to creat an alias for the command, here's what I do in my `.zshrc`:
```
alias ccc="ccclip"
```
`ccc` is very fun to type, but choose whatever alias make you happy.

# Instructions

### Create an account with

```bash
ccclip register -e {your-email}
```
This will prompt you for your password. Your input won't be shown, that's expected (we like to be safe right!?). Just press enter.
Passwords are stored as bcrypt hashes in the server.

### Register your device

```bash
ccclip register-device
```
This will prompt you for your password.
An [X25519](https://pkg.go.dev/crypto/ecdh) key pair is created and stored in your configuration directory.
A "device" is created and associated to your account, the device id is stored in your configuration file.
Only the public key is sent to the server (duh).

### Save and retrieve your clipboard

The root command `ccclip` serves both purposes, if you pipe something to its stdin, it'll read it and save it as your clipboard.
If you don't pass anything then it'll retrieve your current clipboard and output it through stdout.

There's a caveat: When a clipboard is created, the sending device encrypts the data for each of the currently registered devices
individually. This means that if you register a new device, it won't immediately have access to the current clipboard because
its public key wasn't available when the clipboard was created. The benefit of this is that your data is end-to-end encrypted and I'll never
be able to see it nor an attacker if we're compromised :)

If you want your new device to have access to the current clipboard contents then just run `ccclip | ccclip` from an onboarded device.

# Hosting

The cli can be built by running `go build ./cmd/cli` from the root of the repo.

I host a public API at `https://api.ccclip.io` using [fly.io](https://fly.io/) (thanks for the credits!). Feel free to use it.
There's a Dockerfile for the server at cmd/server/Dockerfile though it will probably won't work for you; it's set up to use fly's LiteFS.
It should be easy enough for you to adapt it for an sqlite db on a mounted volume and run without LiteFS. Maybe try other sqlite distros?
