package checker

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labbsr0x/goh/gohclient"

	"github.com/abilioesteves/health-checker/config"
	"github.com/labbsr0x/goh/gohcmd"
	"github.com/sirupsen/logrus"
)

// Checker holds the information needed to run the health checker agent
type Checker struct {
	*config.Builder
	ctx    context.Context
	cancel context.CancelFunc
}

// HealthCheckResponse holds the information for a health check response
type HealthCheckResponse struct {
	Dependencies map[string]DependencyStatus `json:"dependencies"`
}

// DependencyStatus holds the information of a dependency status
type DependencyStatus struct {
	Up  bool   `json:"up"`
	Err string `json:"err"`
}

// InitFromBuilder inits the health checker agent from the Builder object
func (checker *Checker) InitFromBuilder(builder *config.Builder) *Checker {
	checker.Builder = builder
	ctx, cancel := context.WithCancel(context.Background())

	checker.ctx = ctx
	checker.cancel = cancel

	return checker
}

// Run starts the health checker agent
func (checker *Checker) Run() {
	go gohcmd.GracefulStop(checker.cancel)

	duration, _ := time.ParseDuration("15s")
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	logrus.Infof("Started checker agent")
	for {
		select {
		case <-ticker.C:
			response, err := checker.CheckHealth()
			if err != nil {
				logrus.Errorf("Error checking health: %v", err)
				checker.RegisterProblem(err)
			} else {
				checker.RegisterResponse(response)
			}

		case <-checker.ctx.Done():
			logrus.Info("Health checker agent stopped")
			return
		}
	}

}

// RegisterProblem registrates at the appropriate metric that a problem communicating with the target service exists
func (checker *Checker) RegisterProblem(err error) {
	checker.HealthMetric.WithLabelValues(checker.TargetName, "self", err.Error()).Set(0)
}

// RegisterResponse registrates the health check response at the appropriate prometheus metric
func (checker *Checker) RegisterResponse(resp HealthCheckResponse) {
	for k, depStatus := range resp.Dependencies {
		value := 0.0
		if depStatus.Up {
			value = 1.0
		}
		checker.HealthMetric.WithLabelValues(checker.TargetName, k, depStatus.Err).Set(value)
	}
	checker.HealthMetric.WithLabelValues(checker.TargetName, "self", "").Set(1)
}

// CheckHealth calls the health endpoint
func (checker *Checker) CheckHealth() (toReturn HealthCheckResponse, err error) {
	httpClient, err := gohclient.New(nil, checker.TargetHealthURL)

	httpResp, data, err := httpClient.Get("")

	if httpResp != nil {
		if httpResp.StatusCode == http.StatusOK {
			if err == nil {
				if err = json.Unmarshal(data, &toReturn); err == nil {
					return
				}
				err = fmt.Errorf("Health Check '%v': Unable to read response", checker.TargetHealthURL)
			} else {
				err = fmt.Errorf("Health Check '%v': Unable to communicate", checker.TargetHealthURL)
			}
		} else {
			err = fmt.Errorf("Health Check '%v': Not 200 OK; Getting %v", checker.TargetHealthURL, httpResp.StatusCode)
		}
	} else {
		err = fmt.Errorf("Health Check '%v': Not possible to communicate with server: %v", checker.TargetHealthURL, err)
	}

	return
}
