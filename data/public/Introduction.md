---
title: "Welcome to Project R.E.D. - Node Introduction"
authors: ["Standard Codebase"]
contributors: ["Community"]
created_at: "2026-05-25"
updated_at: "2026-05-25"
last_editor: "Standard Codebase"
discussion_hub: "https://discord.gg/example"
---

# Welcome to the Knowledge Node

This is the introductory guide and landing page for this Project R.E.D. node. To navigate to other guides in the repository, click the **☰ Menu Icon** in the top-left corner to open the side panel.

## How to Format Guides

If you are contributing to this node, every markdown file must strictly follow the metadata structure demonstrated at the top of this file. 

### Critical Formatting Rules:
- **No leading spaces:** The opening `---` must be on the very first line of the file.
- **The mandatory blank line:** There must be a completely **blank line** immediately following the closing `---` before any markdown content begins.
- **Required YAML fields:** `title`, `authors`, and `created_at` are mandatory.
- **Optional YAML fields:** `contributors`, `updated_at`, `last_editor`, and `discussion_hub`.
- **List formatting:** `authors` and `contributors` must be formatted as bracketed arrays (e.g., `["Name 1", "Name 2"]`).
- **File extension:** Save your files strictly as `.md`.

## Troubleshooting Rendering Issues

If a newly uploaded guide fails to display content or the metadata ribbon breaks, verify the following:

1. **Syntax Typos:** Check for extra spaces before the colons in your YAML (e.g., use `authors:` instead of `authors :`).
2. **Deprecated Fields:** Ensure you are not using outdated keys like `author_identity`. Use the `authors` list format instead.
3. **URL Pathing:** Ensure you are accessing the URL without the file extension (e.g., `http://localhost/guides/test-guide`, not `test-guide.md`).
4. **Browser Cache:** Perform a hard refresh (`Ctrl+Shift+R` or `Cmd+Shift+R`) to clear cached CSS/HTML.
5. **Engine Logs:** Check the server output by running `docker logs red_light_node` to spot any `goldmark` or YAML parsing errors.
6. **Container Sync:** If you updated the Go engine recently, ensure you run `docker-compose build --no-cache && docker-compose up -d`.
