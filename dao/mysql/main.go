package mysql

import (
	"bonex-middleware/config"
	"bonex-middleware/dao"
	"bonex-middleware/dao/mysql/driver"
	"fmt"
	"github.com/jmoiron/sqlx"
)

const migrationsDir = "./dao/mysql/migrations"

type mysqlDAO struct {
	mysql *driver.Mysql
}

type mysqlDAOTx struct {
	tx *sqlx.Tx
}

func NewMysql(c *config.Config, l driver.Logger) (dao.DbDAO, error) {
	m, err := driver.CreateConnection(&c.Mysql, l)
	if err != nil {
		return nil, err
	}

	err = driver.Migrate(&c.Mysql, migrationsDir)
	if err != nil {
		return nil, err
	}

	return &mysqlDAO{mysql: m}, nil
}

func (this *mysqlDAO) BeginTx() (dao.DAOTx, error) {
	var err error

	dtx := &mysqlDAOTx{}
	dtx.tx, err = this.mysql.Db.Beginx()

	return dtx, err
}

func (this *mysqlDAOTx) CommitTx() error {
	if this.tx == nil {
		return fmt.Errorf("tx not initialized")
	}

	return this.tx.Commit()
}

func (this *mysqlDAOTx) RollbackTx() error {
	if this.tx == nil {
		return nil
	}

	return this.tx.Rollback()
}

// getTx checks if tx param is not empty and returns an *sqlx.Tx
func (this *mysqlDAO) daoTx2Sqlx(tx dao.DAOTx) *sqlx.Tx {
	if tx == nil {
		return nil
	}

	t, ok := tx.(*mysqlDAOTx)
	if !ok {
		return nil
	}

	return t.tx
}
