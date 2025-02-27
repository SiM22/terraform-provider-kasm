//go:build unit
// +build unit

package stats

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"terraform-provider-kasm/internal/client"
)

// mockClient is a mock implementation of the client.Client interface
type mockClient struct {
	mock.Mock
}

func (m *mockClient) GetFrameStats(kasmID string, userID string) (*client.FrameStatsResponse, error) {
	args := m.Called(kasmID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*client.FrameStatsResponse), args.Error(1)
}

func TestStatsResource_Schema(t *testing.T) {
	resource := NewStatsResource()
	_ = resource
}

func TestStatsResource_Configure(t *testing.T) {
	resource := NewStatsResource()
	_ = resource
}

func TestStatsResource_Metadata(t *testing.T) {
	resource := NewStatsResource()
	_ = resource
}
