package handler

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/s1ovac/order-subscribe/internal/cache"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/s1ovac/order-subscribe/internal/subscribe"
	"github.com/sirupsen/logrus"
)

const (
	orderURL = "/order/:id"
)

type handler struct {
	order      order.Repository
	logger     *logrus.Logger
	cache      *cache.Cache
	subscriber *subscribe.Subscriber
	context    context.Context
}

func NewHandler(ctx context.Context, repository *order.Repository, logger *logrus.Logger, cache *cache.Cache, subscriber *subscribe.Subscriber) *handler {
	return &handler{
		order:      *repository,
		logger:     logger,
		cache:      cache,
		subscriber: subscriber,
		context:    ctx,
	}
}

func (h *handler) Register(ctx context.Context, router *httprouter.Router) {
	router.GET(orderURL, h.GetOrderByUUID)
	router.GET("/static/style.css", h.StaticCSS)
}

func (h *handler) GetOrderByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	order, err := h.cache.GetCache(h.context, params.ByName("id"))
	if errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		html, err := template.ParseFiles("/home/s1ovac/github.com/wb-service/order-subscribe/view/notfound.html")
		if err != nil {
			h.logger.Fatal("Can't parse html view file")
		}
		if err := html.Execute(w, order); err != nil {
			h.logger.Fatal("Can't execute html view file")
		}
		return
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		pgErr = err.(*pgconn.PgError)
		switch pgErr.Code {
		case "23505", "22P02":
			w.WriteHeader(http.StatusNotFound)
			html, err := template.ParseFiles("/home/s1ovac/github.com/wb-service/order-subscribe/view/notfound.html")
			if err != nil {
				h.logger.Fatal("Can't parse html view file")
			}
			if err := html.Execute(w, order); err != nil {
				h.logger.Fatal("Can't execute html view file")
			}
			return
		}

		newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
		h.logger.Error(newErr)
		return
	}
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
	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("order-subscribe/view/style.css"))))

	w.WriteHeader(http.StatusOK)
	err = html.Execute(w, order)
	if err != nil {
		h.logger.Fatal("Can't execute html view file")
	}
}

func (h *handler) StaticCSS(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	http.ServeFile(w, r, "/home/s1ovac/github.com/wb-service/order-subscribe/static/style.css")
}
