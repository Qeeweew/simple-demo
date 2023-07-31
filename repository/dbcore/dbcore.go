package dbcore

import (
	"context"
	"reflect"
	"simple-demo/repository/dbinit"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var MySQL *gorm.DB
var globalDB *gorm.DB

// var Redis *redis.Client

func Init() {
	MySQL = dbinit.MySQLInit()

	// TODO: 等之后上缓存的时候再放开
	// Re
}

// https://github.com/win5do/go-microservice-demo/blob/main/pkg/repository/db/dbcore/core.go#L107
// 如果使用跨模型事务则传参
func GetDB(ctx context.Context) *gorm.DB {
	iface := ctx.Value(ctxTransactionKey{})

	if iface != nil {
		tx, ok := iface.(*gorm.DB)
		if !ok {
			logrus.Panicf("unexpect context value type: %s", reflect.TypeOf(tx))
			return nil
		}
		return tx
	}

	return globalDB.WithContext(ctx)
}

// empty struct used as a fixed key
type ctxTransactionKey struct{}

func CtxWithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxTransactionKey{}, tx)
}

type txImpl struct{}

func NewTxImpl() *txImpl {
	return &txImpl{}
}

func (*txImpl) Transaction(ctx context.Context, fn func(txctx context.Context) error) error {
	db := globalDB.WithContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		txctx := CtxWithTransaction(ctx, tx)
		return fn(txctx)
	})
}
