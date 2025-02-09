package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/subscription-tracker/subscription/internal/core/domain"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&domain.Subscription{}, &domain.SubscriptionConfig{}, &domain.SubscriptionConfigPlan{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
