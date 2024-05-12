package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Handler - функция, которая выполняется в транзакции
type Handler func(ctx context.Context) error

// Client клиент для работы с БД
type Client interface {
	DB() DB
	Close() error
}

// TxManager менеджер транзакций, который выполняет указанный пользователем обработчик в транзакции
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Query обертка над запросом, хранящая имя запроса и сам запрос
// Имя запроса используется для логирования и потенциально может использоваться еще где-то, например, для трейсинга
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor интерфейс для работы с транзакциями
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecer комбинирует NamedExecer и QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer интерфейс для работы с именованными запросами с помощью тегов в структурах
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest any, q Query, args ...any) error
	ScanAllContext(ctx context.Context, dest any, q Query, args ...any) error
}

// QueryExecer интерфейс для работы с обычными запросами
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, arguments ...any) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...any) pgx.Row
}

// Pinger интерфейс для проверки соединения с БД
type Pinger interface {
	Ping(ctx context.Context) error
}

// Tr интерфейс для работы с pgx.Conn, pgxpool.Conn or pgxpool.Pool
type Tr interface {
	Begin(ctx context.Context) (pgx.Tx, error)

	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// DB интерфейс для работы с БД
type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}
