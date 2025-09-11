# Browsers MCP server

An MCP server provideing read-only access to browsers configuration files: profiles, bookmarks, history.

Supported browsers: Chrome, Firefox.

## Tools

### list_browsers

List discovered browsers.

### list_profiles

List profiles for a specific browser.

Parameters: 
- `browser` (`string`, optional): the browser name (as returned by `list_browsers`). Required if `list_browsers` returns several browsers.

### list_bookmarks

List bookmarks for a given profile of a given browser.

Parameters:
- `browser` (`string`, optional): the browser name (as returned by `list_browsers`). Required if `list_browsers` returns several browsers.
- `profile` (`string`, optional): the profile name (as returned by `list_profiles`). Required if `list_profiles` returns several profiles for the given browser.

### list_search_engine_queries

List the queries in search engines (supported search engines: Google).

Parameters:
- `browser` (`string`, optional): the browser name (as returned by `list_browsers`). Required if `list_browsers` returns several browsers.
- `profile` (`string`, optional): the profile name (as returned by `list_profiles`). Required if `list_profiles` returns several profiles for the given browser.
- `start_time` (`string`, format `YYYY-MM-DD HH:MM:SS`, optional): list the search engine queries from this time, default is today at midnight.
- `limit` (`number`, optional): the number of results to return, default is 10.

### list_visited_pages_from_search_engine_query

List the pages visited from a search engine query.

Parameters:
- `browser` (`string`, optional): the browser name (as returned by `list_browsers`). Required if `list_browsers` returns several browsers.
- `profile` (`string`, optional): the profile name (as returned by `list_profiles`). Required if `list_profiles` returns several profiles for the given browser.
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
