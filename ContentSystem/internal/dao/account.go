package dao

import (
	"ContentSystem/internal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type AccountDao struct {
	db *gorm.DB
}

func NewAccountDao(db *gorm.DB) *AccountDao {
	return &AccountDao{db: db}
}

// 校验存在
func (a *AccountDao) IsExist(userID string) (bool, error) {
	var account model.Account
	err := a.db.Where("user_id = ?", userID).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		fmt.Println("error isExist:", err)
		return false, err
	}
	return true, nil
}

// 创建账号
func (a *AccountDao) CreateAccount(account model.Account) error {
	if err := a.db.Create(&account).Error; err != nil {
		fmt.Println("AccountDao CreateAccount error:", err)
		return err
	}
	return nil
}

// 查找用户
func (a *AccountDao) FirstByUserId(userID string) (*model.Account, error) {
	var account model.Account
	err := a.db.Where("user_id = ?", userID).First(&account).Error
	if err != nil {
		fmt.Println("error firstByUserId:", err)
		return nil, err
	}
	return &account, nil
}
