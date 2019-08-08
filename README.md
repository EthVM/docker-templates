docker-templates ![version v0.1.0](https://img.shields.io/badge/version-v0.1.0-brightgreen.svg) ![License MIT](https://img.shields.io/badge/license-MIT-blue.svg)
================

Utility oriented to render Docker Compose / Stack files with the power of [go templates]() (and heavily inspired by [dockerize](https://github.com/EthVM/docker-templates) and [docker compose templer](https://github.com/Aisbergg/python-docker-compose-templer)).

**Problem**: Docker Compose / Stack files are very static in nature as you only can use YAML to define them. That's a bummer because you can't add conditionals, neither iterations, scoped blocks...The only way to customise the files are [using environment variables substitution](https://docs.docker.com/compose/environment-variables/) the rest is forbidden.

**Solution**: Use `docker-templates`!

## Roadmap

For now, the initial version is rather limited in scope and features with bugs included! For now it fits properly the bill for the intention we are using it on [EthVM](https://github.com/EthVM/ethvm). But here's a list of next features:

- [ ] Watch for files dynamically and auto-render files.
- [ ] Add test suite to verify correctness (and also to be a good citizen).
- [ ] More complex use cases for rendering partial templates inside other templates.

## Installation

Download the latest version in your container:

* [linux/amd64](https://github.com/EthVM/docker-templates/releases/download/v0.1.0/docker-templates-linux-amd64-v0.1.0.tar.gz)
* [alpine/amd64](https://github.com/EthVM/docker-templates/releases/download/v0.1.0/docker-templates-alpine-linux-amd64-v0.1.0.tar.gz)
* [darwin/amd64](https://github.com/EthVM/docker-templates/releases/download/v0.1.0/docker-templates-darwin-amd64-v0.1.0.tar.gz)

### Docker Base Image

The `ethvm/docker-templates` image is a base image based on `alpine linux`. `docker-templates` is installed in the `$PATH` and can be used directly.

```
FROM ethvm/docker-templates
...
ENTRYPOINT docker-templates ...
```

### Ubuntu Images

``` Dockerfile
RUN apt-get update && apt-get install -y wget

ENV DOCKER_TEMPLATES_VERSION v0.1.0
RUN wget https://github.com/EthVM/docker-templates/releases/download/$DOCKER_TEMPLATES_VERSION/docker-templates-linux-amd64-$DOCKER_TEMPLATES_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf docker-templates-linux-amd64-$DOCKER_TEMPLATES_VERSION.tar.gz \
    && rm docker-templates-linux-amd64-$DOCKER_TEMPLATES_VERSION.tar.gz
```

### For Alpine Images:

``` Dockerfile
RUN apk add --no-cache openssl

ENV DOCKER_TEMPLATES_VERSION v0.1.0
RUN wget https://github.com/EthVM/docker-templates/releases/download/$DOCKER_TEMPLATES_VERSION/docker-templates-alpine-linux-amd64-$DOCKER_TEMPLATES_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf docker-templates-alpine-linux-amd64-$DOCKER_TEMPLATES_VERSION.tar.gz \
    && rm docker-templates-alpine-linux-amd64-$DOCKER_TEMPLATES_VERSION.tar.gz
```

## Usage

### Command line arguments

```text
NAME:
   docker-templates - render Docker Compose / Stack file templates with the power of go templates

USAGE:
   docker-templates [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   render   renders the specified definition file(s)
   help, h  shows a list of commands or help for one command

GLOBAL OPTIONS:
   --stdout           forces output to be written to stdout
   --delims value     template tag delimiters. Default "{{":"}}" (default: "{{:}}")
   --log-level value  log level to emit to the screen (default: 4)
   --help, -h         show help
   --version, -v      print the version
```

### Definition File

The definition file defines what to do. It lists template and the variables to be used for rendering and says where to put the resulting file. The definition file syntax is as follows:

```toml

# Example definition file

[vars]

    # and/or define them here (higher priority)
    [vars.global]
    network_enabled = true
    network_name = "net"
    network_subnet = "172.25.0.0/16"

    # You can include other global variables (lower priority)
    include = [
      "vars/global.toml"
    ]


[[templates]]
src  = "templates/stack.yml.tpl"
dest = "out/stack-1.yml"
include_vars = [ "vars/local.toml" ]
[templates.local_vars]
mariadb_version = "10.2.21"
mariadb_volume_enabled = true

[[templates]]
src  = "templates/stack.yml.tpl"
dest = "out/stack-2.yml"
include_vars = []
[templates.local_vars]
mariadb_version = "11"
mariadb_volume_enabled = false
```

And one example of external vars:

```toml
[vars]
this_is_another_var = 'var'
```

The different sources of variables are merged together in the following order:

1. global `vars`
2. global `include`
3. template `include_vars`
4. template `vars`

### Templates

Templates use Golang [text/template](http://golang.org/pkg/text/template/) and the files are rendered with them. Only the definition file and vars files are written in TOML, but the output of the rendered templates can be anything (be YAML or other format). You can access environment variables within a template with `.Env` like `dockerize` or those defined in the definition file with plain `.` (like `.some_global_var`).

```
{{ .Env.PATH }} is my path
```

There are a few built in functions as well (that have been stolen from `dockerize`):

  * `default $var $default` - Returns a default value for one that does not exist. `{{ default .Env.VERSION "0.1.2" }}`
  * `contains $map $key` - Returns true if a string is within another string
  * `exists $path` - Determines if a file path exists or not. `{{ exists "/etc/default/myapp" }}`
  * `split $string $sep` - Splits a string into an array using a separator string. Alias for [`strings.Split`][go.string.Split]. `{{ split .Env.PATH ":" }}`
  * `replace $string $old $new $count` - Replaces all occurrences of a string within another string. Alias for [`strings.Replace`][go.string.Replace]. `{{ replace .Env.PATH ":" }}`
  * `parseUrl $url` - Parses a URL into it's [protocol, scheme, host, etc. parts][go.url.URL]. Alias for [`url.Parse`][go.url.Parse]
  * `atoi $value` - Parses a string $value into an int. `{{ if (gt (atoi .Env.NUM_THREADS) 1) }}`
  * `add $arg1 $arg` - Performs integer addition. `{{ add (atoi .Env.SHARD_NUM) -1 }}`
  * `isTrue $value` - Parses a string $value to a boolean value. `{{ if isTrue .Env.ENABLED }}`
  * `lower $value` - Lowercase a string.
  * `upper $value` - Uppercase a string.
  * `jsonQuery $json $query` - Returns the result of a selection query against a json document.
  * `loop` - Create for loops.

## Differences

`docker-templates` is pretty similar to `python-docker-compose-templer` so why to create another similar program? Mainly for this reasons:

* I wanted to have a simpler binary that I can add easily to my docker files, without having to install Python in them.
* I like Go :')

And what about `dockerize`?

* `dockerize` is meant to be used mainly with `.Env` variables, but whenever you have a more complex use cases, sometimes falls sort.

## Acknowledgements

Many thanks to jwilder and [Aisbergg](https://github.com/Aisbergg) for creating [dockerize](https://github.com/EthVM/docker-templates) and [python-docker-compose-templer](https://github.com/Aisbergg/python-docker-compose-templer) respectively, from which this project draws 99% inspiration!

## License

*docker-templates* is released under the MIT License. See [LICENSE.txt](LICENSE.txt) for more information.

[go.string.Split]: https://golang.org/pkg/strings/#Split
[go.string.Replace]: https://golang.org/pkg/strings/#Replace
[go.url.Parse]: https://golang.org/pkg/net/url/#Parse
[go.url.URL]: https://golang.org/pkg/net/url/#URL
