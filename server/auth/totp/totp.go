package otp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"time"

	"github.com/pnnh/multiverse-server/config"
	"github.com/sirupsen/logrus"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type RegisterOut struct {
	Image  string
	Secret string
}

func Register(account string) (*RegisterOut, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.ISSUER,
		AccountName: account,
	})
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	out := &RegisterOut{
		Image:  base64.StdEncoding.EncodeToString(buf.Bytes()),
		Secret: key.Secret(),
	}
	return out, nil
}

func Validate(secret, passcode string) error {
	// 通过account从数据库中查询secret
	logrus.Println("Validating TOTP...")
	if len(secret) < 1 {
		return fmt.Errorf("key不存在")
	}
	totpOpts := totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	}
	//valid := totp.Validate(passcode, secret)
	valid, err := totp.ValidateCustom(passcode, secret, time.Now(), totpOpts)
	if err != nil {

		return fmt.Errorf("ValidateCustom出错")
	}
	if valid {
		return nil
	} else {
		return fmt.Errorf("invalid passcode")
	}
}

func GenerateCode(secret string, t time.Time) (string, error) {
	return totp.GenerateCode(secret, t)
}
