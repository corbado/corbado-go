package util

import (
	"github.com/corbado/corbado-go/pkg/generated/common"
)

// Ptr returns a pointer to provided value
func Ptr[T any](v T) *T {
	return &v
}

// RequestID creates an optional Request ID
func RequestID(requestID string) *common.RequestID {
	if requestID != "" {
		return Ptr(requestID)
	}

	return nil
}

// ClientInfo returns client info based on provided user agent and remote address
func ClientInfo(userAgent string, remoteAddress string) *common.ClientInfo {
	if userAgent == "" && remoteAddress == "" {
		return nil
	}

	return &common.ClientInfo{RemoteAddress: remoteAddress, UserAgent: userAgent}
}
