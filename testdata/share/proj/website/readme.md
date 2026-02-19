# Website Template

This template is used to test the 'proj add' command functionality.

## Template Structure

This template creates a basic website with:
- index.html - Homepage with sitename in header
- src/img/favicon.svg - Simple favicon

## Add Definitions

After creating a project, you can add files:

- `proj add html <name>` - Adds an HTML page (requires 'title' variable)
- `proj add css <name>` - Adds a CSS file to src/css/
- `proj add js <name>` - Adds a JS file to src/js/

## Usage

Create a new website:
```sh
proj new website mysite -v sitename="My Site"
```

Add pages to the site:
```sh
proj add html contact -v title="Contact Us"
proj add css styles
proj add js app
```
