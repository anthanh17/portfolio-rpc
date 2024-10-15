package gapi

import (
	"context"
	"fmt"
	"portfolio-profile-rpc/rd_portfolio_rpc"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Row struct {
	Headers []string    `bson:"headers" json:"headers"`
	Values  [][]float64 `bson:"values" json:"values"`
}

type TickerInfo struct {
	Name        string  `json:"name" bson:"name"`
	CompanyName string  `json:"company_name"  bson:"company_name"`
	Allocation  float64 `json:"allocation"  bson:"allocation"`
}

type PortfolioOptimization struct {
	ID                      string `bson:"_id,omitempty" json:"_id,omitempty"`
	PortfolioType           string `bson:"portfolio_type" json:"portfolio_type"`
	OptPortfolios           Row    `bson:"opt_portfolios" json:"opt_portfolios"`
	Cml                     Row    `bson:"cml" json:"cml"`
	MaxSharpeRatioPortfolio struct {
		Tickers           []TickerInfo `bson:"tickers" json:"tickers"`
		ExpectedReturn    float64      `bson:"expected_return" json:"expected_return"`
		StandardDeviation float64      `bson:"standard_deviation" json:"standard_deviation"`
		SharpeRatio       float64      `bson:"sharpe_ratio" json:"sharpe_ratio"`
	} `bson:"max_sharpe_ratio_portfolio" json:"max_sharpe_ratio_portfolio"`
	ProvidedPortfolio struct {
		Tickers           []TickerInfo `bson:"tickers" json:"tickers"`
		ExpectedReturn    float64      `bson:"expected_return" json:"expected_return"`
		StandardDeviation float64      `bson:"standard_deviation" json:"standard_deviation"`
		SharpeRatio       float64      `bson:"sharpe_ratio" json:"sharpe_ratio"`
	} `bson:"provided_portfolio" json:"provided_portfolio"`
	Assets []struct {
		Asset             string  `bson:"asset" json:"asset"`
		MinWeight         float64 `bson:"min_weight" json:"min_weight"`
		MaxWeight         float64 `bson:"max_weight" json:"max_weight"`
		ExpectedReturn    float64 `bson:"expected_return" json:"expected_return"`
		StandardDeviation float64 `bson:"standard_deviation" json:"standard_deviation"`
		SharpeRatio       float64 `bson:"sharpe_ratio" json:"sharpe_ratio"`
	} `bson:"assets" json:"assets"`
	CorrMatrix Row `bson:"corr_matrix" json:"corr_matrix"`
	Portfolios Row `bson:"portfolios" json:"portfolios"`
}

func (s *Server) convertPortfolioOptimizationToResp(po PortfolioOptimization) *rd_portfolio_rpc.GetPortfolioOptRes {
	var resp rd_portfolio_rpc.GetPortfolioOptRes
	// Map PortfolioType
	resp.PortfolioType = po.PortfolioType

	// Map ProvidedPortfolio
	resp.ProvidedPortfolio = &rd_portfolio_rpc.Portfolio{
		ExpectedReturn:    po.ProvidedPortfolio.ExpectedReturn,
		StandardDeviation: po.ProvidedPortfolio.StandardDeviation,
		SharpeRatio:       po.ProvidedPortfolio.SharpeRatio,
	}

	// Map Portfolios
	resp.Portfolios = &rd_portfolio_rpc.Row{
		Headers: po.Portfolios.Headers,
	}
	for _, value := range po.Portfolios.Values {
		var _value *rd_portfolio_rpc.Value = &rd_portfolio_rpc.Value{
			Values: value,
		}
		resp.Portfolios.Values = append(resp.Portfolios.Values, _value)
	}
	// Map Ticker
	for _, ticker := range po.ProvidedPortfolio.Tickers {
		resp.ProvidedPortfolio.Tickers = append(resp.ProvidedPortfolio.Tickers, &rd_portfolio_rpc.TickerInfo{
			Name:        ticker.Name,
			CompanyName: ticker.CompanyName,
			Allocation:  ticker.Allocation,
		})
	}
	// Map MaxSharpeRatioPortfolio
	resp.MaxSharpeRatioPortfolio = &rd_portfolio_rpc.Portfolio{
		ExpectedReturn:    po.MaxSharpeRatioPortfolio.ExpectedReturn,
		StandardDeviation: po.MaxSharpeRatioPortfolio.StandardDeviation,
		SharpeRatio:       po.MaxSharpeRatioPortfolio.SharpeRatio,
	}

	// Map Ticker
	for _, ticker := range po.MaxSharpeRatioPortfolio.Tickers {
		resp.MaxSharpeRatioPortfolio.Tickers = append(resp.MaxSharpeRatioPortfolio.Tickers, &rd_portfolio_rpc.TickerInfo{
			Name:        ticker.Name,
			CompanyName: ticker.CompanyName,
			Allocation:  ticker.Allocation,
		})
	}

	// Map RandomPortfolios
	resp.OptPortfolios = &rd_portfolio_rpc.Row{
		Headers: po.OptPortfolios.Headers,
	}
	for _, value := range po.OptPortfolios.Values {
		var _value *rd_portfolio_rpc.Value = &rd_portfolio_rpc.Value{
			Values: value,
		}
		resp.OptPortfolios.Values = append(resp.OptPortfolios.Values, _value)
	}

	// Map CorrMatrix
	resp.CorrMatrix = &rd_portfolio_rpc.Row{
		Headers: po.CorrMatrix.Headers,
	}
	for _, value := range po.CorrMatrix.Values {
		var _value *rd_portfolio_rpc.Value = &rd_portfolio_rpc.Value{
			Values: value,
		}
		resp.CorrMatrix.Values = append(resp.CorrMatrix.Values, _value)
	}

	// Map Cml
	resp.Cml = &rd_portfolio_rpc.Row{
		Headers: po.Cml.Headers,
	}
	for _, value := range po.Cml.Values {
		var _value *rd_portfolio_rpc.Value = &rd_portfolio_rpc.Value{
			Values: value,
		}
		resp.Cml.Values = append(resp.Cml.Values, _value)
	}

	// Map Assets
	for _, asset := range po.Assets {
		resp.Assets = append(resp.Assets, &rd_portfolio_rpc.Asset{
			Asset:             asset.Asset,
			MinWeight:         asset.MinWeight,
			MaxWeight:         asset.MaxWeight,
			ExpectedReturn:    asset.ExpectedReturn,
			StandardDeviation: asset.StandardDeviation,
			SharpeRatio:       asset.SharpeRatio,
		})
	}

	return &resp
}

func (s *Server) Ping(ctx context.Context, in *rd_portfolio_rpc.Request) (*rd_portfolio_rpc.Response, error) {
	fmt.Printf("\n==> Ping")
	return &rd_portfolio_rpc.Response{}, nil
}

func (s *Server) GetOpt(ctx context.Context, in *rd_portfolio_rpc.GetPortfolioOptReq) (*rd_portfolio_rpc.GetPortfolioOptRes, error) {
	fmt.Printf("\n---> GetOpt request: { %v }\n", in)

	collection := s.mongoClient.Database(s.config.Mongo.Database).Collection(s.config.Mongo.Collection)
	var result PortfolioOptimization
	err := collection.FindOne(
		context.Background(), bson.M{"_id": in.PoId},
	).Decode(&result)
	if err != nil {
		s.logger.Sugar().Infof("\nerror findone: %v\n", err)
		return nil, status.Errorf(codes.Internal, "failed to findone mongo: %s", err)
	}

	fmt.Printf("\n==> GetOpt PoID: %s - UserID: %s", in.PoId, in.UserId)
	return s.convertPortfolioOptimizationToResp(result), nil
}
