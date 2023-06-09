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

func (p *psql) UpdateUserCrosshairCount(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("id = ?", user.ID).Update("crosshairs_registered", user.CrosshairsRegistered+1)
	return user, tx.Error
}

func (p *psql) AddResetPasswordCodeAndTime(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).Update("password_reset_code", user.PasswordResetCode)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).Update("password_reset_code_time", user.PasswordResetCodeTime)
	return user, tx.Error
}

func (p *psql) GetUserByResetpasswordCode(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).Where("password_reset_code = ?", user.PasswordResetCode).First(&user)
	return user, tx.Error
}

func (p *psql) UpdateUserPassword(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).Where("password_reset_code = ?", user.PasswordResetCode).Update("password", user.Password)
	return user, tx.Error
}

func (p *psql) UpdateUserPasswordRaw(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).Update("password", user.Password)
	return user, tx.Error
}

func (p *psql) UpdateVerifyMailResendTime(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("e_mail = ?", user.EMail).Update("request_new_verify_mail_time", user.RequestNewVerifyMailTime)
	return user, tx.Error
}

func (p *psql) UpdateUserAvatarURL(user *database.UserAccount) (*database.UserAccount, error) {
	tx := p.db.Table(tableUsers).Where("id = ?", user.ID).Update("avatar_url", user.AvatarURL)
	return user, tx.Error
}
