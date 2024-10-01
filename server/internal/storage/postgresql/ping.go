package postgresql

func (p *PostgreSQLRepo) Ping() error {
	return p.db.Ping()
}
