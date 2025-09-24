# Browsers MCP server

An MCP server provideing read-only access to browsers configuration files: profiles, bookmarks, history.

Supported browsers: Chrome, Firefox, Safari (default profile only).

## Tools

### list_bookmarks

List bookmarks for a given profile of a given browser.

Parameters:
- `profile` (`string`): the profile name (as indicated in the description of the parameter). Provided only if several browsers or several profiles.

### list_search_engine_queries

List the queries in search engines (supported search engines: Google).

Parameters:
- `profile` (`string`): the profile name (as indicated in the description of the parameter). Provided only if several browsers or several profiles.
- `day` (`string`, format `YYYY-MM-DD`, optional): list the search engine queries during this day, default is today.
- `limit` (`number`, optional): the number of results to return, default is 10.

### list_visited_pages_from_search_engine_query

List the pages visited from a search engine query.

Not supported for Safari browser.

- `profile` (`string`): the profile name (as indicated in the description of the parameter). Provided only if several browsers or several profiles.
- `query` (`string`, required): the query string to list the visited pages for.
- `day` (`string`, format `YYYY-MM-DD`, optional): list the search engine queries during this day, default is today.

### list_source_repos_visits

List the pages visited in sources repositories.

Parameters:
- `profile` (`string`): the profile name (as indicated in the description of the parameter). Provided only if several browsers or several profiles.
- `day` (`string`, format `YYYY-MM-DD`, optional): list the visits during this day, default is today.
- `type` (`string`): Type of pages to list (`provider home`, `organization home`, `repository home`, `issues list`, `pull requests list`, `discussions list`, `issue`, `pull request`, `discussion`)

## Getting Started


### Claude Desktop

If you have npm installed, this is the fastest way to get started with `browsers-mcp-server` on Claude Desktop.

Open your `claude_desktop_config.json` and add the mcp server to the list of `mcpServers`:
``` json
{
  "mcpServers": {
    "browsers": {
      "command": "npx",
      "args": [
        "-y",
        "browsers-mcp-server@latest"
      ]
    }
  }
}
```


### Cursor

Install the extension manually by editing the `mcp.json` file:

```json
{
  "mcpServers": {
    "browsers-mcp-server": {
      "command": "npx",
      "args": ["-y", "browsers-mcp-server@latest"]
    }
  }
}
```

## Troubleshooting

You can output logs to a specific file with the `--log-file` flag, and indicate the log level with `--log-level=debug|info|warn|error` (default `warn`). By default, no logs are written.

