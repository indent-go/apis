package pgquery_test

import (
	"testing"

	"go.indent.com/apis/pkg/providers/postgres/pgquery"
)

func TestInvalidSQL(t *testing.T) {
	_, err := pgquery.Analyze("9E2h7hkMHqguzf6KURjxLqvTG")
	if err == nil {
		t.Errorf("invalid sql: should have errored")
	}
}

func TestSimpleQuery(t *testing.T) {
	res, _ := pgquery.Analyze("select * from users;")
	if len(res.Actions) == 0 {
		t.Errorf("missing: AnalyzeResult.Actions")
	}
	if len(res.Tables) == 0 {
		t.Errorf("missing: AnalyzeResult.Tables")
	}
	if len(res.Columns) == 0 {
		t.Errorf("missing: AnalyzeResult.Columns")
	}
}

func TestOneColumn(t *testing.T) {
	res, _ := pgquery.Analyze("select username from users;")
	if len(res.Actions) == 0 {
		t.Errorf("missing: AnalyzeResult.Actions")
	}
	if len(res.Tables) == 0 {
		t.Errorf("missing: AnalyzeResult.Tables")
	}
	if len(res.Columns) == 0 {
		t.Errorf("missing: AnalyzeResult.Columns")
	}
}

func TestMultiColumn(t *testing.T) {
	t.Skip()
	pgquery.Analyze("select id, username, password from users;")
}

func TestSimpleJoin(t *testing.T) {
	t.Skip()
	pgquery.Analyze(`
		select id, username, password
			from users
			left outer join companies on (users.company_id = companies.id);
	`)
}
