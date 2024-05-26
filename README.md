
<img src="https://github.com/jellyterra/artworks/raw/master/logo/pagine.svg" width="410.4" height="140" alt="Pagine logo" />

# Pagine
Template-driven generator for building websites of any scale.

Latest version: v1.0.0

- Template.
- Separation of **content**, **template** and **page form**.

### Planned features

- Job pipeline for reducing redundant content generation.
- Directories as collections of web pages.

## Install

```shell
$ go install github.com/webpagine/go-pagine/cmd/pagine
$ pagine --gen
```

Serve as HTTP server and automatically generate when files change:

```shell
$ pagine --serve --listen :8080 --public /tmp/public
```

> [!NOTE]
> Incremental generation is not implemented yet.<br/>
> Set the `--public` under `/tmp/` is recommended to reduce hard disk writes.

## Get Started

Example structure:
```
.
├── pagine.toml
├── contents/
│   └── my_first_post.md
├── data/
│   ├── header_all.toml
│   └── header_specific.toml
├── posts/
│   └── my_first_post.html.pagine
└── templates/
    ├── header.html
    ├── footer.html
    └── post.html
```

### Site

- Top level directory contents, such as website metadata.
- Global elements, such as page frame, navigation bar.

For example: `/pagine.toml`
```toml
ignore = [  "/\\.*", "/*toml", "/contents/*", "/templates/*" ]
```

### Template

Current implementation of template depends on Go `text/template` library.

For Go templates, refer to the [tutorial](https://gohugo.io/templates/introduction/) by Hugo team.

For example: `/templates/post.html`
```html
<html lang="{{ .data.lang }}">
<head>
    <title>{{ .data.title }}</title>
</head>
<body>
<p>{{ .contents.my_first_post }}</p>
<p>{{ .data.time }}</p>
</body>
</html>
```

### Page

Page is a set of attributions of single page.

- Templates to be used.
- Data definitions to be used in template.
- Different contents to be used in template.

For example: `/posts/my_first_post.html.pagine`
```toml
# Templates to be used in this page.
[templates]
header = "/templates/header.html"
footer = "/templates/footer.html"

# Main template (top-level) is required.
main   = "/templates/post.html"

# Include data definitions from extern TOMLs.
[include]
header = [
    "/data/header_all.toml",
    "/data/header_specific.toml",
]

# Contents to be parsed and generated to HTML.
[contents]
my_first_post = "/contents/my_first_post.md"

# Define data for template "main".
[define.main]
lang  = "en"
title = "My First Post"
time  = 2024-05-01

# Define data for template "header".
[define.header]
logo = "/assets/img/logo.svg"
```

### Content

For each supported rich text format, there is a parser and an HTML generator. Pagine detects format by file name suffix `.md`.

The latest implementation accepts:
- Markdown

For example: `/contents/my_first_post.md`
```markdown
# My First Post

It is a post in Markdown.
```

## Deploy manually

```shell
$ pagine --gen --public ../public
```

## Deploy via CI/CD

### GitHub Actions

GitHub Actions workflow configuration can be found in [Get Started](https://github.com/webpagine/get-started) repository.

# FAQ

### Why another generator? Isn't Hugo enough?

Pagine is **not** Hugo, and is not aim to replace Hugo.

Pagine does not embed page configurations in Markdown file, they are separated and should be separated.

And Pagine does not focus on Markdown only, I hope to support various kinds of source.

### Can I use Pagine for building complex web application?

It can only help you get rid of repetitive work about static contents.

Templates can increase productivity as long as Pagine well integrated with external tools.

So, **it depends**.

### Co-operate with external tools such as npx?

It is possible. This step should be transparent to external tools.

Run `npx` in *public* directory after the generation by Pagine.

### What is the origin of the logo and name?

It is **neither** a browser engine, a layout engine **nor** a rendering engine.

Page Gen × Engine ⇒ Pagine. It has similar pronunciation to "pagen".

The logo is an opened book with a bookmark.

### Rewrite it in other PL?

*I expected somebody would ask.*

It will not be taken unless it does bring obvious advantages.

Thus: NO. It is not planned currently.
