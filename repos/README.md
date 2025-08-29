# wof

## wof repos

```
$> ./bin/wof repos -h
List repositories for the Who's On First (or other GitHub) organization.
Usage:
	 ./bin/wof [options]
  -ensure-commits
    	Ensure that 1 or more files have been updated in the last commit
  -exclude value
    	Exclude repositories with this prefix
  -exclude-archived
    	Exclude repos that have been archived.
  -forked
    	Only include repositories that have been forked
  -not-forked
    	Only include repositories that have not been forked
  -org string
    	The name of the organization to clone repositories from (default "whosonfirst-data")
  -prefix value
    	Limit repositories to only those with this prefix
  -token string
    	A valid GitHub API access token
  -updated-since string
    	A valid Unix timestamp or an ISO8601 duration string (months are currently not supported)
  -verbose
    	Enable verbose (debug) logging
```

### Examples

For example:

```
$> ./bin/wof repos -prefix whosonfirst-data-admin-
whosonfirst-data-admin-ad
whosonfirst-data-admin-ae
whosonfirst-data-admin-af
whosonfirst-data-admin-ag
whosonfirst-data-admin-ai
whosonfirst-data-admin-al
whosonfirst-data-admin-am
whosonfirst-data-admin-an
whosonfirst-data-admin-ao
whosonfirst-data-admin-aq
whosonfirst-data-admin-ar
whosonfirst-data-admin-as
whosonfirst-data-admin-at
whosonfirst-data-admin-au
whosonfirst-data-admin-aw
whosonfirst-data-admin-ax
... and so on
```
