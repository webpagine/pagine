
<img src="https://github.com/jellyterra/artworks/raw/master/logo/pagine_v2.svg" height="140" alt="Pagine logo" />

# Pagine v2.2

Pagine is an high-performance website constructor that makes full use of multicore hardware.

Build jobs can be completed very fast.

## Features

- Parallel hierarchy processing and unit execution. Everything is executed in parallel from beginning to end.
- Hierarchical metadata propagation which makes metadata management easy.
- Manage templates and assets via Git. Every template can be distributed and used without modification.
- In-template builtin functions
- Interact with Pagine in templates.
- Update on file change while running as HTTP server.

Supported rich text formats:

- [Markdown](https://markdownguide.org) with [MathJax](https://www.mathjax.org) support
- [Asciidoc](https://asciidoc.org)

## Install

### Binaries

Find the executable that matches your OS and architecture in [releases](https://github.com/webpagine/pagine/v2/releases).

### Build from source

```shell
$ go install github.com/webpagine/pagine/v2/cmd/pagine@latest
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
    - "/*html"                    # Matched files will be added as template file.

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
  <link rel="stylesheet" href="{{ (getAttr).templateBase }}/css/base.css" />
</head>
<body>
{{ template "header.html" .header }}
<main>{{ render .content }}</main>
</body>
{{ template "footer.html" .footer }}
</html>
```

### Env

"Environment" is the configuration of the details of the entire process.

```yaml
ignore:
  - "/.git*" # Pattern matching. Matched files will not be **copied** to the public/destination.

use:
  genesis: "/templates/genesis"
  another: "/templates/something_else" # Load and set alias for the template.
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
  - template: "genesis:page" # Which template to use.
    output: "/index.html"    # Where to save the result.
    define:
      title: "Pagine"        # Unit-specified metadata.
      content: "README.md"

  - template: "genesis:page"
    output: "/404.html"
    define:
      title: "Page not found"
      content: "404.md"
```

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

| Func      | Description                                                                                      |
|-----------|--------------------------------------------------------------------------------------------------|
| `getEnv`  | Get environment information for debug use.                                                       | 
| `getAttr` | Get meta information in the form of map about units, hierarchy and templates provided by engine. |

| Environment | Description                                 |
|-------------|---------------------------------------------|
| `isServing` | Return true if Pagine is running as server. |

| Attribution    | Description                                     |
|----------------|-------------------------------------------------|
| `unitBase`     | Unit's level's base dir path.                   |
| `templateBase` | It tells the template where it has been stored. | 

| Func          | Description                                           |
|---------------|-------------------------------------------------------|
| `getMetadata` | It returns the root node of metadata of the template. |

> [!TIP]
> `(getEnv).isServing` can be used to enable some debug code in templates such as page realtime update.

### Data processing

| Func         | Args                        | Result |
|--------------|-----------------------------|--------|
| `hasPrefix`  | str: String, prefix: String | Bool   |
| `trimPrefix` | str: String, prefix: String | String |

| Func             | Args                                            | Result                                           | Description                                                                       |
|------------------|-------------------------------------------------|--------------------------------------------------|-----------------------------------------------------------------------------------|
| `divideSliceByN` | slice: []Any, n: Int                            | [][]Any                                          | Divide a slice into *len(slice) / N* slices                                       |
| `mapAsSlice`     | map: map[String]Any, **key**, **value**: String | []map[String]{ **key**: String, **value**: Any } | Convert map to a slice of map that contains two keys named **key** and **value**. |

### Content

Path starts from where the unit is.

| Func             | Args                              | Description                             |
|------------------|-----------------------------------|-----------------------------------------|
| `apply`          | path: String, data map[string]Any | Invoke a template.                      |
| `embed`          | path: String                      | Embed file raw content.                 |
| `render`         | path: String                      | Invoke renderer by file extension name. |
| `renderAsciidoc` | path: String                      | Render and embed Asciidoc content.      |
| `renderMarkdown` | path: String                      | Render and embed Markdown content.      |

| Format   | File Extension Name |
|----------|---------------------|
| Markdown | `md`                |
| Asciidoc | `adoc`              |

## Deploy

### Manually

```shell
$ pagine --public ../public
```

Upload `public` to your server.

### Deploy to GitHub Pages via GitHub Actions (recommended)

GitHub Actions workflow configuration can be found in [Get Started](https://github.com/webpagine/get-started) repository.

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
