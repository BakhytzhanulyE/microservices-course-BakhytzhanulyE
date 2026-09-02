package interceptor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

// TestIsClientError фиксирует границу между «клиент не прав» и «сервис сломался»:
// от неё зависит уровень лога, а значит и то, на что будут срабатывать алерты.
func TestIsClientError(t *testing.T) {
	clientErrors := []codes.Code{
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.Unauthenticated,
		codes.FailedPrecondition,
		codes.OutOfRange,
		codes.Canceled,
	}
	for _, code := range clientErrors {
		require.Truef(t, isClientError(code), "%s — отказ по бизнес-правилу, а не авария", code)
	}

	serverErrors := []codes.Code{
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss,
		codes.Unknown,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Unimplemented,
	}
	for _, code := range serverErrors {
		require.Falsef(t, isClientError(code), "%s — настоящая проблема, её нельзя прятать в WARN", code)
	}
}
