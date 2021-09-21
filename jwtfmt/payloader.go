package jwtfmt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type JWTPayloader interface {
	Body(v interface{}) error
	Subject() JwtSubject
	Token() string
}

type JwtPayload struct {
	rawToken string
}

type JwtSubject struct {
	Provider string
	Subject  string
}

func NewJWTPayload(jwt string) JWTPayloader {
	return &JwtPayload{
		rawToken: jwt,
	}
}

func (p *JwtPayload) Body(v interface{}) error {

	s := strings.Split(p.rawToken, " ")[1]
	parts := strings.Split(s, ".")

	if len(parts) < 3 {
		return errors.New("NOT JWT TYPE")
	}

	b, err := decodeSegment(parts[1])
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

func (p *JwtPayload) Subject() JwtSubject {
	s := struct {
		Subject string `json:"sub"`
	}{}

	_ = p.Body(&s)
	sections := strings.Split(s.Subject, "|")

	if len(sections) > 1 {
		return JwtSubject{
			Provider: sections[0],
			Subject:  sections[1],
		}
	}

	return JwtSubject{
		Subject: sections[0],
	}
}

func (p *JwtPayload) Token() string {
	return p.rawToken
}

func (p *JwtPayload) GetType() string {
	s := strings.Split(p.rawToken, " ")
	return strings.ToLower(s[0])
}

func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
