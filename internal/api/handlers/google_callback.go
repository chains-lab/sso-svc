package handlers

import (
	"context"
	"encoding/json"

	"github.com/chains-lab/chains-auth/internal/api/responses"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/proto-storage/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (a Service) GoogleCallback(
	ctx context.Context,
	req *auth.GoogleCallbackRequest,
) (*auth.TokensPairResponse, error) {
	requestID := uuid.New()
	log := Log(ctx, requestID)

	// 1) Провека наличия кода
	code := req.Code
	if code == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing code")
	}

	// 2) Обмен кода на токен
	token, err := a.cfg.GoogleOAuth().Exchange(ctx, code)
	if err != nil {
		log.WithError(err).Errorf("error exchanging code for token: %s", code)
		return nil, status.Errorf(codes.Internal, "failed to exchange code")
	}

	// 3) Получаем UserInfo из Google
	client := a.cfg.GoogleOAuth().Client(ctx, token)
	httpResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.WithError(err).Error("error fetching userinfo from Google")
		return nil, status.Errorf(codes.Internal, "failed to fetch user info")
	}
	defer httpResp.Body.Close()

	var ui struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(httpResp.Body).Decode(&ui); err != nil {
		log.WithError(err).Error("error decoding Google userinfo")
		return nil, status.Errorf(codes.Internal, "invalid user info")
	}

	md, _ := metadata.FromIncomingContext(ctx)
	ua := ""
	if vals := md.Get("user-agent"); len(vals) > 0 {
		ua = vals[0]
	}

	// 5) Ваша бизнес-логика: логиним/регистрируем пользователя
	_, tokensPair, err := a.app.Login(ctx, ui.Email, roles.User, ua)
	if err != nil {
		// конвертим в gRPC-ошибку через ваш презентер
		return nil, responses.AppError(ctx, requestID, err)
	}

	// 6) Возвращаем пару токенов
	return responses.TokensPair(tokensPair), nil
}
