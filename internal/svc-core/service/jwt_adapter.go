package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/konsultin/project-goes-here/dto"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/httpk"
	"github.com/konsultin/project-goes-here/internal/svc-core/pkg/valk"
	"github.com/konsultin/errk"
	"github.com/konsultin/logk"
	logkOption "github.com/konsultin/logk/option"
)

type JwtAdapter struct {
	Algorithm *jwt.SigningMethodHMAC
	Issuer    string
	Secret    string
}

type IssueJwtPayload struct {
	Subject     string
	Audience    []string
	Lifetime    int64
	SessionId   string
	SubjectType int32
	CreatedAt   sql.NullTime
	metadata    map[string]string
}

type JwtResponse struct {
	Aud  []string               `json:"aud"`
	Ent  int32                  `json:"ent"`
	Exp  int64                  `json:"exp"`
	Iat  int64                  `json:"iat"`
	Iss  string                 `json:"iss"`
	Jti  string                 `json:"jti"`
	Sub  string                 `json:"sub"`
	Meta map[string]interface{} `json:"meta"`
}

func (s *Service) NewJwtAdapter() *JwtAdapter {
	return &JwtAdapter{
		Algorithm: jwt.SigningMethodHS512,
		Issuer:    s.config.JwtIssuer,
		Secret:    s.config.JwtSecret,
	}
}

func (ja *JwtAdapter) Issue(options IssueJwtPayload) (*dto.Session, error) {
	// Retrieve created at from options
	var createdAt time.Time
	if options.CreatedAt.Valid {
		createdAt = options.CreatedAt.Time
	} else {
		createdAt = time.Now()
	}

	// Calculate expired Refresh time based on createdAt
	exp := createdAt.Add(time.Second * time.Duration(options.Lifetime)).Unix()

	// prepare signing token
	signingMethod := ja.Algorithm
	token := jwt.New(signingMethod)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exp
	claims["iss"] = ja.Issuer
	claims["aud"] = options.Audience
	claims["jti"] = options.SessionId
	claims["sub"] = options.Subject
	claims["iat"] = createdAt.Unix()
	claims["ent"] = options.SubjectType
	claims["meta"] = options.metadata

	// create string token
	tokenString, err := token.SignedString([]byte(ja.Secret))
	if err != nil {
		logk.Get().Error("failed to signedString for JWT Token", logkOption.Error(err))
		return nil, errk.Trace(err)
	}

	return &dto.Session{
		Token:     tokenString,
		ExpiredAt: exp,
	}, nil
}

func (ja *JwtAdapter) Validate(token string, options *dto.ValidateJwt_Payload) (*JwtResponse, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errk.Trace(fmt.Errorf("token method is invalid"))
		}
		return []byte(ja.Secret), nil
	}

	mapClaims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &mapClaims, keyFunc)
	if err != nil {
		return nil, err
	}

	// convert map claims to json
	jsonData, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, errk.Trace(err)
	}

	var claims JwtResponse
	err = json.Unmarshal(jsonData, &claims)
	if err != nil {
		return nil, errk.Trace(err)
	}

	isValid := false
	for _, val := range options.Audience {
		if valk.InArrayString(val, claims.Aud) {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, httpk.UnauthorizedError.Wrap(fmt.Errorf("audience did not match. Expected=%s Actual=%s",
			strings.Join(options.Audience, ","),
			strings.Join(claims.Aud, ","),
		))
	}

	return &claims, nil
}

// ValidateWithoutAudience validates JWT token without checking audience
func (ja *JwtAdapter) ValidateWithoutAudience(token string) (*JwtResponse, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errk.Trace(fmt.Errorf("token method is invalid"))
		}
		return []byte(ja.Secret), nil
	}

	mapClaims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &mapClaims, keyFunc)
	if err != nil {
		return nil, err
	}

	// convert map claims to json
	jsonData, err := json.Marshal(mapClaims)
	if err != nil {
		return nil, errk.Trace(err)
	}

	var claims JwtResponse
	err = json.Unmarshal(jsonData, &claims)
	if err != nil {
		return nil, errk.Trace(err)
	}

	return &claims, nil
}
