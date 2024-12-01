package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	bearer        string = "bearer"
	authorization string = "authorization"
	jwtpayload    string = "jwtpayload"
)

type userInfoKey struct{}

type UserInfo struct {
	ID         uint64 `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	Name       string `json:"name,omitempty"`
	Mobile     string `json:"phone_number,omitempty"`
	Email      string `json:"email,omitempty"`
	MerchantID uint64 `json:"merchant_id,omitempty"`
	ClientID   string `json:"client_id,omitempty"`
}

func extractTokenFromAuthHeader(auth string) (string, bool) {
	if tokenType, token, ok := strings.Cut(auth, " "); ok {
		if strings.EqualFold(tokenType, bearer) {
			return token, true
		}
	}

	return "", false
}

func extractPayloadFromToken(token string) (string, bool) {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return "", false
	}

	return tokenParts[1], true
}

func withUserInfoClaims(ctx context.Context, payload string) context.Context {
	var userInfo UserInfo
	if claims, err := base64.RawURLEncoding.DecodeString(payload); err != nil {
		log.Printf("error while decoding jwt payload: %v", err)

		return ctx
	} else if err := json.Unmarshal(claims, &userInfo); err != nil {
		log.Printf("error while unmarshalling jwt payload: %v", err)

		return ctx
	}

	return context.WithValue(ctx, userInfoKey{}, userInfo)
}

// WithUserInfoContext returns context with user info value extracted from auth token metadata
func WithUserInfoContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	var payload string
	payloadHeader, ok := md[jwtpayload]
	if ok {
		payload = payloadHeader[0]
	} else {
		authHeader, ok := md[authorization]
		if !ok {
			return ctx
		}

		token, ok := extractTokenFromAuthHeader(authHeader[0])
		if ok {
			payload, ok = extractPayloadFromToken(token)
			if !ok {
				return ctx
			}
		}
	}

	return withUserInfoClaims(ctx, payload)
}

// WithUserInfoRequestContext returns http request with user info context value extracted from auth token header
func WithUserInfoRequestContext(req *http.Request) *http.Request {
	ctx := req.Context()

	var payload string
	payloadHeader := req.Header.Get(jwtpayload)
	if len(payloadHeader) > 0 {
		payload = payloadHeader
	} else {
		authHeader := req.Header.Get(authorization)
		if len(authHeader) > 0 {
			token, ok := extractTokenFromAuthHeader(authHeader)
			if ok {
				payload, ok = extractPayloadFromToken(token)
				if !ok {
					return req
				}
			} else {
				return req
			}
		} else {
			return req
		}
	}

	ctx = withUserInfoClaims(ctx, payload)

	return req.WithContext(ctx)
}

// UserInfoFromContext returns user info from context if it exists
func UserInfoFromContext(ctx context.Context) (UserInfo, bool) {
	if val := ctx.Value(userInfoKey{}); val != nil {
		userInfo, ok := val.(UserInfo)
		if ok {
			return userInfo, ok
		}
	}

	return UserInfo{}, false
}

// WithUserInfo return context with added user info
func WithUserInfo(ctx context.Context, userInfo UserInfo) context.Context {
	return context.WithValue(ctx, userInfoKey{}, userInfo)
}
