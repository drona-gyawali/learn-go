package sqllite

import (
	"database/sql"
	"fmt"

	"github.com/drona-gyawali/learn-go/internal/config"
	"github.com/drona-gyawali/learn-go/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqllite struct {
	Db *sql.DB
}


func New(cfg *config.Config) (*Sqllite, error) {
	db , err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	    id    INTEGER PRIMARY KEY AUTOINCREMENT,
	    name  TEXT,
	    email TEXT,
	    age   INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqllite{
		Db: db,
	}, nil
}

func (s *Sqllite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt , err := s.Db.Prepare("INSERT INTO STUDENTS (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	id , err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil 
} 


func (s *Sqllite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM STUDENTS WHERE ID = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("No student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
	
}



func (s *Sqllite) GetStudentList() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM STUDENTS")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil , err
	}

	defer rows.Close()

	var students  []types.Student

	for rows.Next() {
		var student  types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil , err
		}
		students = append(students, student)
	}
	return students, nil
}


func (s *Sqllite) UpdateStudent(data types.Student) (int64, error) {

	query := "UPDATE STUDENTS SET name = ?, email = ?, age = ? WHERE id = ?"
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return 0 , err
	}

	defer stmt.Close()
	var student types.Student
	res, err := stmt.Exec(data.Name, data.Email, data.Age)
	if err != nil {
		return 0, fmt.Errorf("No query found")
	}

	r , err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if r == 0 {
		return 0, fmt.Errorf("no student found with ID %d", data.Id)
	}

	// _id, err := strconv.ParseInt(student.Id, 10, 64)
	// if err != nil {
	// 	return 0, err
	// }
	return student.Id , nil
}