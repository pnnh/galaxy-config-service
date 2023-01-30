package authorizationserver

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/ory/fosite"
	"github.com/pnnh/multiverse-server/config"
	"github.com/sirupsen/logrus"

	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/token/jwt"
)

// fosite requires four parameters for the server to get up and running:
// 1. config - for any enforcement you may desire, you can do this using `compose.Config`. You like PKCE, enforce it!
// 2. store - no auth service is generally useful unless it can remember clients and users.
//    fosite is incredibly composable, and the store parameter enables you to build and BYODb (Bring Your Own Database)
// 3. secret - required for code, access and refresh token generation.
// 4. privateKey - required for id/jwt token generation.

var (
	// Check the api documentation of `compose.Config` for further configuration options.
	fositeConfig = &fosite.Config{
		AccessTokenLifespan: time.Minute * 30,
		GlobalSecret:        secret,
		// ...
	}

	// This is the example storage that contains:
	// * an OAuth2 Client with id "my-client" and secrets "foobar" and "foobaz" capable of all oauth2 and open id connect grant and response types.
	// * a User for the resource owner password credentials grant type with username "peter" and password "secret".
	//
	// You will most likely replace this with your own logic once you set up a real world application.
	//store = storage.NewExampleStore()
	// todo 测试目的
	// store = &storage.MemoryStore{
	// 	IDSessions: make(map[string]fosite.Requester),
	// // 	Clients: map[string]fosite.Client{
	// 		"pwa": &fosite.DefaultClient{
	// 			ID:             "pwa",
	// 			Secret:         []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`),            // = "foobar"
	// 			RotatedSecrets: [][]byte{[]byte(`$2y$10$X51gLxUQJ.hGw1epgHTE5u0bt64xM0COU7K9iAp.OFg8p2pUd.1zC `)}, // = "foobaz",
	// 			RedirectURIs:   []string{"https://127.0.0.1:3500/login/callback"},
	// 			ResponseTypes:  []string{"id_token", "code", "token", "id_token token", "code id_token", "code token", "code id_token token"},
	// 			GrantTypes:     []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
	// 			Scopes:         []string{"fosite", "openid", "photos", "offline"},
	// 		},

	// // 	},
	// 	Users: map[string]storage.MemoryUserRelation{
	// 		"peter": {
	// 			// This store simply checks for equality, a real storage implementation would obviously use
	// 			// a hashing algorithm for encrypting the user password.
	// 			Username: "peter",
	// 			Password: "secret",
	// 		},
	// 	},
	// 	AuthorizeCodes:         map[string]storage.StoreAuthorizeCode{},
	// 	AccessTokens:           map[string]fosite.Requester{},
	// 	RefreshTokens:          map[string]storage.StoreRefreshToken{},
	// 	PKCES:                  map[string]fosite.Requester{},
	// 	AccessTokenRequestIDs:  map[string]string{},
	// 	RefreshTokenRequestIDs: map[string]string{},
	// 	IssuerPublicKeys:       map[string]storage.IssuerPublicKeys{},
	// }
	store = NewDatabaseStore()

	// This secret is used to sign authorize codes, access and refresh tokens.
	// It has to be 32-bytes long for HMAC signing. This requirement can be configured via `compose.Config` above.
	// In order to generate secure keys, the best thing to do is use crypto/rand:
	//
	// ```
	// package main
	//
	// import (
	//	"crypto/rand"
	//	"encoding/hex"
	//	"fmt"
	// )
	//
	// func main() {
	//	var secret = make([]byte, 32)
	//	_, err := rand.Read(secret)
	//	if err != nil {
	//		panic(err)
	//	}
	// }
	// ```
	//
	// If you require this to key to be stable, for example, when running multiple fosite servers, you can generate the
	// 32byte random key as above and push it out to a base64 encoded string.
	// This can then be injected and decoded as the `var secret []byte` on server start.
	secret = []byte("some-cool-secret-that-is-32bytes")

	// privateKey is used to sign JWT tokens. The default strategy uses RS256 (RSA Signature with SHA-256)
	//privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	privateKey *rsa.PrivateKey
)

func init() {
	// privateBytes, err := ioutil.ReadFile("./files/cert/rs256-private.pem")
	// if err != nil {
	// 	logrus.Fatalln("privatekey init", err)
	// 	return
	// }
	// // key, err := ssh.ParseRawPrivateKey(bytes)
	// // if err != nil {
	// // 	logrus.Fatalln("ParseRawPrivateKey ", err)
	// // 	return
	// // }
	// // privateKey = key.(*rsa.PrivateKey)

	// block, _ := pem.Decode(privateBytes) //将密钥解析成私钥实例
	// if block == nil {
	// 	panic("private key error!")
	// }
	// priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	// if err != nil {
	// 	panic(err)
	// }
	// privateKey = priv

	privateString := config.LoadAwsConfig("rs256-private.pem", "default")
	block, _ := pem.Decode([]byte(privateString)) //将密钥解析成私钥实例
	if block == nil {
		logrus.Fatalln("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		logrus.Fatalln("privateKeyBytes", err)
	}
	privateKey = priv

	oauth2 = compose.ComposeAllEnabled(fositeConfig, store, privateKey)

	/// 以下手动生成====
	// privateKey, _ = rsa.GenerateKey(rand.Reader, 2048)

	// pkcs1PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	// block:=&pem.Block{
	// 	Type:"RSA PRIVATE KEY",
	// 	Bytes:pkcs1PrivateKey,
	// }
	// privateBuffer := bytes.NewBufferString("")
	// ok := pem.Encode(privateBuffer,block)
	// fmt.Printf("privateBuffer %v\n%s", ok, privateBuffer.String())

	// PublicKey := &privateKey.PublicKey

	// pkixPublicKey, err := x509.MarshalPKIXPublicKey(PublicKey)
	// if err != nil {
	// 	panic(fmt.Sprintf("MarshalPKIXPublicKey(PublicKey) %v", err))
	// }
	// block1:=&pem.Block{
	// 	Type:"RSA PUBLIC KEY",
	// 	Bytes:pkixPublicKey,
	// }
	// publicBuffer := bytes.NewBufferString("")
	// ok = pem.Encode(publicBuffer, block1)
	// fmt.Printf("publicBuffer %v\n%s", ok, publicBuffer.String())
}

// Build a fosite instance with all OAuth2 and OpenID Connect handlers enabled, plugging in our configurations as specified above.
//var oauth2 = compose.ComposeAllEnabled(config, store, privateKey)
var oauth2 fosite.OAuth2Provider

// A session is passed from the `/auth` to the `/token` endpoint. You probably want to store data like: "Who made the request",
// "What organization does that person belong to" and so on.
// For our use case, the session will meet the requirements imposed by JWT access tokens, HMAC access tokens and OpenID Connect
// ID Tokens plus a custom field

// newSession is a helper function for creating a new session. This may look like a lot of code but since we are
// setting up multiple strategies it is a bit longer.
// Usually, you could do:
//
//	session = new(fosite.DefaultSession)
func newSession(user string) *openid.DefaultSession {
	return &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{
			Issuer:      "https://fosite.my-application.com",
			Subject:     user,
			Audience:    []string{"https://my-client.my-application.com"},
			ExpiresAt:   time.Now().Add(time.Hour * 6),
			IssuedAt:    time.Now(),
			RequestedAt: time.Now(),
			AuthTime:    time.Now(),
		},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}
}
