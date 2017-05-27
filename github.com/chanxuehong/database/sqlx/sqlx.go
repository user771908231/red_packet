package sqlx

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/chanxuehong/database/sql"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB

	sqlStmtCachePtrMutex sync.Mutex     // used only by writers
	sqlStmtCachePtr      unsafe.Pointer // *sqlStmtCache

	stmtCachePtrMutex sync.Mutex     // used only by writers
	stmtCachePtr      unsafe.Pointer // *stmtCache

	namedStmtCachePtrMutex sync.Mutex     // used only by writers
	namedStmtCachePtr      unsafe.Pointer // *namedStmtCache
}

type (
	sqlStmtCache   map[string]sql.Stmt  // map[query]sql.Stmt
	stmtCache      map[string]Stmt      // map[query]Stmt
	namedStmtCache map[string]NamedStmt // map[query]NamedStmt
)

func NewDB(db *sqlx.DB) *DB {
	return &DB{
		DB: db,
	}
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return NewDB(db), nil
}

func (db *DB) Prepare(query string) (stmt sql.Stmt, err error) {
	var m sqlStmtCache

	if p := (*sqlStmtCache)(atomic.LoadPointer(&db.sqlStmtCachePtr)); p != nil {
		m = *p
		if stmt = m[query]; stmt.Stmt != nil {
			return
		}
	}

	db.sqlStmtCachePtrMutex.Lock()
	defer db.sqlStmtCachePtrMutex.Unlock()

	if p := (*sqlStmtCache)(atomic.LoadPointer(&db.sqlStmtCachePtr)); p != nil {
		m = *p
		if stmt = m[query]; stmt.Stmt != nil {
			return
		}
	}

	stmtx, err := db.DB.Prepare(query)
	if err != nil {
		return
	}
	stmt = sql.Stmt{Stmt: stmtx}

	m2 := make(sqlStmtCache, len(m)+1)
	for k, v := range m {
		m2[k] = v
	}
	m2[query] = stmt

	atomic.StorePointer(&db.sqlStmtCachePtr, unsafe.Pointer(&m2))
	return
}

func (db *DB) Preparex(query string) (stmt Stmt, err error) {
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

	stmtx, err := db.DB.Preparex(query)
	if err != nil {
		return
	}
	stmt = Stmt{Stmt: stmtx}

	m2 := make(stmtCache, len(m)+1)
	for k, v := range m {
		m2[k] = v
	}
	m2[query] = stmt

	atomic.StorePointer(&db.stmtCachePtr, unsafe.Pointer(&m2))
	return
}

func (db *DB) PrepareNamed(query string) (stmt NamedStmt, err error) {
	var m namedStmtCache

	if p := (*namedStmtCache)(atomic.LoadPointer(&db.namedStmtCachePtr)); p != nil {
		m = *p
		if stmt = m[query]; stmt.Stmt != nil {
			return
		}
	}

	db.namedStmtCachePtrMutex.Lock()
	defer db.namedStmtCachePtrMutex.Unlock()

	if p := (*namedStmtCache)(atomic.LoadPointer(&db.namedStmtCachePtr)); p != nil {
		m = *p
		if stmt = m[query]; stmt.Stmt != nil {
			return
		}
	}

	stmtx, err := db.DB.PrepareNamed(query)
	if err != nil {
		return
	}
	stmt = NamedStmt{NamedStmt: stmtx}

	m2 := make(namedStmtCache, len(m)+1)
	for k, v := range m {
		m2[k] = v
	}
	m2[query] = stmt

	atomic.StorePointer(&db.namedStmtCachePtr, unsafe.Pointer(&m2))
	return
}

// ================================================================================================================

type Stmt struct {
	*sqlx.Stmt
}

func (Stmt) Close() error {
	return nil
}

type NamedStmt struct {
	*sqlx.NamedStmt
}

func (NamedStmt) Close() error {
	return nil
}
