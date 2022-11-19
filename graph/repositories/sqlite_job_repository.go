package repositories

import (
	"database/sql"
	"fmt"

	u "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tgrindinger/go-graphql-3sat-solver/graph/model"
)

type SqliteJobRepository struct {
	db *sql.DB
}

func NewSqliteJobRepository(dbName string) *SqliteJobRepository {
	repo := &SqliteJobRepository{}
	repo.openDatabase(dbName)
	return repo
}

func (r* SqliteJobRepository) FindJob(uuid u.UUID) (*model.Job, error) {
	job, err := r.queryJob(uuid)
	if err != nil {
		return nil, err
	}
	job.Clauses, err = r.queryClauses(uuid)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r* SqliteJobRepository) InsertJob(job *model.Job) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to create insert job transaction: %v", err)
	}
	err = r.insertJobRow(job, tx)
	if err != nil {
		return err
	}
	err = r.insertClauseRows(job, tx)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (r* SqliteJobRepository) MarkDone(job *model.Job) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("unable to create mark done transaction: %v", err)
	}
	statement, err := tx.Prepare("UPDATE jobs SET done = ? WHERE uuid = ?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(true, job.Uuid.String())
	tx.Commit()
	return err
}

func (r* SqliteJobRepository) insertJobRow(job *model.Job, tx *sql.Tx) error {
	statement, err := tx.Prepare("INSERT INTO jobs (uuid, done, name) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to create insert job statement: %v", err)
	}
	defer statement.Close()
	_, err = statement.Exec(job.Uuid, job.Done, job.Name)
	if err != nil {
		return fmt.Errorf("failed to execute insert job statement: %v", err)
	}
	return err
}

func (r* SqliteJobRepository) insertClauseRows(job *model.Job, tx *sql.Tx) error {
	statement, err := tx.Prepare("INSERT INTO clauses (uuid, var1, var1negated, var2, var2negated, var3, var3negated) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to create insert clause statement: %v", err)
	}
	defer statement.Close()
	for _, clause := range job.Clauses {
		_, err = statement.Exec(job.Uuid.String(), clause.Var1.Name, clause.Var1.Negated, clause.Var2.Name, clause.Var2.Negated, clause.Var3.Name, clause.Var3.Negated)
		if err != nil {
			return fmt.Errorf("failed to execute insert clause statement: %v", err)
		}
	}
	return nil
}

func (r* SqliteJobRepository) queryJob(uuid u.UUID) (*model.Job, error) {
	jobRow, err := r.db.Query("SELECT uuid, done, name FROM jobs where uuid = ?", uuid.String())
	if err != nil {
		errDesc := fmt.Errorf("failed to query job: %v", err)
		return nil, errDesc
	}
	defer jobRow.Close()
	job := &model.Job{}
	found := jobRow.Next()
	if !found {
		return nil, fmt.Errorf("unable to find job with uuid %s", uuid.String())
	}
	jobRow.Scan(&job.Uuid, &job.Done, &job.Name)
	return job, nil
}

func (r* SqliteJobRepository) queryClauses(uuid u.UUID) ([]*model.Clause, error) {
	clauseRows, err := r.db.Query("SELECT var1, var1negated, var2, var2negated, var3, var3negated FROM clauses WHERE UUID = ?", uuid.String())
	if err != nil {
		errDesc := fmt.Errorf("failed to query clauses: %v", err)
		return nil, errDesc
	}
	defer clauseRows.Close()
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

func (r *SqliteJobRepository) openDatabase(dbName string) {
	var err error
	r.db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		panic(fmt.Sprintf("Unable to open jobs database: %v", err))
	}
	r.initJobsTable()
	r.initClausesTable()
}

func (r *SqliteJobRepository) initJobsTable() {
	statement, err := r.db.Prepare("CREATE TABLE IF NOT EXISTS jobs (id INTEGER PRIMARY KEY, uuid STRING, done BOOLEAN, name STRING)")
	if err != nil {
		panic(fmt.Sprintf("Unable to create jobs table statement: %v", err))
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		panic(fmt.Sprintf("unable to execute create jobs table statement: %v", err))
	}
}

func (r *SqliteJobRepository) initClausesTable() {
	statement, err := r.db.Prepare("CREATE TABLE IF NOT EXISTS clauses (id INTEGER PRIMARY KEY, uuid STRING, var1 STRING, var1negated BOOLEAN, var2 STRING, var2negated BOOLEAN, var3 STRING, var3negated BOOLEAN)")
	if err != nil {
		panic(fmt.Sprintf("Unable to create clauses table statement: %v", err))
	}
	defer statement.Close()
	_, err = statement.Exec()
	if err != nil {
		panic(fmt.Sprintf("unable to execute create clauses table statement: %v", err))
	}
}
