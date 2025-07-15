// Copyright (c) James Hiester.
// Licensed under the MIT license.

package mocks

import (
	"context"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/jameshiester/terraform-provider-bland/internal/provider"
	utils "github.com/jameshiester/terraform-provider-bland/internal/util"
)

func TestName() string {
	pc, _, _, _ := runtime.Caller(1)
	nameFull := runtime.FuncForPC(pc).Name()
	nameEnd := filepath.Ext(nameFull)
	name := strings.TrimPrefix(nameEnd, ".")
	return name
}

var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"bland": providerserver.NewProtocol6WithError(provider.NewBlandProvider(utils.UnitTestContext(context.Background(), ""), true)()),
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"bland": providerserver.NewProtocol6WithError(provider.NewBlandProvider(context.Background(), false)()),
}

func ActivateEnvironmentHttpMocks() {
}
