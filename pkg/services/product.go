package services

import (
	"github.com/storyofhis/go-grpc-product-svc/pkg/db"
	"context"
	pb "github.com/storyofhis/go-grpc-product-svc/pkg/pb"
	"github.com/storyofhis/go-grpc-product-svc/pkg/models"
	"net/http"
)

type Server struct {
	H db.Handler
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var product models.Product

	// request
	product.Name = req.Name
	product.Stock = req.Stock
	product.Price = req.Price

	if res := s.H.DB.Create(&product); res.Error != nil {
		return &pb.CreateProductResponse{
			Status: http.StatusConflict, 
			Error: res.Error.Error(),
		}, nil
	}

	return &pb.CreateProductResponse{
		Status: http.StatusCreated,
		Id: product.Id,
	}, nil
}

func (s *Server) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	var product models.Product

	if res := s.H.DB.First(&product, req.Id); res.Error != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound, 
			Error: res.Error.Error(),
		}, nil
	}

	// request 
	data := &pb.FindOneData {
		Id: product.Id,
		Name: product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &pb.FindOneResponse{
		Status: http.StatusOK,
		Data: data,
	}, nil
}

func (s *Server) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
    var product models.Product

    if result := s.H.DB.First(&product, req.Id); result.Error != nil {
        return &pb.DecreaseStockResponse{
            Status: http.StatusNotFound,
            Error:  result.Error.Error(),
        }, nil
    }

    if product.Stock <= 0 {
        return &pb.DecreaseStockResponse{
            Status: http.StatusConflict,
            Error:  "Stock too low",
        }, nil
    }

    var log models.StockDecreaseLog

    if result := s.H.DB.Where(&models.StockDecreaseLog{OrderId: req.OrderId}).First(&log); result.Error == nil {
        return &pb.DecreaseStockResponse{
            Status: http.StatusConflict,
            Error:  "Stock already decreased",
        }, nil
    }

    product.Stock = product.Stock - 1

    s.H.DB.Save(&product)

    log.OrderId = req.OrderId
    log.ProductRefer = product.Id

    s.H.DB.Create(&log)

    return &pb.DecreaseStockResponse{
        Status: http.StatusOK,
    }, nil
}
