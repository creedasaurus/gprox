<p align="center">
  <img alt="gprox logo" src="gprox.svg" height="350" />
</p>

[![Download from webinstall.dev](https://img.shields.io/static/v1?label=webi&message=gprox&color=6c71c4&labelColor=fdf6e3&logoWidth=10&logo=data:image/svg+xml;base64,PHN2ZyBhcmlhLWhpZGRlbj0idHJ1ZSIgZGF0YS1wcmVmaXg9ImZhcyIgZGF0YS1pY29uPSJib2x0IiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAzMjAgNTEyIiB3aWR0aD0iMzQiIGhlaWdodD0iNTYiPjxwYXRoIGZpbGw9IiM2YzcxYzQiIGQ9Ik0yOTYgMTYwSDE4MC42bDQyLjYtMTI5LjhDMjI3LjIgMTUgMjE1LjcgMCAyMDAgMEg1NkM0NCAwIDMzLjggOC45IDMyLjIgMjAuOGwtMzIgMjQwQy0xLjcgMjc1LjIgOS41IDI4OCAyNCAyODhoMTE4LjdMOTYuNiA0ODIuNWMtMy42IDE1LjIgOCAyOS41IDIzLjMgMjkuNSA4LjQgMCAxNi40LTQuNCAyMC44LTEybDE3Ni0zMDRjOS4zLTE1LjktMi4yLTM2LTIwLjctMzZ6Ii8+PC9zdmc+Cg==)](https://webinstall.dev/gprox/)
[![GitHub license](https://img.shields.io/github/license/creedasaurus/gprox)](https://github.com/creedasaurus/gprox/blob/main/LICENSE)
[![Github Releases Stats of gprox](https://img.shields.io/github/downloads/creedasaurus/gprox/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=creedasaurus&repository=gprox)

---

This is a very simple local HTTPS proxy. While working on a recent project, I needed to be able to develop and test web changes locally using HTTPS. I was using the nice little NodeJS package [local-ssl-proxy](https://github.com/cameronhunter/local-ssl-proxy) with good success. But I decided I wanted to write my own that wouldn't require Node or another external dependency. And I wanted it to be a really simple install so my teammates could also use it without much effort. So I wrote `gprox`. It's still a bit of a WIP, but it's essentially a port of `local-ssl-proxy` written in Go, making it easy to compile and distribute for different systems.

_I feel like I should mention, this is only for local development... please dont proxy things in production using this_

### Install

The fastest way to get started:

```
curl -sS https://webinstall.dev/gprox | bash
```

Or if you happen to have `go` installed, you can use:

```sh
go get github.com/creedasaurus/gprox
```

### Run

The default just starts up and serves https from `9001` to `9000`.

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
One small feature of `gprox` is that it has some really basic certs built in (again, not for production use -- just local development). So as soon as you run it, it'll use those by default unless you specify others from the flags. Another quick note about the built in certs -- if you want to save them to files for any reason, you can use the flag `-d, --dropcerts` which will do just that.

