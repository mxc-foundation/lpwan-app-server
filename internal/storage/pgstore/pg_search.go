package pgstore

import (
	"context"
	"database/sql"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq/hstore"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/user"
)

// GlobalSearch performs a search on organizations, applications, gateways
// and devices.
func (ps *PgStore) GlobalSearch(ctx context.Context, userID int64, globalAdmin bool, search string, limit, offset int) ([]user.SearchResult, error) {
	var result []user.SearchResult

	query, tags := parseSearchQuery(search)
	query = "%" + search + "%"

	tagsHstore := hstore.Hstore{
		Map: make(map[string]sql.NullString),
	}
	for k, v := range tags {
		tagsHstore.Map[k] = sql.NullString{String: v, Valid: true}
	}

	err := sqlx.SelectContext(ctx, ps.db, &result, `
		select
			'device' as kind,
			greatest(similarity(d.name, $1), similarity(encode(d.dev_eui, 'hex'), $1)) as score,
			o.id as organization_id,
			o.name as organization_name,
			a.id as application_id,
			a.name as application_name,
			d.dev_eui as device_dev_eui,
			d.name as device_name,
			null as gateway_mac,
			null as gateway_name
		from device d
		inner join application a
			on a.id = d.application_id
		inner join organization o
			on o.id = a.organization_id
		left join organization_user ou
			on ou.organization_id = o.id
		left join "user" u
			on u.id = ou.user_id
		where
			($3 = true or u.id = $4)
			and (d.name ilike $2 or encode(d.dev_eui, 'hex') ilike $2 or ($7 != hstore('') and d.tags @> $7))
		union
		select
			'gateway' as kind,
			greatest(similarity(g.name, $1), similarity(encode(g.mac, 'hex'), $1)) as score,
			o.id as organization_id,
			o.name as organization_name,
			null as application_id,
			null as application_name,
			null as device_dev_eui,
			null as device_name,
			g.mac as gateway_mac,
			g.name as gateway_name
		from
			gateway g
		inner join organization o
			on o.id = g.organization_id
		left join organization_user ou
			on ou.organization_id = o.id
		left join "user" u
			on u.id = ou.user_id
		where
			($3 = true or u.id = $4)
			and (g.name ilike $2 or encode(g.mac, 'hex') ilike $2 or ($7 != hstore('') and g.tags @> $7))
		union
		select
			'organization' as kind,
			similarity(o.name, $1) as score,
			o.id as organization_id,
			o.name as organization_name,
			null as application_id,
			null as application_name,
			null as device_dev_eui,
			null as device_name,
			null as gateway_mac,
			null as gateway_name
		from
			organization o
		left join organization_user ou
			on ou.organization_id = o.id
		left join "user" u
			on u.id = ou.user_id
		where
			($3 = true or u.id = $4)
			and o.name ilike $2
		union
		select
			'application' as kind,
			similarity(a.name, $1) as score,
			o.id as organization_id,
			o.name as organization_name,
			a.id as application_id,
			a.name as application_name,
			null as device_dev_eui,
			null as device_name,
			null as gateway_mac,
			null as gateway_name
		from
			application a
		inner join organization o
			on o.id = a.organization_id
		left join organization_user ou
			on ou.organization_id = o.id
		left join "user" u
			on u.id = ou.user_id
		where
			($3 = true or u.id = $4)
			and a.name ilike $2
		order by
			score desc
		limit $5
		offset $6`,
		search,
		query,
		globalAdmin,
		userID,
		limit,
		offset,
		tagsHstore,
	)
	if err != nil {
		return nil, handlePSQLError(Select, err, "select error")
	}

	return result, nil
}

var searchTagRegexp = regexp.MustCompile(`([^ ]+):([^ ]+)`)

// parseSearchQuery returns the query and tags.
// Example: "foo bar:test" will return the query "foo" with tags {"bar": "test"}.
func parseSearchQuery(query string) (string, map[string]string) {
	matches := searchTagRegexp.FindAllStringSubmatch(query, -1)
	if len(matches) == 0 {
		return query, nil
	}

	tags := make(map[string]string)
	for _, t := range matches {
		if len(t) != 3 {
			continue
		}

		tags[t[1]] = t[2]
	}

	query = strings.TrimSpace(searchTagRegexp.ReplaceAllString(query, ""))

	return query, tags
}
