// Copyright (c) James Hiester
// SPDX-License-Identifier: MPL-2.0

package secret

type createSecretDataDto struct {
	ID string `json:"secret_id"`
}

type createSecretResponseDto struct {
	Data createSecretDataDto `json:"data"`
}

type secretConfigDto struct {
	URL             string             `json:"url"`
	Method          string             `json:"method"`
	Response        string             `json:"response"`
	Body            *string            `json:"body,omitempty"`
	RefreshInterval int32              `json:"refresh_interval"`
	Headers         *map[string]string `json:"headers,omitempty"`
}

type createSecretDto struct {
	Name   string           `json:"name"`
	Value  *string          `json:"secret,omitempty"`
	Config *secretConfigDto `json:"config,omitempty"`
}

type readSecretDataDto struct {
	Secret secretDto `json:"secret"`
}

type readSecretDto struct {
	Data readSecretDataDto `json:"data"`
}

type updateSecretDto struct {
	Name   string           `json:"name"`
	Value  *string          `json:"value,omitempty"`
	Config *secretConfigDto `json:"config,omitempty"`
}

type secretDto struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Value  *string          `json:"value,omitempty"`
	Config *secretConfigDto `json:"config,omitempty"`
}
