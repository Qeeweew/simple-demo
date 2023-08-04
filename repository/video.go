package repository

import (
	"simple-demo/common/model"

	"gorm.io/gorm"
)

type videoRepository struct {
	*gorm.DB
}

func NewVideoRepository(db *gorm.DB) model.VideoRepository {
	return &videoRepository{
		db,
	}
}

func (v *videoRepository) Save(video *model.Video) error {
	return v.Create(video).Error
}

func (v *videoRepository) FindListByUserId(userId uint, videos *[]model.Video) error {
	return v.Transaction(func(tx *gorm.DB) (err error) {
		err = v.Preload("Author").Where(model.Video{AuthorId: userId}).Find(videos).Error
		return
	})
}

func (v *videoRepository) FeedList(limit uint, videos *[]model.Video) error {
	return v.Transaction(func(tx *gorm.DB) (err error) {
		err = v.Preload("Author").Limit(int(limit)).Order("created_at DESC").Find(videos).Error
		return
	})
}

func (v *videoRepository) FillExtraData(userId uint, video *model.Video) (err error) {
	return v.Transaction(func(tx *gorm.DB) (err error) {
		err = NewUserRepository(tx).FillExtraData(userId, &video.Author)
		if err != nil {
			return
		}
		var (
			FavoriteCount int64
			CommentCount  int64
			IsFavorite    bool
		)
		FavoriteCount, err = NewFavoriteRepository(tx).VideoFavoriteCount(video.Id)
		if err != nil {
			return
		}
		CommentCount, err = NewCommentRepository(tx).VideoCommentCount(video.Id)
		if err != nil {
			return
		}
		if userId != 0 {
			IsFavorite, err = NewFavoriteRepository(tx).IsFavorite(userId, video.Id)
		}
		video.Extra = &model.VideoExtra{
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    IsFavorite,
		}
		return
	})
}
