package order

import (
	"context"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

const (
	orderURL = "/:id"
)

type handler struct {
	order  Repository
	logger *logrus.Logger
}

func NewHandler(repository *Repository, logger *logrus.Logger) *handler {
	return &handler{
		order:  *repository,
		logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(orderURL, h.GetOrderByUUID)
}

func (h *handler) GetOrderByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	order, err := h.order.FindOne(context.TODO(), params.ByName("id"))
	if err != nil {
		w.WriteHeader(400)
		h.logger.Fatal(err)
	}
	html, err := template.ParseFiles("/home/s1ovac/github.com/wb-service/order-subscribe/view/view.html")
	if err != nil {
		h.logger.Fatal("Can't parse html view file")
	}
	// allBytes, err := json.Marshal(order)
	// if err != nil {
	// 	h.logger.Fatal(err)
	// }

	w.WriteHeader(http.StatusOK)
	err = html.Execute(w, order)
	if err != nil {
		h.logger.Fatal("Can't execute html view file")
	}
}
