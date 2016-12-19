package sql

import (
	"database/sql"
	"sync"
	"sync/atomic"
	"unsafe"
)

type DB struct {
	*sql.DB

	stmtCachePtrMutex sync.Mutex     // used only by writers
	stmtCachePtr      unsafe.Pointer // *stmtCache
}

type stmtCache map[string]Stmt // map[query]Stmt

func NewDB(db *sql.DB) *DB {
	return &DB{
		DB: db,
	}
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return NewDB(db), nil
}

func (db *DB) Prepare(query string) (stmt Stmt, err error) {
	var m stmtCache

	if p := (*stmtCache)(atomic.LoadPointer(&db.stmtCachePtr)); p != nil {
		m = *p
		if stmt = m[query]; stmt.Stmt != nil {
			return
		}
	}

	db.stmtCachePtrMutex.Lock()
	defer db.stmtCachePtrMutex.Unlock()

	if p := (*stmtCache)(atomic.LoadPointer(&db.stmtCachePtr)); p != nil {
		m = *p
		if stmt = m[query]; stmt.Stmt != nil {
			return
		}
	}

	sqlStmt, err := db.DB.Prepare(query)
	if err != nil {
		return
	}
	stmt = Stmt{Stmt: sqlStmt}

	m2 := make(stmtCache, len(m)+1)
	for k, v := range m {
		m2[k] = v
	}
	m2[query] = stmt

	atomic.StorePointer(&db.stmtCachePtr, unsafe.Pointer(&m2))
	return
}

// ================================================================================================================

type Stmt struct {
	*sql.Stmt
}

func (Stmt) Close() error {
	return nil
}
