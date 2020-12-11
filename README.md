# gprox

This is a really simple static local SSL proxy. While working I needed to be able to develop and test web changes locally using HTTPS. I was using the nice little NodeJS package [local-ssl-proxy](https://github.com/cameronhunter/local-ssl-proxy) with good success. But I decided I wanted to write my own that wouldn't require Node or any other dependency. And I wanted it to be a really simply install so my teammates could also use it without much effort. So I wrote `gprox`. It's still a little bit of a WIP, but it's essentially a port of `local-ssl-proxy` written in Go, making it easy to compile and distribute for different systems.

**I feel like I should mention, this is only for local development... please dont proxy things in production using this**

### Install

Currently I'm working on writing the installer for https://webinstall.dev/, which is a great way to install things quick and without needing `root`.

Until then, you can download the proper binary for your system from the [releases](https://github.com/creedasaurus/gprox/releases) page and install as you please.

Or if you happen to have `go` installed, you can use:

```sh
go get github.com/creedasaurus/gprox
```

### Run

One small feature of `gprox` is that it has some really basic certs built in (again, not for production use -- just local development). So as soon as you run it, it'll use those as defaults and start up a small proxy server for you from port `9001` to `9000`.

```
gprox
```

And you should see the output like:

```
9:12PM INF Running proxy! from=https://localhost:9001 to=http://localhost:9000
```

### Configure

There are a few options that are currently available (and a few that I'm still working on). You can use the `-h, --help` flag to see the options:

```
❯ gprox --help
Usage:
  gprox [OPTIONS]

Application Options:
  -n, --hostname=  The hostname to be used for the local proxy (default: localhost)
  -s, --source=    The source port that you will hit to go through the proxy (default: 9001)
  -t, --target=    The port you are targeting (default: 9000)
  -c, --cert=      Path to a .cert file
  -k, --key=       Path to a .key file
  -o, --config=
  -d, --dropcerts  Save the built-in cert/key files to disk
      --version

Help Options:
  -h, --help       Show this help message
```

So you can specify the source and target ports of course. But you can also generate your own `.crt` and `.key` files using `openssl` and some command like:

```
openssl req -x509 \
    -out local.example.crt \
    -keyout local.example.key \
    -newkey rsa:2048 \
    -nodes -sha256 \
    -subj '/CN=local.example.com' \
    -extensions EXT -config <( \
        printf "[dn]\nCN=local.example.com\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:local.example.com\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
```

This will create `local.example.crt` and `local.example.key` files for you. If you want to change the hostname that the local proxy is using in order to match the hostname for the certs you just generated, you can add another name for localhost to your `/etc/hosts` file like: 

```
# /etc/hosts file (requires sudo to edit -- pls do at your own risk)
# add the following line:
127.0.0.1       local.example.com
```

Now, you can run `gprox` with this hostname and certs to match using:

```
gprox --hostname local.example.com \ 
    --cert local.example.crt \
    --key local.example.key \ 
    --target 8080
```

Which will proxy `https://local.example.com:9001` to `http://local.exmaple.com:8080` (localhost).

#### Built-in certs

Another quick note about the built in certs. If you want to import the built-in certs for any reason, you can use the flag `-d, --dropcerts` to save the built-ins to local files. 

