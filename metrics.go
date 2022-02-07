package dtrack

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type PortfolioMetrics struct {
	FirstOccurrence int `json:"firstOccurrence"`
	LastOccurrence  int `json:"lastOccurrence"`

	InheritedRiskScore   float64 `json:"inheritedRiskScore"`
	Vulnerabilities      int     `json:"vulnerabilities"`
	VulnerableProjects   int     `json:"vulnerableProjects"`
	VulnerableComponents int     `json:"vulnerableComponents"`
	Projects             int     `json:"projects"`
	Components           int     `json:"components"`
	Suppressed           int     `json:"suppressed"`

	Critical   int `json:"critical"`
	High       int `json:"high"`
	Medium     int `json:"medium"`
	Low        int `json:"low"`
	Unassigned int `json:"unassigned"`

	FindingsTotal     int `json:"findingsTotal"`
	FindingsAudited   int `json:"findingsAudited"`
	FindingsUnaudited int `json:"findingsUnaudited"`

	PolicyViolationsTotal     int `json:"policyViolationsTotal"`
	PolicyViolationsFail      int `json:"policyViolationsFail"`
	PolicyViolationsWarn      int `json:"policyViolationsWarn"`
	PolicyViolationsInfo      int `json:"policyViolationsInfo"`
	PolicyViolationsAudited   int `json:"policyViolationsAudited"`
	PolicyViolationsUnaudited int `json:"policyViolationsUnaudited"`

	PolicyViolationsSecurityTotal     int `json:"policyViolationsSecurityTotal"`
	PolicyViolationsSecurityAudited   int `json:"policyViolationsSecurityAudited"`
	PolicyViolationsSecurityUnaudited int `json:"policyViolationsSecurityUnaudited"`

	PolicyViolationsLicenseTotal     int `json:"policyViolationsLicenseTotal"`
	PolicyViolationsLicenseAudited   int `json:"policyViolationsLicenseAudited"`
	PolicyViolationsLicenseUnaudited int `json:"policyViolationsLicenseUnaudited"`

	PolicyViolationsOperationalTotal     int `json:"policyViolationsOperationalTotal"`
	PolicyViolationsOperationalAudited   int `json:"policyViolationsOperationalAudited"`
	PolicyViolationsOperationalUnaudited int `json:"policyViolationsOperationalUnaudited"`
}

type MetricsService struct {
	client *Client
}

func (m MetricsService) LatestPortfolioMetrics(ctx context.Context) (*PortfolioMetrics, error) {
	req, err := m.client.newRequest(ctx, http.MethodGet, "/api/v1/metrics/portfolio/current")
	if err != nil {
		return nil, err
	}

	var metrics PortfolioMetrics
	_, err = m.client.doRequest(req, &metrics)
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}

func (m MetricsService) PortfolioMetricsSince(ctx context.Context, date time.Time) ([]PortfolioMetrics, error) {
	req, err := m.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/metrics/portfolio/since/%s", date.Format("20060102")))
	if err != nil {
		return nil, err
	}

	var metrics []PortfolioMetrics
	_, err = m.client.doRequest(req, &metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (m MetricsService) PortfolioMetricsSinceDays(ctx context.Context, days uint) ([]PortfolioMetrics, error) {
	req, err := m.client.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/metrics/portfolio/%d/days", days))
	if err != nil {
		return nil, err
	}

	var metrics []PortfolioMetrics
	_, err = m.client.doRequest(req, &metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (m MetricsService) RefreshPortfolioMetrics(ctx context.Context) error {
	req, err := m.client.newRequest(ctx, http.MethodGet, "/api/v1/metrics/portfolio/refresh")
	if err != nil {
		return err
	}

	_, err = m.client.doRequest(req, nil)
	if err != nil {
		return err
	}

	return nil
}
