// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type HamonixBusinessAsset struct {
	ID          int64   `json:"id"`
	PortfolioID string  `json:"portfolio_id"`
	TickerID    int32   `json:"ticker_id"`
	Price       float64 `json:"price"`
	Allocation  float64 `json:"allocation"`
}

type HamonixBusinessEqAccount struct {
	ID        string      `json:"id"`
	AdvisorID pgtype.Text `json:"advisor_id"`
	Code      string      `json:"code"`
}

type HamonixBusinessEqAdvisor struct {
	ID          string      `json:"id"`
	Code        pgtype.Text `json:"code"`
	Description pgtype.Text `json:"description"`
}

type HamonixBusinessEqBackoffice struct {
	ID           string      `json:"id"`
	WhitelableID pgtype.Text `json:"whitelable_id"`
	Name         string      `json:"name"`
	Description  pgtype.Text `json:"description"`
}

type HamonixBusinessEqBranch struct {
	ID          string      `json:"id"`
	Code        string      `json:"code"`
	Description pgtype.Text `json:"description"`
}

type HamonixBusinessEqOrganization struct {
	ID           string      `json:"id"`
	BackofficeID pgtype.Text `json:"backoffice_id"`
	Code         string      `json:"code"`
	Description  pgtype.Text `json:"description"`
}

type HamonixBusinessEqWhitelable struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Url         string      `json:"url"`
	Description pgtype.Text `json:"description"`
}

type HamonixBusinessPAdvisor struct {
	ID          int64       `json:"id"`
	PortfolioID string      `json:"portfolio_id"`
	AdvisorID   pgtype.Text `json:"advisor_id"`
}

type HamonixBusinessPBranch struct {
	ID          int64       `json:"id"`
	PortfolioID string      `json:"portfolio_id"`
	BranchID    pgtype.Text `json:"branch_id"`
}

type HamonixBusinessPCategory struct {
	ID          int64       `json:"id"`
	PortfolioID string      `json:"portfolio_id"`
	CategoryID  pgtype.Text `json:"category_id"`
}

type HamonixBusinessPOrganization struct {
	ID             int64       `json:"id"`
	PortfolioID    string      `json:"portfolio_id"`
	OrganizationID pgtype.Text `json:"organization_id"`
}

type HamonixBusinessPortfolio struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Privacy   string    `json:"privacy"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HamonixBusinessPortfolioCategory struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type HamonixBusinessTicker struct {
	ID          int64     `json:"id"`
	Symbol      string    `json:"symbol"`
	Description string    `json:"description"`
	Exchange    string    `json:"exchange"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type HamonixBusinessTickerPrice struct {
	TickerID int64       `json:"ticker_id"`
	Open     float64     `json:"open"`
	High     float64     `json:"high"`
	Low      float64     `json:"low"`
	Close    float64     `json:"close"`
	Date     pgtype.Date `json:"date"`
}

type HamonixBusinessUCategory struct {
	ID         int64       `json:"id"`
	CategoryID pgtype.Text `json:"category_id"`
	UserID     string      `json:"user_id"`
}

type HamonixBusinessUPortfolio struct {
	ID          int64       `json:"id"`
	UserID      string      `json:"user_id"`
	PortfolioID pgtype.Text `json:"portfolio_id"`
}

type HamonixBusinessUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
