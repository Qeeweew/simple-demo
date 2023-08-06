package dbcore

import (
	"context"
	"reflect"
	"simple-demo/common/log"
	"simple-demo/repository/dbinit"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var globalDB *gorm.DB

var redisDB *redis.Client

func Init() {
	globalDB = dbinit.MySQLInit()
	redisDB = dbinit.RedisInit()
}

func GetRedisClient() *redis.Client {
	return redisDB
}

// https://github.com/win5do/go-microservice-demo/blob/main/pkg/repository/db/dbcore/core.go#L107
// 如果使用跨模型事务则传参
func GetDB(ctx context.Context) *gorm.DB {
	iface := ctx.Value(ctxTransactionKey{})

	if iface != nil {
		tx, ok := iface.(*gorm.DB)
		if !ok {
			log.Logger.Sugar().Panicf("unexpect context value type: %s", reflect.TypeOf(tx))
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
