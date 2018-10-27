package driver

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
	"github.com/wedancedalot/squirrel"
	"bonex-middleware/config"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type Logger interface {
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

type Mysql struct {
	Db        *sqlx.DB
	Logger    Logger
	DebugMode bool
}

var errorsRegexp = regexp.MustCompile(`^Error (?P<code>\d+)`)

var ErrDuplicate = fmt.Errorf("DB duplicate error")
var ErrNoRows = fmt.Errorf("DB no rows in resultset")

// CreateConnection creates a mysql connection
func CreateConnection(c *config.MysqlConfig, l Logger) (*Mysql, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", c.User, c.Password, c.Host, c.Port, c.Database))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpenConns)

	return &Mysql{db, l, c.DebugMode}, nil
}

// Migrate makes migration for DB from migrationsDir
func Migrate(c *config.MysqlConfig, migrationsDir string) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}

	dir := filepath.Join(filepath.Dir(ex), migrationsDir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		//return errors.New("Migrations dir does not exist: " + dir)
		dir = migrationsDir
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return errors.New("Migrations dir does not exist: " + dir)
		}
	}

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&multiStatements=true&parseTime=true", c.User, c.Password, c.Host, c.Port, c.Database))
	if err != nil {
		return err
	}

	migrations := &migrate.FileMigrationSource{
		Dir: dir,
	}

	_, err = migrate.Exec(db.DB, "mysql", migrations, migrate.Up)
	return err
}

// Find into dest from querybuilder
func (this *Mysql) Find(dest interface{}, b squirrel.SelectBuilder, tx ...*sqlx.Tx) error {
	q, params, err := b.ToSql()
	if err != nil {
		return err
	}

	start := time.Now()

	if len(tx) > 0 && tx[0] != nil {
		err = tx[0].Select(dest, q, params...)
	} else {
		err = this.Db.Select(dest, q, params...)
	}
	exTime := time.Since(start)
	if exTime > time.Millisecond * 100 {
		this.warn("long running query found!:", "query", q, "params", params, "execution time", exTime)
	} else if this.DebugMode {
		this.debug(q, "params", params, "execution time", exTime)

	}

	return err
}

// Find into dest from querybuilder
func (this *Mysql) FindRaw(dest interface{}, q string, params ...interface{}) error {
	start := time.Now()
	err := this.Db.Select(dest, q, params...)

	if this.DebugMode {
		this.debug(q, "params", params, "execution time", time.Since(start))
	}

	return err
}

func (this *Mysql) FindFirstRaw(dest interface{}, q string, params ...interface{}) error {
	start := time.Now()
	err := this.Db.Get(dest, q, params...)

	if this.DebugMode {
		this.debug(q, "params", params, "execution time", time.Since(start))
	}

	return err
}

// Find first row into dest from querybuilder
func (this *Mysql) FindFirst(dest interface{}, b squirrel.SelectBuilder, tx ...*sqlx.Tx) (err error) {
	q, params, err := b.ToSql()
	if err != nil {
		return
	}

	start := time.Now()

	if len(tx) > 0 && tx[0] != nil {
		err = tx[0].Get(dest, q, params...)
	} else {
		err = this.Db.Get(dest, q, params...)
	}

	if this.DebugMode {
		this.debug(q, "params", params, "execution time", time.Since(start))
	}

	return this.parseError(err)
}

// Insert from querybuilder
func (this *Mysql) Insert(b squirrel.InsertBuilder, tx ...*sqlx.Tx) (uint64, error) {
	q, args, err := b.ToSql()
	if err != nil {
		return 0, err
	}

	result, err := this.exec(q, args, tx...)
	if err != nil {
		return 0, this.parseError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), err
}

// Exec query
func (this *Mysql) Exec(q string, args []interface{}, err error, tx ...*sqlx.Tx) (uint64, error) {
	if err != nil {
		return 0, this.parseError(err)
	}

	result, err := this.exec(q, args, tx...)
	if err != nil {
		return 0, this.parseError(err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return 0, this.parseError(err)
	}

	return uint64(affectedRows), nil
}

func (this *Mysql) CallFunc(q string, params []interface{}, tx ...*sqlx.Tx) error {
	var err error

	start := time.Now()
	if len(tx) > 0 && tx[0] != nil {
		err = tx[0].QueryRow(q).Scan(params...)
	} else {
		err = this.Db.QueryRow(q).Scan(params...)
	}

	if this.DebugMode {
		this.debug(q, "execution time", time.Since(start))
	}

	return err
}

func (this *Mysql) exec(q string, params []interface{}, tx ...*sqlx.Tx) (sql.Result, error) {
	var result sql.Result
	var err error

	start := time.Now()

	if len(tx) > 0 && tx[0] != nil {
		result, err = tx[0].Exec(q, params...)
	} else {
		result, err = this.Db.Exec(q, params...)
	}

	if this.DebugMode {
		this.debug(q, "params", params, "execution time", time.Since(start))
	}

	return result, err
}

func (this *Mysql) parseError(err error) error {
	if err == nil {
		return nil
	}

	// Just a wrapper not to use sql lib directly from code
	if err == sql.ErrNoRows {
		return ErrNoRows
	}

	matches := this.matchStringGroups(errorsRegexp, err.Error())
	code, ok := matches["code"]
	if !ok {
		return err
	}

	switch code {
	case "1062":
		return ErrDuplicate
	default:
		return err
	}
}

// matchStringGroups matches regexp with capture groups. Returns map string string
func (this *Mysql) matchStringGroups(re *regexp.Regexp, s string) map[string]string {
	m := re.FindStringSubmatch(s)
	n := re.SubexpNames()

	r := make(map[string]string, len(m))

	if len(m) > 0 {
		m, n = m[1:], n[1:]
		for i, _ := range n {
			r[n[i]] = m[i]
		}
	}

	return r
}

func (this *Mysql) warn(msg string, args ...interface{}) {
	if this.Logger != nil {
		this.Logger.Warn(msg, args...)
	}
}

func (this *Mysql) debug(msg string, args ...interface{}) {
	if this.Logger != nil {
		this.Logger.Debug(msg, args...)
	}
}
