
<img src="https://github.com/jellyterra/artworks/raw/master/logo/pagine_v2.svg" height="140" alt="Pagine logo" />

# Pagine v2.2.0

Pagine is an high-performance website constructor that makes full use of multicore hardware.

Build jobs can be completed very fast.

## Features

- Parallel hierarchy processing and unit execution. Everything is executed in parallel from beginning to end.
- Hierarchical metadata propagation which makes metadata management easy.
- Manage templates and assets via Git. Every template can be distributed and used without modification.
- In-template builtin functions
- Interact with Pagine in templates.
- Multi-stage workflow.

As server: update web content on file change.

Supported rich text formats:

- [Asciidoc](https://asciidoc.org)
- [Markdown](https://markdownguide.org)
- [MathJax](https://www.mathjax.org) support provided in [Documentor](https://github.com/webpagine/documentor)

## Install

### Binaries

Find the executable that matches your OS and architecture in [releases](https://github.com/webpagine/pagine/v2/releases).

### Build from source

```shell
$ go install github.com/webpagine/pagine/v2/cmd/pagine@v2.2.0
```

> [!TIP]
> Install Pagine via Go mod is recommended for non-amd64 platforms. Choose Go toolchain for your platform [here](https://go.dev/dl).

## Usage

Usage of pagine:
- `-public` string
  - Location of public directory. (default `/tmp/$(basename $PWD).public`)
- `-root` string
  - Site root. (default `$PWD`)
- `-serve` string
  - Specify the port to listen and serve as HTTP.


### Generate

```shell
$ cd ~/web/my_site
$ pagine
Generation complete.
```

### Run as HTTP server

```shell
$ cd ~/web/my_site
$ pagine --serve :12450
```

It automatically executes generation when file changes are detected by `inotify`.

> [!NOTE]
> Incremental generation is not implemented yet.<br/>
> Set the `--public` under `/tmp` is recommended to reduce hard disk writes.

Since v2.1.0, the server provides a WebSocket interface at `/ws` to provide event monitoring for the client, such as page updates.

> [!CAUTION]
> Exposing your Pagine server to the public network might be risky!
> You should deploy your final pages via static page services.

## Structure

### Template

Template is a set of page frames (Go template file) and assets (e.g. SVGs, stylesheets and scripts).

Manifest of one template looks like:
```yaml
manifest:
  canonical: "com.symboltics.pagine.genesis" # Canonical name
  patterns:
    - "/*html"                               # Matched files will be added as template file.

templates:
  - name: "page"        # Export as name `page`
    export: "page.html" # Export `page.html`

  - name: "post"        # Export as name `post`
    export: "post.html" # Export `post.html`
```

To the Go templates files syntax, see [text/template](https://pkg.go.dev/text/template).

Example:  `page.html`
```html
<html lang="{{ .lang }}">
<head>
  <title>{{ .title }}</title>
  <link rel="stylesheet" href="{{ api.Attr.templateBase }}/css/base.css" />
</head>
<body>
{{ template "header.html" . }}
<main>{{ render.ByFileExtName .content }}</main>
</body>
{{ template "footer.html" . }}
</html>
```

### Env

"Environment" is the configuration of the details of the entire process.

```yaml
ignore: # Pattern matching. Matched files will not be **copied** to the public/destination.
  - "/.git*"

use: # Load and set alias for the template.
  genesis: "/templates/genesis"
  documentor: "/templates/documentor"
```

Installing templates via Git submodule is recommended. Such as:

```shell
$ git submodule add https://github.com/webpagine/genesis templates/genesis
```

### Level

Each "level" contains its metadata. And a set of units to be executed.

For directories, metadata sets are stored in `metadata.toml` in the form of map, and units are stored in `unit.toml`

Each template has its alias that defined in `env` as the namespace.

Levels can override fields propagated from parents.

Example: `/metadata.yaml`
```yaml
genesis:
  title: "Pagine"

  head:
    icon: "/favicon.ico"

    nav:
      items:
        - name: "Documentation"
          link: "/docs/"
```

### Unit

Example: `/unit.yaml`
```yaml
unit:
  - template: "genesis"      # Which template to use.
    template_key: "page"     # Which the key refers to.
    output: "/index.html"    # Where to save the result.
    define:
      title: "Pagine"        # Unit-specified metadata.
      content: "README.md"

  - template: "genesis"
    template_key: "page"
    output: "/404.html"
    define:
      title: "Page not found"
      content: "404.md"
```

## Multi-stage Workflow

All paths are based on the directory where the `workflow.yaml` is located.

### Workflow

Workflows are executed in parallel.

Example: `workflow.yaml`
```yaml
stage:
  - job:
      - type: "tsc/v1"
        title: "Build TypeScript"
        path: "/script/"
```

### Stage

Stages are executed in order.
Each stage contains a set of jobs.

### Job

Jobs are executed in parallel.

## Builtin Job Types

### TypeScript compiler (tsc)

```yaml
type: "tsc/v1"
title: "Build TypeScript"
path: "/script/"
```

Make sure `tsc` had been installed on target machine.

## Builtin functions

### Arithmetic

| Func  | Args      | Result |
|-------|-----------|--------|
| `add` | a, b: Int | Int    |
| `sub` | a, b: Int | Int    |
| `mul` | a, b: Int | Int    |
| `div` | a, b: Int | Int    |
| `mod` | a,b : Int | Int    |

### Engine API

| Func            |                                                        | Description                                                                                      |
|-----------------|--------------------------------------------------------|--------------------------------------------------------------------------------------------------|
| `api.Attr`      |                                                        | Get meta information in the form of map about units, hierarchy and templates provided by engine. |
| `api.Apply`     | templateName, templateKey: String, data map[string]Any | Invoke a template.                                                                               |
| `api.ApplyFile` | apiVer, path: String, data map[string]any              | Invoke single template file.                                                                     | 
| `This`          |                                                        | It returns the root node of metadata of the template.                                            |

| api.Attr.      | Description                                     |
|----------------|-------------------------------------------------|
| `isServing`    | True if Pagine is running as server.            |
| `templateBase` | It tells the template where it has been stored. |
| `unitBase`     | Unit's level's base dir path.                   | 

> [!TIP]
> `api.Attr.isServing` can be used to enable some debug code in templates such as page realtime update.

### String processing

| Func                 | Args                        | Result |
|----------------------|-----------------------------|--------|
| `strings.HasPrefix`  | str: String, prefix: String | Bool   |
| `strings.TrimPrefix` | str: String, prefix: String | String |

### Content rendering

Path starts from where the unit is.

| Func                    | Args                   | Description                                         |
|-------------------------|------------------------|-----------------------------------------------------|
| `render.FileByExtName`  | path: String           | Query MIME type by file extension name then render. |
| `render.FileByMimeType` | mimeType, path: String | Render file content by given MIME type.             |

| Format   | MIME Type       | File Extension Name |
|----------|-----------------|---------------------|
| Asciidoc | `text/asciidoc` | `adoc`, `asciidoc`  |
| Markdown | `text/markdown` | `md`, `markdown`    |

## Deploy

### Manually

```shell
$ pagine --public ../public
```

Upload `public` to your server.

### Deploy to GitHub Pages via GitHub Actions (recommended)

GitHub Actions workflow file content: `.github/workflows/pagine.yml`
```yaml
name: Deploy

on:
  push:
    branches:
      - master

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

defaults:
  run:
    shell: bash

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Install Pagine
        run: go install github.com/webpagine/pagine/v2/cmd/pagine@v2.2.0

      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0

      - name: Setup Pages
        id: pages
        uses: actions/configure-pages@v4

      - name: Build with Pagine
        run: ~/go/bin/pagine --public ./public/

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public/

      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## FAQ

### Why another generator? Isn't Hugo enough?

Pagine is **not** Hugo, and is not aim to replace Hugo.

Pagine does not embed page configurations in Markdown file, they are separated and should be separated.

And Pagine does not focus on Markdown only, I hope to support various kinds of source.

### Can I use Pagine for building complex web application?

Yes. It can help you get rid of repetitive work about static contents.

And templates can also increase productivity as long as Pagine well integrated with external tools via ***workflow***.

### Co-operate with external tools such as TypeScript compiler or npx?

It can be done via ***workflow***.

### What is the origin of the logo and name?

It is **neither** a browser engine, a layout engine **nor** a rendering engine.

Page Gen × Engine ⇒ Pagine. It has similar pronunciation to "pagen".

The logo is an opened book with a bookmark.

### Rewrite it in other PL?

*I expected somebody would ask.*

It will not be taken unless it does bring obvious advantages.

Thus: NO. It is not planned currently.
