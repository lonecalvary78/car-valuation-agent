package handler

import (
	"bytes"
	"car-valuation-agent/internal/infrastructure/middleware"
	"car-valuation-agent/internal/infrastructure/util"
	"car-valuation-agent/internal/model"
	"context"
	"errors"
	"io"
	"iter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/adk/v2/agent"
	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/session"
	"google.golang.org/genai"
)

var errRunnerFailed = errors.New("runner failed")

func newTestHandler(t *testing.T, run func(agent.InvocationContext) iter.Seq2[*session.Event, error]) Handler {
	t.Helper()

	fakeAgent, err := agent.New(agent.Config{
		Name: "car-valuer-agent",
		Run:  run,
	})
	require.NoError(t, err)

	agentRunner, err := util.OfRunner(fakeAgent, memory.InMemoryService(), session.InMemoryService(), true)
	require.NoError(t, err)

	return New(*agentRunner, 6*time.Second)
}

func withAuthenticatedUser(r *http.Request) *http.Request {
	user := middleware.User{ID: generateNewUUID(), Username: "test-user"}
	return r.WithContext(context.WithValue(r.Context(), middleware.UserKey, user))
}

func runOfSingleText(text string) func(agent.InvocationContext) iter.Seq2[*session.Event, error] {
	return func(ic agent.InvocationContext) iter.Seq2[*session.Event, error] {
		return func(yield func(*session.Event, error) bool) {
			event := session.NewEvent(ic, ic.InvocationID())
			event.Content = genai.NewContentFromText(text, genai.RoleModel)
			yield(event, nil)
		}
	}
}

func TestAskToAgent(t *testing.T) {
	t.Run("returns the parsed valuation on a valid request", func(t *testing.T) {
		expectedResponse := constructResponse()
		responseJSON, err := util.ToJSON(expectedResponse)
		require.NoError(t, err)

		h := newTestHandler(t, runOfSingleText(string(responseJSON)))

		body := bytes.NewBufferString(`{"message":"How much for a 2010 Toyota Prius?"}`)
		r := withAuthenticatedUser(httptest.NewRequestWithContext(t.Context(), "POST", "/v1/valuations", body))
		w := httptest.NewRecorder()

		h.askToAgent(w, r)

		require.Equal(t, 200, w.Code)

		responseBody, _ := io.ReadAll(w.Result().Body)
		var carValuationResponse model.CarValuationResponse
		require.NoError(t, util.FromJSON(string(responseBody), &carValuationResponse))
		require.Equal(t, expectedResponse.Brand, carValuationResponse.Brand)
		require.Equal(t, expectedResponse.Model, carValuationResponse.Model)
		require.Equal(t, expectedResponse.Year, carValuationResponse.Year)
		require.Equal(t, expectedResponse.Market, carValuationResponse.Market)
		require.Equal(t, expectedResponse.Price.CurrencyCode, carValuationResponse.Price.CurrencyCode)
		require.True(t, expectedResponse.Price.LowerPrice.Equal(carValuationResponse.Price.LowerPrice))
		require.True(t, expectedResponse.Price.HigherPrice.Equal(carValuationResponse.Price.HigherPrice))
	})

	t.Run("returns 400 when message is missing", func(t *testing.T) {
		runInvoked := false
		h := newTestHandler(t, func(ic agent.InvocationContext) iter.Seq2[*session.Event, error] {
			runInvoked = true
			return func(yield func(*session.Event, error) bool) {}
		})

		body := bytes.NewBufferString(`{}`)
		r := withAuthenticatedUser(httptest.NewRequestWithContext(t.Context(), "POST", "/v1/valuations", body))
		w := httptest.NewRecorder()

		h.askToAgent(w, r)

		require.Equal(t, 400, w.Code)
		require.False(t, runInvoked, "the agent should not run when validation fails")
	})

	t.Run("returns 401 when there is no authenticated user in the request context", func(t *testing.T) {
		runInvoked := false
		h := newTestHandler(t, func(ic agent.InvocationContext) iter.Seq2[*session.Event, error] {
			runInvoked = true
			return func(yield func(*session.Event, error) bool) {}
		})

		body := bytes.NewBufferString(`{"message":"How much for a 2010 Toyota Prius?"}`)
		r := httptest.NewRequestWithContext(t.Context(), "POST", "/v1/valuations", body)
		w := httptest.NewRecorder()

		h.askToAgent(w, r)

		require.Equal(t, 401, w.Code)
		require.False(t, runInvoked, "the agent should not run without an authenticated user")
	})

	t.Run("returns 400 when the request body is not valid JSON", func(t *testing.T) {
		h := newTestHandler(t, runOfSingleText(""))

		body := bytes.NewBufferString("not-json")
		r := withAuthenticatedUser(httptest.NewRequestWithContext(t.Context(), "POST", "/v1/valuations", body))
		w := httptest.NewRecorder()

		h.askToAgent(w, r)

		require.Equal(t, 400, w.Code)
	})

	t.Run("returns 500 when the runner yields an error", func(t *testing.T) {
		h := newTestHandler(t, func(ic agent.InvocationContext) iter.Seq2[*session.Event, error] {
			return func(yield func(*session.Event, error) bool) {
				yield(nil, errRunnerFailed)
			}
		})

		body := bytes.NewBufferString(`{"message":"How much for a 2010 Toyota Prius?"}`)
		r := withAuthenticatedUser(httptest.NewRequestWithContext(t.Context(), "POST", "/v1/valuations", body))
		w := httptest.NewRecorder()

		h.askToAgent(w, r)

		require.Equal(t, 500, w.Code)
	})

	t.Run("returns 500 when the agent text is not valid CarValuationResponse JSON", func(t *testing.T) {
		h := newTestHandler(t, runOfSingleText("not-json"))

		body := bytes.NewBufferString(`{"message":"How much for a 2010 Toyota Prius?"}`)
		r := withAuthenticatedUser(httptest.NewRequestWithContext(t.Context(), "POST", "/v1/valuations", body))
		w := httptest.NewRecorder()

		h.askToAgent(w, r)

		require.Equal(t, 500, w.Code)
	})
}

func generateNewUUID() string {
	return uuid.New().String()
}

func constructResponse() model.CarValuationResponse {
	return model.CarValuationResponse{
		Brand:  "Toyota",
		Model:  "Prius",
		Year:   2010,
		Market: "Malaysia",
		Price: model.QuotedPrice{
			CurrencyCode: "MYR",
			LowerPrice:   decimal.NewFromFloat(30000),
			HigherPrice:  decimal.NewFromFloat(45000),
		},
	}
}
