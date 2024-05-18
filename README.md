
<img src="https://github.com/jellyterra/artworks/raw/master/logo/pagine.svg" width="410.4" height="140" alt="Pagine logo" />

# Pagine
Static web page generator for blogs and showcases.

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
├── posts/
│   └── my_first_post.html.pagine
└── templates/
    └── post.html
```

### Site

- Top level directory contents, such as website metadata.
- Global elements, such as page frame, navigation bar.

For example: `/pagine.toml`
```toml
ignore = [  "/\\.*", "/*toml", "/*pagine", "/contents/*", "/templates/*" ]
```

### Template

Current implementation of template depends on Go `text/template` library.

For Go templates, refer to the [tutorial](https://gohugo.io/templates/introduction/) by Hugo team.

For example: `/templates/post.html`
```html
<html lang="{{ $.lang }}">
<head>
    <title>{{ $.title }}</title>
</head>
<body>
<p>{{ $.content }}</p>
</body>
</html>
```

### Page

Page is a set of attributions of single page.

- Template to use.
- Customized data used in template.
- Defines contents at different parts in template.

For example: `/posts/my_first_post.html.pagine`
```toml
template = "/templates/post.html"

[content]
content = "/contents/my_first_post.md"

[data]
lang = "en"
title = "My First Post"
```

`[content]` Content sources to be rendered to HTML.

- The key defines the name used in template such as `{{ $.content }}`.
- The value refers to the path of content source.

`[data]` Customized values.

- The key defines the name used in template such as `{{ $.title }}`.
- The value is the value. It remains AS IS.

### Content

For each supported rich text format, there is a parser and an HTML generator. Pagine detects format by file name suffix `.md`.

The latest implementation accepts:
- Markdown

For example: `/contents/my_first_post.md`
```markdown
# My First Post

It is a post in Markdown.
```

## Deploy

```shell
$ pagine --gen
```

Currently:
- Relative CI/CD is not implemented.
- The only approach to deploy is to upload entire generated site to the server.

## FAQ

### Why another generator? Isn't Hugo enough?

Pagine is **not** Hugo, and is not aim to replace Hugo.

Pagine does not embed page configurations in Markdown file, they are separated and should be separated.

And Pagine does not focus on Markdown only, I hope to support various kinds of source.

### Can I use Pagine for building complex web application?

It can only help you get rid of repetitive work about static contents.

Templates can increase productivity as long as Pagine well integrated with external tools.

So, **it depends**.

### Co-operate with external tools such as npx?

It is planned but not implemented in the latest version.

### What is the origin of the logo and name?

It is **neither** a browser engine, a layout engine **nor** a rendering engine.

Page Gen × Engine ⇒ Pagine. It has similar pronunciation to "pagen".

The logo is an opened book with a bookmark.

### Rewrite it in other PL?

*I expected somebody would ask.*

It will not be taken unless it does bring obvious advantages.

Thus: NO. It is not planned currently.
