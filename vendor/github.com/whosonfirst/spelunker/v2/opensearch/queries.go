package opensearch

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/whosonfirst/spelunker/v2"
)

// Something something something do all of this with templates...

func (s *OpenSearchSpelunker) matchAllQuery() string {
	return `{"query": { "match_all": {} }}`
}

func (s *OpenSearchSpelunker) idQuery(id int64) string {
	return fmt.Sprintf(`{"query": { "ids": { "values": [ %d ] } } }`, id)
}

// Null Island

func (s *OpenSearchSpelunker) visitingNullIslandQuery(filters []spelunker.Filter) string {

	q := s.visitingNullIslandQueryCriteria(filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s *OpenSearchSpelunker) visitingNullIslandFacetedQuery(filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.visitingNullIslandQueryCriteria(filters)
	str_aggs := s.facetsToAggregations(facets)

	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s *OpenSearchSpelunker) visitingNullIslandQueryCriteria(filters []spelunker.Filter) string {

	terms := []string{
		`{ "term": { "geom:latitude":  0.0  } }`,
		`{ "term": { "geom:longitude":  0.0  } }`,
	}

	str_terms := strings.Join(terms, ",")

	q := fmt.Sprintf(`{ "bool": { "must": [ %s ] } }`, str_terms)

	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

// Descendants

func (s *OpenSearchSpelunker) descendantsQuery(id int64, filters []spelunker.Filter) string {

	q := s.descendantsQueryCriteria(id, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s *OpenSearchSpelunker) descendantsFacetedQuery(id int64, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.descendantsQueryCriteria(id, filters)
	str_aggs := s.facetsToAggregations(facets)

	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s *OpenSearchSpelunker) descendantsQueryCriteria(id int64, filters []spelunker.Filter) string {

	q := fmt.Sprintf(`{ "term": { "wof:belongsto":  %d  } }`, id)

	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s *OpenSearchSpelunker) hasPlacetypeQuery(pt string, filters []spelunker.Filter) string {

	q := s.hasPlacetypeQueryCriteria(pt, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s *OpenSearchSpelunker) hasPlacetypeFacetedQuery(pt string, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.hasPlacetypeQueryCriteria(pt, filters)
	str_aggs := s.facetsToAggregations(facets)

	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s *OpenSearchSpelunker) hasPlacetypeQueryCriteria(pt string, filters []spelunker.Filter) string {

	q := fmt.Sprintf(`{ "term": { "wof:placetype":  "%s"  } }`, pt)

	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)
}

func (s *OpenSearchSpelunker) hasAlternatePlacetypeQuery(pt string, filters []spelunker.Filter) string {

	q := s.hasAlternatePlacetypeQueryCriteria(pt, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s *OpenSearchSpelunker) hasAlternatePlacetypeFacetedQuery(pt string, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.hasAlternatePlacetypeQueryCriteria(pt, filters)
	str_aggs := s.facetsToAggregations(facets)

	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s *OpenSearchSpelunker) hasAlternatePlacetypeQueryCriteria(pt string, filters []spelunker.Filter) string {

	q := fmt.Sprintf(`{ "term": { "wof:placetype_alt":  "%s"  } }`, pt)

	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)
}

func (s *OpenSearchSpelunker) hasConcordanceQuery(namespace string, predicate string, value any, filters []spelunker.Filter) string {

	q := s.hasConcordanceQueryCriteria(namespace, predicate, value, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s *OpenSearchSpelunker) hasConcordanceFacetedQuery(namespace string, predicate string, value any, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.hasConcordanceQueryCriteria(namespace, predicate, value, filters)
	str_aggs := s.facetsToAggregations(facets)

	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s *OpenSearchSpelunker) hasConcordanceQueryCriteria(namespace string, predicate string, value any, filters []spelunker.Filter) string {

	var q string

	str_value := fmt.Sprintf("%v", value)

	// Basically we need to index "magic 8"s...

	switch {
	case namespace != "" && predicate != "" && str_value != "":
		q = fmt.Sprintf(`{ "term": { "wof:concordances.%s:%s":  { "value": "%s", "case_insensitive": true } } }`, namespace, predicate, str_value)
	case namespace != "" && predicate != "":
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags.keyword":  { "value": "%s:%s=*", "case_insensitive": true }  } }`, namespace, predicate)
	case predicate != "" && str_value != "":
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags.keyword":  { "value": "*:%s=%s", "case_insensitive": true }  } }`, predicate, value)
	case namespace != "" && str_value != "":
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags.keyword":  { "value": "%s:*=%s", "case_insensitive": true }  } }`, namespace, value)
	case namespace != "":
		q = fmt.Sprintf(`{ "prefix": { "wof:concordances_machinetags.keyword":  { "value": "%s:*", "case_insensitive": true }  } }`, namespace)
	case predicate != "":
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags.keyword":  { "value": "*:%s", "case_insensitive": true }  } }`, predicate)
	case value != nil:
		q = fmt.Sprintf(`{ "wildcard": { "wof:concordances_machinetags.keyword":  { "value": "*:*=%s", "case_insensitive": true }  } }`, value)
	default:

	}

	slog.Info("Concordance", "namespace", namespace, "predicate", predicate, "value", value)
	slog.Info(q)

	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)
}

func (s *OpenSearchSpelunker) getRecentQuery(d time.Duration, filters []spelunker.Filter) string {

	q := s.getRecentQueryCriteria(d, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s *OpenSearchSpelunker) getRecentFacetedQuery(d time.Duration, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.getRecentQueryCriteria(d, filters)
	str_aggs := s.facetsToAggregations(facets)

	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s *OpenSearchSpelunker) getRecentQueryCriteria(d time.Duration, filters []spelunker.Filter) string {

	now := time.Now()
	ts := now.Unix()

	then := ts - int64(d.Seconds())

	q := fmt.Sprintf(`{ "range": { "wof:lastmodified": { "gte": %d  } } }`, then)

	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)
}

func (s *OpenSearchSpelunker) matchAllFacetedQuery(facets []*spelunker.Facet) string {

	str_aggs := s.facetsToAggregations(facets)
	return fmt.Sprintf(`{"query": { "match_all": {} }, "aggs": { %s } }`, str_aggs)
}

// https://opensearch.org/docs/latest/aggregations/
// https://opensearch.org/docs/latest/aggregations/bucket/terms/

func (s *OpenSearchSpelunker) searchQuery(search_opts *spelunker.SearchOptions, filters []spelunker.Filter) string {

	q := s.searchQueryCriteria(search_opts, filters)
	return fmt.Sprintf(`{"query": %s  }`, q)
}

func (s *OpenSearchSpelunker) searchFacetedQuery(search_opts *spelunker.SearchOptions, filters []spelunker.Filter, facets []*spelunker.Facet) string {

	q := s.searchQueryCriteria(search_opts, filters)
	str_aggs := s.facetsToAggregations(facets)

	slog.Info(str_aggs)
	return fmt.Sprintf(`{"query": %s, "aggs": { %s } }`, q, str_aggs)
}

func (s *OpenSearchSpelunker) searchQueryCriteria(search_opts *spelunker.SearchOptions, filters []spelunker.Filter) string {

	// This is a short-term fix to address these issues:
	// https://github.com/whosonfirst/spelunker/v2-opensearch/issues/4
	// https://github.com/whosonfirst/spelunker/v2-httpd/issues/20
	// In advance of addressing the actual problem:
	// https://github.com/whosonfirst/whosonfirst-opensearch/issues/2
	lower_q := strings.ToLower(search_opts.Query)

	// https://github.com/whosonfirst/spelunker/v2-opensearch/issues/6
	// switch to https://opensearch.org/docs/latest/query-dsl/full-text/query-string/
	// https://opensearch.org/docs/latest/query-dsl/full-text/simple-query-string/

	q := fmt.Sprintf(`{ "simple_query_string": { "query": "%s", "fields": ["search"], "default_operator": "AND" } }`, lower_q)

	if len(filters) == 0 {
		return q
	}

	must := []string{
		q,
	}

	return s.mustQueryWithFiltersCriteria(must, filters)
	return fmt.Sprintf(`{"query": %s }`, q)
}

func (s *OpenSearchSpelunker) facetsToAggregations(facets []*spelunker.Facet) string {

	count_facets := len(facets)
	aggs := make([]string, count_facets)

	for i, f := range facets {

		var facet_field string

		switch f.String() {
		case "isdeprecated":
			// This flag is derived in go-whosonfirst-spelunker/document and
			// assigned in go-whosonfirst-opensearch
			facet_field = "mz:is_deprecated"
		case "iscurrent":
			facet_field = "mz:is_current"
		case "placetypealt":
			facet_field = "wof:placetype_alt"
		default:
			facet_field = fmt.Sprintf("wof:%s", f)
		}

		aggs[i] = fmt.Sprintf(`"%s": { "terms": { "field": "%s", "size": 1000 } }`, f, facet_field)
	}

	return strings.Join(aggs, ",")
}

func (s *OpenSearchSpelunker) mustQueryWithFiltersCriteria(must []string, filters []spelunker.Filter) string {

	for _, f := range filters {

		switch f.Scheme() {
		case "placetype":
			must = append(must, fmt.Sprintf(`{ "term": { "wof:placetype": "%s" } }`, f.Value()))
		case "placetypealt":
			must = append(must, fmt.Sprintf(`{ "term": { "wof:placetype_alt": "%s" } }`, f.Value()))
		case "country":
			must = append(must, fmt.Sprintf(`{ "term": { "wof:country": "%s" } }`, f.Value()))
		case "iscurrent":
			must = append(must, fmt.Sprintf(`{ "term": { "mz:is_current": "%d" } }`, f.Value()))
		case "isdeprecated":
			must = append(must, fmt.Sprintf(`{ "term": { "mz:is_deprecated": "%d" } }`, f.Value()))
		default:
			slog.Warn("Unsupported filter scheme", "scheme", f.Scheme())
		}
	}

	str_must := strings.Join(must, ",")

	return fmt.Sprintf(`{ "bool": { "must": [ %s ] } }`, str_must)
}
