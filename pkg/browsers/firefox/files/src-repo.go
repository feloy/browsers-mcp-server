package files

import (
	"github.com/charmbracelet/log"
	"github.com/feloy/browsers-mcp-server/pkg/api"
)

func ListVisitedPagesFromSourceRepos(profile string, isRelative bool, options api.ListVisitedPagesFromSourceReposOptions) ([]api.VisitedPageFromSourceRepos, error) {
	log.Debug("source repository visits", "profile", profile, "options", options)

	type queryResult struct {
		Times        int
		URL          string
		Organization string
		Repository   string
		Pagetype     string
		Name         string
	}

	db, err := getDb(profile, isRelative)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	startTime := toDbDate(options.StartTime)
	endTime := toDbDate(options.EndTime)
	rows, err := db.Query(`with recursive 
  cte0 (title, pathAndQuery) as (
    SELECT 
      title as title,
      SUBSTR(url, 20, INSTR(url||'#', '#')-20) as pathAndQuery
    FROM moz_historyvisits hv
    INNER JOIN moz_places p ON p.id = hv.place_id 
    WHERE url LIKE 'https://github.com/%'
    AND url NOT LIKE 'https://github.com/search?%'
  	AND visit_date >= ?
	  AND visit_date < ?
  ),
  cte1 (title, path) AS (
    SELECT 
      title,
      SUBSTR(pathAndQuery, 1, INSTR(pathAndQuery||'?', '?') - 1) AS path
	  FROM cte0
  ),
  cte2 (title, path, organization, rest2) as (
    select 
      title,
      path,
      SUBSTR(path, 1, INSTR(path||'/', '/') - 1) as organization2,
      SUBSTR(path, INSTR(path||'/', '/') + 1) as rest2
    from cte1
  ),
  cte3 (title, path, organization, repository, rest3) as (
    select 
      title,
      path,
      organization,
      SUBSTR(rest2, 1, INSTR(rest2||'/', '/') - 1) as repository,
      SUBSTR(rest2, INSTR(rest2||'/', '/') + 1) as rest3
    from cte2
  ),
  cte4 (title, path, organization, repository, page, rest4) as (
    select 
      title,
      path,
      organization,
      repository,
      SUBSTR(rest3, 1, INSTR(rest3||'/', '/') - 1) as page,
      SUBSTR(rest3, INSTR(rest3||'/', '/') + 1) as rest4
    from cte3
  ),
  cte5 (title, path, organization, repository, page, name, rest5) as (
    select 
      title,
      path,
      organization,
      repository,
      page,
      SUBSTR(rest4, 1, INSTR(rest4||'/', '/') - 1) as name,
      SUBSTR(rest4, INSTR(rest4||'/', '/') + 1) as rest5
    from cte4
  ),
  cte6 (title, url, organization, repository, pagetype, name, rest5) as (
    select 
      title,
      'https://github.com/' || path,
      organization,
      repository,
      case 
        when organization = '' then 'provider home'
        when repository = '' then 'organization home'
        when page = '' then 'repository home'
        when page = 'issues' and name = '' then 'issues list'
        when page = 'pulls' and name = '' then 'pull requests list'
        when page = 'discussions' and name = '' then 'discussions list'
        when page = 'issues' and name != '' then 'issue'
        when page = 'pull' and name != '' then 'pull request'
        when page = 'discussions' and name != '' then 'discussion'
        else 'other details'  
      end as pagetype,
      name,
      rest5
    from cte5
  )
select count(*) as c, url, organization, repository, pagetype, name from cte6
where (? = '' OR ? = pagetype) AND pagetype != 'other details'
group by url, organization, repository, pagetype, name
order by c desc;
`, startTime, endTime, options.Type, options.Type)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visitedPages []api.VisitedPageFromSourceRepos
	for rows.Next() {
		var queryResult queryResult
		err = rows.Scan(&queryResult.Times, &queryResult.URL, &queryResult.Organization, &queryResult.Repository, &queryResult.Pagetype, &queryResult.Name)
		if err != nil {
			return nil, err
		}

		var namePtr *string
		if queryResult.Name != "" {
			namePtr = &queryResult.Name
		}
		visitedPages = append(visitedPages, api.VisitedPageFromSourceRepos{
			Times:        queryResult.Times,
			Provider:     "github",
			URL:          queryResult.URL,
			Organization: queryResult.Organization,
			Repository:   queryResult.Repository,
			Type:         api.SourceRepoPageType(queryResult.Pagetype),
			Number:       namePtr,
		})
	}
	return visitedPages, nil
}
