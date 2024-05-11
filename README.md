# Pagine
Static website generator.

## Hierarchy

### Site

Site decides

- The shorts of various content root paths (relative).
- Top level directory contents.

### Page

Page is a set of attributions of single page.

- Template to use.
- Customized data used in template.
- Defines contents at different parts in template.

### Content

For each rich text format, there is an HTML generator (except HTML itself).
