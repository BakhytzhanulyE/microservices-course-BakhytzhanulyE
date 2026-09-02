package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/model"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/service"
)

// fakeService подменяет бизнес-логику: HTTP-слой проверяем отдельно от неё.
type fakeService struct {
	order     model.Order
	createErr error
	getErr    error
	payErr    error
	cancelErr error

	transactionUUID string
}

var _ service.OrderService = (*fakeService)(nil)

func (f *fakeService) Create(_ context.Context, _ model.CreateOrderParams) (model.Order, error) {
	return f.order, f.createErr
}

func (f *fakeService) Get(_ context.Context, _ string) (model.Order, error) {
	return f.order, f.getErr
}

func (f *fakeService) Pay(_ context.Context, _ model.PayOrderParams) (string, error) {
	return f.transactionUUID, f.payErr
}

func (f *fakeService) Cancel(_ context.Context, _ string) error {
	return f.cancelErr
}

const validUserUUID = "6f3b2d1e-1111-4a2b-9c3d-4e5f60718293"

func do(t *testing.T, svc *fakeService, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	router := NewAPI(svc).Router()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	return rec
}

func TestHealth(t *testing.T) {
	t.Parallel()

	rec := do(t, &fakeService{}, http.MethodGet, "/health", "")

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
}

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		svc        *fakeService
		wantStatus int
	}{
		{
			name:       "валидный запрос создаёт заказ",
			body:       `{"user_uuid":"` + validUserUUID + `","part_uuids":["p1"]}`,
			svc:        &fakeService{order: model.Order{UUID: "o1", TotalPrice: 42}},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "битый JSON — 400",
			body:       `{`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "лишнее поле в теле — 400",
			body:       `{"user_uuid":"` + validUserUUID + `","part_uuids":["p1"],"lishnee":1}`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "user_uuid не UUID — 400",
			body:       `{"user_uuid":"не-uuid","part_uuids":["p1"]}`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "неизвестная деталь — 400",
			body:       `{"user_uuid":"` + validUserUUID + `","part_uuids":["p1"]}`,
			svc:        &fakeService{createErr: model.ErrPartsNotFound},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "каталог недоступен — 503",
			body:       `{"user_uuid":"` + validUserUUID + `","part_uuids":["p1"]}`,
			svc:        &fakeService{createErr: model.ErrInventoryUnavailable},
			wantStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := do(t, tt.svc, http.MethodPost, "/api/v1/orders", tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestGetOrder(t *testing.T) {
	t.Parallel()

	t.Run("существующий заказ отдаётся с полями", func(t *testing.T) {
		t.Parallel()

		svc := &fakeService{order: model.Order{
			UUID:       "o1",
			UserUUID:   validUserUUID,
			PartUUIDs:  []string{"p1", "p2"},
			TotalPrice: 150.5,
			Status:     model.OrderStatusPendingPayment,
		}}

		rec := do(t, svc, http.MethodGet, "/api/v1/orders/o1", "")
		require.Equal(t, http.StatusOK, rec.Code)

		var resp OrderResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, "o1", resp.OrderUUID)
		assert.Equal(t, "PENDING_PAYMENT", resp.Status)
		assert.Len(t, resp.PartUUIDs, 2)
	})

	t.Run("несуществующий заказ — 404", func(t *testing.T) {
		t.Parallel()

		rec := do(t, &fakeService{getErr: model.ErrOrderNotFound}, http.MethodGet, "/api/v1/orders/нет", "")

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPayOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		svc        *fakeService
		wantStatus int
	}{
		{
			name:       "оплата проходит",
			body:       `{"payment_method":"CARD"}`,
			svc:        &fakeService{transactionUUID: "tx-1"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "неизвестный способ оплаты — 400",
			body:       `{"payment_method":"NALICHNIE"}`,
			svc:        &fakeService{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "повторная оплата — 409",
			body:       `{"payment_method":"SBP"}`,
			svc:        &fakeService{payErr: model.ErrOrderNotPendingPayment},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "платёжный сервис недоступен — 502",
			body:       `{"payment_method":"SBP"}`,
			svc:        &fakeService{payErr: model.ErrPaymentFailed},
			wantStatus: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := do(t, tt.svc, http.MethodPost, "/api/v1/orders/o1/pay", tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

func TestCancelOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		svc        *fakeService
		wantStatus int
	}{
		{name: "отмена проходит", svc: &fakeService{}, wantStatus: http.StatusNoContent},
		{name: "оплаченный заказ — 409", svc: &fakeService{cancelErr: model.ErrOrderAlreadyPaid}, wantStatus: http.StatusConflict},
		{name: "нет заказа — 404", svc: &fakeService{cancelErr: model.ErrOrderNotFound}, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := do(t, tt.svc, http.MethodPost, "/api/v1/orders/o1/cancel", "")

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}

// Заголовок с идентификатором запроса должен возвращаться клиенту:
// по нему потом ищут логи всех сервисов, участвовавших в запросе.
func TestRequestIDHeader(t *testing.T) {
	t.Parallel()

	rec := do(t, &fakeService{}, http.MethodGet, "/health", "")

	assert.NotEmpty(t, rec.Header().Get("X-Request-Id"))
}
