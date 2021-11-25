package jwtfmt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type JWTPayloader interface {
	Body(v interface{}) error
	Subject() JwtSubject
	Token() string
}

type JwtPayload struct {
	tokenType string
	rawToken  string
}

type JwtSubject struct {
	Provider string
	Subject  string
}

func NewJWTPayload(jwt string) JWTPayloader {
	p := new(JwtPayload)
	p.rawToken = jwt

	s := strings.Split(jwt, " ")
	if len(s) > 1 {
		p.tokenType = s[0]
		p.rawToken = s[1]
	}

	return p
}

func (p *JwtPayload) Body(v interface{}) error {

	parts := strings.Split(p.rawToken, ".")

	if len(parts) < 3 {
		return errors.New("invalid jwt format")
	}

	b, err := decodeSegment(parts[1])
	if err != nil {
		return fmt.Errorf("jwt body could not be decoded: %s", err.Error())
	}

	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("jwt payload failed to be parsed: %s", err.Error())
	}

	return nil
}

func (p *JwtPayload) Subject() JwtSubject {
	s := struct {
		Subject string `json:"sub"`
	}{}

	if err := p.Body(&s); err != nil {
		return JwtSubject{}
	}

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
