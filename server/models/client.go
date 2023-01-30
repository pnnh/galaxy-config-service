package models

import (
	"database/sql"
	"fmt"
	"github.com/pnnh/quantum-go/services/sqlxsvc"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ory/fosite"
)

type ClientTable struct {
	Pk             string         `json:"pk"`
	ID             string         `json:"id"`
	Secret         string         `json:"secret"`
	RotatedSecrets sql.NullString `json:"rotated_secrets" db:"rotated_secrets"`
	RedirectURIs   sql.NullString `json:"redirect_uris" db:"redirect_uris"`
	GrantTypes     sql.NullString `json:"grant_types" db:"grant_types"`
	ResponseTypes  sql.NullString `json:"response_types" db:"response_types"`
	Scopes         sql.NullString `json:"scopes"`
	Audience       sql.NullString `json:"audience"`
	Public         sql.NullInt32  `json:"public"`
}

func (t *ClientTable) ToModel() *ClientModel {
	model := &ClientModel{
		ID:     t.ID,
		Secret: []byte(t.Secret),
		//RotatedSecrets: make([][]byte, 0),
		//RedirectURIs:   []string{t.RedirectURIs.String},
		//GrantTypes:     []string{t.GrantTypes.String},
		//ResponseTypes:  []string{t.ResponseTypes.String},
		//Scopes:         []string{t.Scopes.String},
		//Audience:       []string{t.Audience.String},
	}

	for _, v := range strings.Split(t.RotatedSecrets.String, ",") {
		model.RotatedSecrets = append(model.RotatedSecrets, []byte(v))
	}

	// for _, v := range strings.Split(t.RedirectURIs.String, ",") {
	// 	model.RedirectURIs = append(model.RedirectURIs,  v)
	// }
	model.RedirectURIs = append(model.RedirectURIs, strings.Split(t.RedirectURIs.String, ",")...)
	// for _, v := range strings.Split(t.GrantTypes.String, ",") {
	// 	model.GrantTypes = append(model.GrantTypes,  v)
	// }
	model.GrantTypes = append(model.GrantTypes, strings.Split(t.GrantTypes.String, ",")...)
	// for _, v := range strings.Split(t.ResponseTypes.String, ",") {
	// 	model.ResponseTypes = append(model.ResponseTypes,  v)
	// }
	model.ResponseTypes = append(model.ResponseTypes, strings.Split(t.ResponseTypes.String, ",")...)
	// for _, v := range strings.Split(t.Scopes.String, ",") {
	// 	model.Scopes = append(model.Scopes,  v)
	// }
	model.Scopes = append(model.Scopes, strings.Split(t.Scopes.String, ",")...)
	// for _, v := range strings.Split(t.Audience.String, ",") {
	// 	model.Audience = append(model.Audience,  v)
	// }
	model.Audience = append(model.Audience, strings.Split(t.Audience.String, ",")...)

	if t.Public.Valid && t.Public.Int32 == 1 {
		model.Public = true
	}

	return model
}

type ClientModel struct {
	ID             string   `json:"id"`
	Secret         []byte   `json:"client_secret,omitempty"`
	RotatedSecrets [][]byte `json:"rotated_secrets,omitempty"`
	RedirectURIs   []string `json:"redirect_uris"`
	GrantTypes     []string `json:"grant_types"`
	ResponseTypes  []string `json:"response_types"`
	Scopes         []string `json:"scopes"`
	Audience       []string `json:"audience"`
	Public         bool     `json:"public"`
}

func (c *ClientModel) GetID() string {
	return c.ID
}

func (c *ClientModel) IsPublic() bool {
	return c.Public
}

func (c *ClientModel) GetAudience() fosite.Arguments {
	return c.Audience
}

func (c *ClientModel) GetRedirectURIs() []string {
	return c.RedirectURIs
}

func (c *ClientModel) GetHashedSecret() []byte {
	return c.Secret
}

func (c *ClientModel) GetRotatedHashes() [][]byte {
	return c.RotatedSecrets
}

func (c *ClientModel) GetScopes() fosite.Arguments {
	return c.Scopes
}

func (c *ClientModel) GetGrantTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// JSON array containing a list of the OAuth 2.0 Grant Types that the Client is declaring
	// that it will restrict itself to using.
	// If omitted, the default is that the Client will use only the authorization_code Grant Type.
	if len(c.GrantTypes) == 0 {
		return fosite.Arguments{"authorization_code"}
	}
	return fosite.Arguments(c.GrantTypes)
}

func (c *ClientModel) GetResponseTypes() fosite.Arguments {
	// https://openid.net/specs/openid-connect-registration-1_0.html#ClientMetadata
	//
	// JSON array containing a list of the OAuth 2.0 response_type values that the Client is declaring
	// that it will restrict itself to using. If omitted, the default is that the Client will use
	// only the code Response Type.
	if len(c.ResponseTypes) == 0 {
		return fosite.Arguments{"code"}
	}
	return fosite.Arguments(c.ResponseTypes)
}

func GetClient(id string) (*ClientModel, error) {
	sqlText := `select pk, id, secret, rotated_secrets, redirect_uris, response_types, grant_types, scopes,
		audience, public
	from applications where id = :id;`

	sqlParams := map[string]interface{}{"id": id}
	var sqlResults []*ClientTable

	rows, err := sqlxsvc.NamedQuery(sqlText, sqlParams)
	if err != nil {
		return nil, fmt.Errorf("NamedQuery: %w", err)
	}
	if err = sqlx.StructScan(rows, &sqlResults); err != nil {
		return nil, fmt.Errorf("StructScan: %w", err)
	}

	for _, v := range sqlResults {
		model := v.ToModel()
		return model, nil
	}

	return nil, nil
}
