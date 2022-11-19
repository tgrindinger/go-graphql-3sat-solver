package repositories

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	u "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

var dbName string = "testJobs.db"

func init() {
	os.Remove(dbName)
}

func TestInsertJob(t *testing.T) {
	cases := []struct {
		desc string
		job *model.Job
	}{
		{ "no clauses", jobWithoutClauses(u.New()) },
		{ "one clause", jobWithOneClause(u.New()) },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// arrange
			sut := NewSqliteJobRepository(dbName)

			// act
			err := sut.InsertJob(tc.job)

			// assert
			if err != nil {
				t.Fatalf("failed to insert job: %v", err)
			}
			verifyJobRow(t, tc.job)
			verifyClauseRows(t, tc.job)
		})
	}
}

func TestFindJob(t *testing.T) {
	cases := []struct {
		desc string
		want *model.Job
	}{
		{ "no clauses", jobWithoutClauses(u.New()) },
		{ "one clause", jobWithOneClause(u.New()) },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// arrange
			sut := NewSqliteJobRepository(dbName)
			err := sut.InsertJob(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			// act
			got, err := sut.FindJob(tc.want.Uuid)

			// assert
			if err != nil {
				t.Fatalf("failed to find job: %v", err)
			}
			verifyJobsAreEqual(t, got, tc.want)
		})
	}
}

func verifyJobsAreEqual(t testing.TB, got *model.Job, want *model.Job) {
	if got.Uuid != want.Uuid || got.Done != want.Done || got.Name != want.Name {
		t.Fatalf("got (%s %t %s) want (%s %t %s)", got.Uuid.String(), got.Done, got.Name, want.Uuid.String(), want.Done, want.Name)
	}
	if len(got.Clauses) != len(want.Clauses) {
		t.Fatalf("wrong number of clauses: got %d want %d", len(got.Clauses), len(want.Clauses))
	}
	for index, clause := range want.Clauses {
		gotClause := got.Clauses[index]
		if gotClause.Var1.Name != clause.Var1.Name || gotClause.Var1.Negated != clause.Var1.Negated ||
				gotClause.Var2.Name != clause.Var2.Name || gotClause.Var2.Negated != clause.Var2.Negated ||
				gotClause.Var3.Name != clause.Var3.Name || gotClause.Var3.Negated != clause.Var3.Negated {
			t.Errorf("clause %d got (%s %t %s %t %s %t) want (%s %t %s %t %s %t)", index,
				gotClause.Var1.Name, gotClause.Var1.Negated, gotClause.Var2.Name, gotClause.Var2.Negated, gotClause.Var3.Name, gotClause.Var3.Negated,
				clause.Var1.Name, clause.Var1.Negated, clause.Var2.Name, clause.Var2.Negated, clause.Var3.Name, clause.Var3.Negated)
		}
	}
}

func verifyJobRow(t testing.TB, job *model.Job) {
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

func verifyClauseRows(t testing.TB, job *model.Job) {
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
