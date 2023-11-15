package util

import (
	"github.com/corbado/corbado-go/pkg/sdk/entity/common"
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

func ClientInfo(userAgent string, remoteAddress string) *common.ClientInfo {
	if userAgent == "" && remoteAddress == "" {
		return nil
	}

	return &common.ClientInfo{RemoteAddress: remoteAddress, UserAgent: userAgent}
}
