package application

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/mochisuna/load-test-sample/domain"
	"github.com/mochisuna/load-test-sample/domain/repository"

	"github.com/rs/xid"
)

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService inject userRepo
func NewUserService(userRepo repository.UserRepository) *userService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Refer(userID domain.UserID) (*domain.User, error) {
	log.Println("application.Refer")
	return s.userRepo.Get(userID)
}

func (s *userService) Register(ctx context.Context, user *domain.User) error {
	log.Println("application.Register")

	user.ID = userID()
	user.Name = name()
	user.SecretKey = secretKey(string(user.ID))
	now := int(time.Now().Unix())
	user.CreatedAt = now
	user.UpdatedAt = now

	err := s.userRepo.WithTransaction(ctx, func(tx *sql.Tx) error {
		if err := s.userRepo.Create(user); err != nil {
			return err
		}
		return nil
	})
	return err
}

// 雑にランダムな値を生成するシリーズ
func secretKey(id string) string {
	return id + ":" + xid.New().String() + "-" + xid.New().String()
}

func name() string {
	// 内容に意味はないので適当に苗字と名前を合体させる（メジャーらしい苗字・名前）
	family := []string{
		"佐藤", "鈴木", "高橋", "田中", "伊藤", "渡辺", "山本", "中村", "小林", "加藤",
		"吉田", "山田", "佐々木", "山口", "松本", "井上", "木村", "林", "斎藤", "清水",
	}
	name := []string{
		"蓮", "陽翔", "陽太", "樹", "悠人", "湊", "大翔", "蒼", "朝陽", "陽斗",
		"陽葵", "芽依", "莉子", "葵", "澪", "結菜", "凛", "結愛", "琴音", "陽菜",
	}
	rand.Seed(time.Now().UnixNano())
	return family[rand.Intn(len(family))] + name[rand.Intn(len(name))]
}

func userID() domain.UserID {
	id := xid.New().String()
	return domain.UserID(id)
}
