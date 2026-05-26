package main

import (
	"encoding/json"
	"errors"
	"log/slog"

	rhsmapi "github.com/jirihnidek/rhsm-server/varlink/com/redhat/rhsm"
	"github.com/jirihnidek/rhsm2"
)

var AppName = "rhsm-server"

type ClientError struct {
	Message string
}

func (e *ClientError) Error() string {
	return e.Message
}

type ServerError struct {
	Message string
}

func (e *ServerError) Error() string {
	return e.Message
}

// GetStatus retrieves the current status of the Red Hat Subscription Management (RHSM) server.
func GetStatus(ipcSender *string, locale *string, correlationID *string) (*rhsm2.RHSMStatus, error) {
	rhsmClient, err := rhsm2.GetRHSMClient(&AppName, nil)
	if err != nil {
		return nil, &ClientError{Message: err.Error()}
	}

	// Create client information from provided parameters
	clientInfo := rhsm2.RequestMetadata{IPCSender: ipcSender, Locale: locale, CorrelationId: correlationID}
	status, err := rhsmClient.GetServerStatus(&clientInfo)
	if err != nil {
		return nil, &ServerError{Message: err.Error()}
	}

	return status, nil
}

// IsSystemRegistered checks if the system is registered with RHSM.
// When it is not possible to retrieve the consumer UUID, it returns false.
func IsSystemRegistered() (bool, error) {
	rhsmClient, err := rhsm2.GetRHSMClient(&AppName, nil)
	if err != nil {
		return false, &ClientError{Message: err.Error()}
	}

	_, err = rhsmClient.GetConsumerUUID()
	if err != nil {
		return false, err
	}

	return true, nil
}

type RhsmBackend struct{}

func NewRhsmBackend() *RhsmBackend {
	return &RhsmBackend{}
}

// Ping checks the status of the RHSM server.
func (b *RhsmBackend) Ping(in *rhsmapi.PingIn) (*rhsmapi.PingOut, error) {
	slog.Debug("Ping() method called")
	var rhsmServerStatus *rhsm2.RHSMStatus
	var err error
	if in.Metadata != nil {
		rhsmServerStatus, err = GetStatus(
			in.Metadata.UserAgent,
			in.Metadata.Locale,
			in.Metadata.CorrelationId,
		)
	} else {
		rhsmServerStatus, err = GetStatus(nil, nil, nil)
	}
	if err != nil {
		var typeClientErr *ClientError
		var typeServerErr *ServerError
		switch {
		case errors.As(err, &typeClientErr):
			return nil, &rhsmapi.InvalidClientConnectionError{Message: typeClientErr.Message}
		case errors.As(err, &typeServerErr):
			return nil, &rhsmapi.FailedServerResponseError{Message: typeServerErr.Message}
		default:
			slog.Error("Failed to get RHSM status", "error", err)
			return nil, err
		}
	}
	status, err := json.Marshal(rhsmServerStatus)
	if err != nil {
		return nil, &rhsmapi.FailedServerResponseError{Message: err.Error()}
	}
	slog.Debug("RHSM status retrieved", "status", string(status))
	return &rhsmapi.PingOut{Status: status}, nil
}

// IsRegistered checks if the system is registered with RHSM.
func (b *RhsmBackend) IsRegistered(in *rhsmapi.IsRegisteredIn) (*rhsmapi.IsRegisteredOut, error) {
	slog.Debug("IsRegistered() method called")
	registered, err := IsSystemRegistered()
	if err != nil {
		// When it is not possible to determine registration status, then log the reason
		// and return false
		slog.Debug("Failed to determine registration status", "error", err)
		return &rhsmapi.IsRegisteredOut{Registered: false}, nil
	}
	slog.Debug("System registration status determined", "registered", registered)
	return &rhsmapi.IsRegisteredOut{Registered: registered}, nil
}
