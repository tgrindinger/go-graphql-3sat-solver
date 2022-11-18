package repositories

import (
	"database/sql"
	"fmt"
	"testing"

	u "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

func TestInsert(t *testing.T) {
	cases := []struct {
		desc string
		jobFunc func (uuid u.UUID) *model.Job
	}{
		{ "no clauses", jobWithoutClauses },
		{ "one clause", jobWithOneClause },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// arrange
			dbName := "testJobs.db"
			uuid := u.New()
			job := tc.jobFunc(uuid)
			sut := NewSqliteJobRepository("testJobs.db")

			// act
			err := sut.InsertJob(job)

			// assert
			if err != nil {
				t.Fatalf("failed to insert job: %v", err)
			}
			verifyJobRow(t, job, dbName)
			verifyClauseRows(t, job, dbName)
		})
	}
}

func verifyJobRow(t testing.TB, job *model.Job, dbName string) {
	db, _ := sql.Open("sqlite3", dbName)
	defer db.Close()
	rows, _ := db.Query("SELECT uuid, done, name FROM jobs WHERE uuid = ?", job.Uuid.String())
	hasRow := rows.Next()
	if !hasRow {
		t.Fatal("returned 0 rows")
	}
	got := &model.Job{}
	var gotUuid string
	err := rows.Scan(&gotUuid, &got.Done, &got.Name)
	if err != nil {
		t.Fatalf("unable to read row: %v", err)
	}
	hasRow = rows.Next()
	if hasRow {
		t.Fatal("returned more than 1 row")
	}
	if gotUuid != job.Uuid.String() || got.Done != job.Done || got.Name != job.Name {
		t.Fatalf("got (%s, %t, %s) want (%s, %t, %s)", gotUuid, got.Done, got.Name, job.Uuid.String(), job.Done, job.Name)
	}
}

func verifyClauseRows(t testing.TB, job *model.Job, dbName string) {
	db, _ := sql.Open("sqlite3", dbName)
	defer db.Close()
	rows, _ := db.Query("SELECT var1, var1negated, var2, var2negated, var3, var3negated FROM clauses WHERE uuid = ?", job.Uuid.String())
	hasRow := rows.Next()
	for _, clause := range job.Clauses {
		if !hasRow {
			t.Fatal("not enough clause rows")
		}
		verifyClauseRow(t, clause, rows)
		hasRow = rows.Next()
	}
	if hasRow {
		t.Fatal("returned too many clause rows")
	}
}

func verifyClauseRow(t testing.TB, clause *model.Clause, rows *sql.Rows) {
	got := &model.Clause{
		Var1: &model.Variable{},
		Var2: &model.Variable{},
		Var3: &model.Variable{},
	}
	err := rows.Scan(&got.Var1.Name, &got.Var1.Negated, &got.Var2.Name, &got.Var2.Negated, &got.Var3.Name, &got.Var3.Negated)
	if err != nil {
		t.Fatalf("unable to read row: %v", err)
	}
	if got.Var1.Name != clause.Var1.Name || got.Var1.Negated != clause.Var1.Negated ||
			got.Var2.Name != clause.Var2.Name || got.Var2.Negated != clause.Var2.Negated ||
			got.Var3.Name != clause.Var3.Name || got.Var3.Negated != clause.Var3.Negated {
		t.Fatalf("got (%s, %t, %s, %t, %s, %t) want (%s, %t, %s, %t, %s, %t)",
			got.Var1.Name, got.Var1.Negated, got.Var2.Name, got.Var2.Negated, got.Var3.Name, got.Var3.Negated,
			clause.Var1.Name, clause.Var1.Negated, clause.Var2.Name, clause.Var2.Negated, clause.Var3.Name, clause.Var3.Negated)
	}
}

func jobWithoutClauses(uuid u.UUID) *model.Job {
	return &model.Job{
		Uuid: uuid,
		Name: fmt.Sprintf("test-%s", uuid.String()),
		Clauses: []*model.Clause{},
	}
}

func jobWithOneClause(uuid u.UUID) *model.Job {
	return &model.Job{
		Uuid: uuid,
		Name: fmt.Sprintf("test-%s", uuid.String()),
		Clauses: []*model.Clause{ {
				Var1: &model.Variable{
					Name: "v1",
					Negated: true,
				},
				Var2: &model.Variable{
					Name: "v2",
					Negated: false,
				},
				Var3: &model.Variable{
					Name: "v3",
					Negated: true,
				},
			},
		},
	}
}
