package repositories

import (
	"database/sql"
	"fmt"

	u "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type SqliteJobRepository struct {
	firstTime bool
	dbName string
}

func NewSqliteJobRepository(dbName string) *SqliteJobRepository {
	return &SqliteJobRepository{
		firstTime: true,
		dbName: dbName,
	}
}

func (r* SqliteJobRepository) FindJob(uuid u.UUID) (*model.Job, error) {
	db := r.openDatabase()
	job, err := r.queryJob(uuid, db)
	if err != nil {
		return nil, err
	}
	job.Clauses, err = r.queryClauses(uuid, db)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r* SqliteJobRepository) InsertJob(job *model.Job) error {
	db := r.openDatabase()
	defer db.Close()
	err := r.insertJobRow(job, db)
	if err != nil {
		return err
	}
	err = r.insertClauseRows(job, db)
	if err != nil {
		return err
	}
	return nil
}

func (r* SqliteJobRepository) MarkDone(job *model.Job) error {
	db := r.openDatabase()
	statement, err := db.Prepare("UPDATE jobs SET done = ?")
	if err != nil {
		return err
	}
	_, err = statement.Exec(job.Done)
	return err
}

func (r* SqliteJobRepository) insertJobRow(job *model.Job, db *sql.DB) error {
	statement, err := db.Prepare("INSERT INTO jobs (uuid, done, name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = statement.Exec(job.Uuid, job.Done, job.Name)
	return err
}

func (r* SqliteJobRepository) insertClauseRows(job *model.Job, db *sql.DB) error {
	for _, clause := range job.Clauses {
		statement, err := db.Prepare("INSERT INTO clauses (uuid, var1, var1negated, var2, var2negated, var3, var3negated) VALUES (?, ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return err
		}
		_, err = statement.Exec(job.Uuid.String(), clause.Var1.Name, clause.Var1.Negated, clause.Var2.Name, clause.Var2.Negated, clause.Var3.Name, clause.Var3.Negated)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r* SqliteJobRepository) queryJob(uuid u.UUID, db *sql.DB) (*model.Job, error) {
	jobRow, err := db.Query("SELECT uuid, done, name FROM jobs where uuid = ?", uuid.String())
	if err != nil {
		return nil, err
	}
	job := &model.Job{}
	found := jobRow.Next()
	if !found {
		return nil, fmt.Errorf("unable to find job with uuid %s", uuid.String())
	}
	jobRow.Scan(&job.Uuid, &job.Done, &job.Name)
	return job, nil
}

func (r* SqliteJobRepository) queryClauses(uuid u.UUID, db *sql.DB) ([]*model.Clause, error) {
	clauseRows, err := db.Query("SELECT var1, var1negated, var2, var2negated, var3, var3negated FROM clauses WHERE UUID = ?", uuid.String())
	if err != nil {
		return nil, err
	}
	clauses := []*model.Clause{}
	for clauseRows.Next() {
		clause := &model.Clause{
			Var1: &model.Variable{},
			Var2: &model.Variable{},
			Var3: &model.Variable{},
		}
		clauseRows.Scan(&clause.Var1.Name, &clause.Var1.Negated, &clause.Var2.Name, &clause.Var2.Negated, &clause.Var3.Name, &clause.Var3.Negated)
		clauses = append(clauses, clause)
	}
	return clauses, nil
}

func (r *SqliteJobRepository) openDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", r.dbName)
	if err != nil {
		panic(fmt.Sprintf("Unable to open jobs database: %v", err))
	}
	if r.firstTime {
		r.initJobsTable(db)
		r.initClausesTable(db)
		r.firstTime = false
	}
	return db
}

func (r *SqliteJobRepository) initJobsTable(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS jobs (id INTEGER PRIMARY KEY, uuid STRING, done BOOLEAN, name STRING)")
	if err != nil {
		panic(fmt.Sprintf("Unable to create jobs table: %v", err))
	}
	statement.Exec()
	statement.Close()
}

func (r *SqliteJobRepository) initClausesTable(db *sql.DB) {
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS clauses (id INTEGER PRIMARY KEY, uuid STRING, var1 STRING, var1negated BOOLEAN, var2 STRING, var2negated BOOLEAN, var3 STRING, var3negated BOOLEAN)")
	if err != nil {
		panic(fmt.Sprintf("Unable to create clauses table: %v", err))
	}
	statement.Exec()
	statement.Close()
}
