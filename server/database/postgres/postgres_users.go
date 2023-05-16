package postgres

import "github.com/devusSs/crosshairs/database"

func (p *psql) AddUser(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Create(user)
	return user, tx.Error
}

func (p *psql) GetUserByVerificationCode(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("verification_code = ?", user.VerificationCode).First(&user)
	return user, tx.Error
}

func (p *psql) UpdateUserVerification(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).Update("verified_mail", user.VerifiedMail)
	return user, tx.Error
}

func (p *psql) GetUserByEmail(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).First(&user)
	return user, tx.Error
}

func (p *psql) UpdateUserLogin(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).Updates(database.UserAccount{LoginIP: user.LoginIP, LastLogin: user.LastLogin})
	return user, tx.Error
}

func (p *psql) GetUserByUID(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("id = ?", user.ID).First(&user)
	return user, tx.Error
}
