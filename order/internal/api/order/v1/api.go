package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/order/internal/service"
	"github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/logger"
	middleware "github.com/BakhytzhanulyE/microservices-course-BakhytzhanulyE/platform/pkg/middleware/http"
)

// API — HTTP-обработчики заказов.
type API struct {
	orderService service.OrderService
}

// NewAPI создаёт HTTP-слой поверх бизнес-логики заказов.
func NewAPI(orderService service.OrderService) *API {
	return &API{orderService: orderService}
}

// Router собирает маршрутизатор со всеми middleware и ручками.
func (a *API) Router() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logging)
	router.Use(middleware.Metrics(routePattern))
	router.Use(middleware.Recoverer)

	router.Get("/health", a.Health)

	router.Route("/api/v1/orders", func(r chi.Router) {
		r.Post("/", a.CreateOrder)
		r.Get("/{order_uuid}", a.GetOrder)
		r.Post("/{order_uuid}/pay", a.PayOrder)
		r.Post("/{order_uuid}/cancel", a.CancelOrder)
	})

	return router
}

// routePattern отдаёт шаблон маршрута вместо конкретного пути.
// Без этого каждый UUID стал бы отдельной меткой в метриках.
func routePattern(r *http.Request) string {
	if rctx := chi.RouteContext(r.Context()); rctx != nil {
		return rctx.RoutePattern()
	}

	return ""
}

// Health отвечает 200 OK — простейшая проверка живости сервиса.
func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("OK")); err != nil {
		logger.Error(r.Context(), "health: не удалось записать ответ", zap.Error(err))
	}
}
