module github.com/shipmight/helm-playground

// Based on:
//   https://github.com/helm/helm/blob/v3.18.2/go.mod
// Commands:
//   go get github.com/BurntSushi/toml@v1.5.0
//   go get github.com/Masterminds/sprig/v3@v3.3.0
//   go get sigs.k8s.io/yaml@v1.4.0
// Update these commands with the versions and then run them!
// Also keep go version in sync!

go 1.24

toolchain go1.24.4

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/Masterminds/sprig/v3 v3.3.0
	sigs.k8s.io/yaml v1.4.0
)

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/huandu/xstrings v1.5.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
)
