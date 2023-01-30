package models

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pnnh/multiverse-server/server/helpers"
	"github.com/pnnh/quantum-go/services/sqlxsvc"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AccountTable struct {
	Pk          string         `json:"pk"`    // 主键标识
	Account     string         `json:"uname"` // 账号
	Password    string         `json:"upass"` // 密码
	CreateAt    time.Time      `json:"createat"`
	UpdateAt    time.Time      `json:"updateat"`
	Nickname    string         `json:"nickname"`
	Mail        sql.NullString `json:"mail"`
	Credentials sql.NullString `json:"credentials"`
}

type WebauthnCredentials struct {
	CredentialsSlice []webauthn.Credential
}

type AccountModel struct {
	AccountTable
	WebauthnCredentials
}

// NewAccountModel creates and returns a new User
func NewAccountModel(name string, displayName string) *AccountModel {
	user := &AccountModel{
		AccountTable: AccountTable{
			Pk:       helpers.NewPostId(),
			Account:  name,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
			Nickname: displayName,
		},
	}

	return user
}

// WebAuthnID returns the user's ID
func (u *AccountModel) WebAuthnID() []byte {
	// buf := make([]byte, binary.MaxVarintLen64)
	// binary.PutUvarint(buf, uint64(u.Id))
	// return buf
	return []byte(u.Pk)
}

// WebAuthnName returns the user's username
func (u *AccountModel) WebAuthnName() string {
	//return u.Name
	return u.Account
}

// WebAuthnDisplayName returns the user's display name
func (u *AccountModel) WebAuthnDisplayName() string {
	return u.Nickname
}

// WebAuthnIcon is not (yet) implemented
func (u *AccountModel) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *AccountModel) AddCredential(cred webauthn.Credential) {
	u.CredentialsSlice = append(u.CredentialsSlice, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u *AccountModel) WebAuthnCredentials() []webauthn.Credential {
	if len(u.CredentialsSlice) < 1 && u.Credentials.Valid {
		decodeBytes, err := base64.StdEncoding.DecodeString(u.Credentials.String)
		if err != nil {
			logrus.Errorln("WebAuthnCredentials DecodeString: %w", err)
			return u.CredentialsSlice
		}
		webauthnCredentials := &WebauthnCredentials{}
		if err := json.Unmarshal(decodeBytes, webauthnCredentials); err != nil {
			logrus.Errorln("WebAuthnCredentials Unmarshal error: %w", err)
			return u.CredentialsSlice
		}
		u.WebauthnCredentials = *webauthnCredentials
	}
	return u.CredentialsSlice
}

func (model *AccountModel) MarshalCredentials() (string, error) {
	if len(model.CredentialsSlice) < 1 {
		return "", fmt.Errorf("credentialsSlice为空")
	}
	data, err := json.Marshal(model.WebauthnCredentials)
	if err != nil {
		return "", fmt.Errorf("MarshalCredentials error: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u *AccountModel) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.CredentialsSlice {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}

func GetAccount(account string) (*AccountModel, error) {
	sqlText := `select pk, account, mail, nickname, credentials from accounts where account = :account and status = 1;`

	sqlParams := map[string]interface{}{"account": account}
	var sqlResults []*AccountTable

	rows, err := sqlxsvc.NamedQuery(sqlText, sqlParams)
	if err != nil {
		return nil, fmt.Errorf("NamedQuery: %w", err)
	}
	if err = sqlx.StructScan(rows, &sqlResults); err != nil {
		return nil, fmt.Errorf("StructScan: %w", err)
	}

	for _, v := range sqlResults {
		return &AccountModel{AccountTable: *v}, nil
	}

	return nil, nil
}

func PutAccount(model *AccountModel) error {
	sqlText := `insert into accounts(pk, createat, updateat, account, password, nickname, status)
	values(:pk, :createat, :updateat, :account, :password, :nickname, 1)`

	sqlParams := map[string]interface{}{"pk": model.Pk, "createat": model.CreateAt, "updateat": model.UpdateAt,
		"account": model.Account, "password": "", "nickname": model.Nickname}

	_, err := sqlxsvc.NamedExec(sqlText, sqlParams)
	if err != nil {
		return fmt.Errorf("PutAccount: %w", err)
	}
	return nil
}

func UpdateAccountCredentials(model *AccountModel) error {
	sqlText := `update accounts set credentials = :credentials where pk = :pk;`

	credentials, err := model.MarshalCredentials()
	if err != nil {
		return fmt.Errorf("MarshalCredentiials: %w", err)
	}

	sqlParams := map[string]interface{}{"pk": model.Pk, "credentials": credentials}

	_, err = sqlxsvc.NamedExec(sqlText, sqlParams)
	if err != nil {
		return fmt.Errorf("PutAccount: %w", err)
	}
	return nil
}
