package middleware

import (
	"context"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = "Bearer "
)

type JWTMiddleware struct {
	PublicKey string
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{
		PublicKey: secret,
	}
}

func (m *JWTMiddleware) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return nil, err
		}

		value := md[authorizationHeader]

		if len(value) == 0 || !strings.HasPrefix(value[0], bearerPrefix) {
			return nil, status.Errorf(codes.Internal, "bearer prefix not found: %v", err)
		}

		tokenStr := strings.TrimPrefix(value[0], bearerPrefix)

		pubKey, errParse := jwt.ParseRSAPublicKeyFromPEM([]byte(m.PublicKey))

		if errParse != nil {
			return nil, status.Errorf(codes.Internal, "invalid public key: %v", err)
		}

		_, errParseTokenAlg := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
				return nil, status.Errorf(codes.Internal, "unexpected signing method: %v", t.Header["alg"])
			}

			return pubKey, nil
		})

		if errParseTokenAlg != nil {
			return nil, status.Errorf(codes.Internal, "unexpected signing method: %v", errParseTokenAlg)
		}

		return handler(ctx, req)
	}
}
