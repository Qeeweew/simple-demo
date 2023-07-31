package repository

import (
	"context"
	"simple-demo/common/model"
	"simple-demo/repository/dbcore"
)

type tableVisitor struct{}

func NewTableVistor() *tableVisitor {
	return &tableVisitor{}
}

func (*tableVisitor) User(ctx context.Context) model.UserRepository {
	return NewUserRepository(dbcore.GetDB(ctx))
}

func (*tableVisitor) Video(ctx context.Context) model.VideoRepository {
	return NewVideoRepository(dbcore.GetDB(ctx))
}

func (*tableVisitor) Relation(ctx context.Context) model.RelationRepository {
	return NewRelationRepository(dbcore.GetDB(ctx))
}
