package server

import (
	"catalog/cmd/indexing"
	categoryProducer "catalog/cmd/producers/category"
	productProducer "catalog/cmd/producers/product"
	propertyProducer "catalog/cmd/producers/property"
	"catalog/internal/controllers/v1/category"
	"catalog/internal/controllers/v1/filter"
	"catalog/internal/controllers/v1/product"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	/** Api **/
	mux.HandleFunc("GET /v1/products", product.GetList)
	mux.HandleFunc("GET /v1/categories", category.GetList)
	mux.HandleFunc("GET /v1/filters", filter.GetList)

	/** Создание индексов и отправка тестовых данных в кафку **/
	mux.HandleFunc("GET /createIndexes", indexing.Run)
	mux.HandleFunc("GET /sendCategories", categoryProducer.SendMessages)
	mux.HandleFunc("GET /sendProducts", productProducer.SendMessages)
	mux.HandleFunc("GET /sendProperties", propertyProducer.SendMessages)

	return mux
}
