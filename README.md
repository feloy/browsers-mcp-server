# Browsers MCP server

An MCP server provideing read-only access to browsers configuration files: profiles, bookmarks, history.

Supported browsers: Chrome, Firefox.

## Tools

### list_bookmarks, list_bookmarks_browserName

List bookmarks for a given profile of a given browser. `list_bookmarks` is provided if only one browser is found. A set of `list_bookmarks_browserName` are provided when several browsers are detected, one for each browser.

Parameters:
- `profile` (`string`): the profile name (as indicated in the description of the parameter). Provided only if the browser has several profiles.

### list_search_engine_queries, list_search_engine_queries_browserName

List the queries in search engines (supported search engines: Google). `list_search_engine_queries` is provided if only one browser is found. A set of `list_search_engine_queries_browserName` are provided when several browsers are detected, one for each browser.

Parameters:
- `profile` (`string`): the profile name (as indicated in the description of the parameter). Provided only if the browser has several profiles.
- `start_time` (`string`, format `YYYY-MM-DD HH:MM:SS`, optional): list the search engine queries from this time, default is today at midnight.
- `limit` (`number`, optional): the number of results to return, default is 10.

### list_visited_pages_from_search_engine_query, list_visited_pages_from_search_engine_query_browserName

List the pages visited from a search engine query. `list_visited_pages_from_search_engine_query` is provided if only one browser is found. A set of `list_visited_pages_from_search_engine_query_browserName` are provided when several browsers are detected, one for each browser.

- `profile` (`string`): the profile name (as indicated in the description of the parameter). Provided only if the browser has several profiles.
- `query` (`string`, required): the query string to list the visited pages for.
- `start_time` (`string`, format `YYYY-MM-DD HH:MM:SS`, optional): list the search engine queries from this time, default is today at midnight.


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

