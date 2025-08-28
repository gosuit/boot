# Boot

Boot is a library that helps to inject dependencies.

We separate module dependencies into other modules and configurations. Di containers handle the injection of other modules well, but there are certain problems with the configuration. We are in favor of centralized configuration reading into a single structure that contains structures with the configs needed for the application. But the di container cannot "extract" the config needed by the constructor from the general config. Boot helps to solve this problem.

## Installation

```zsh
go get github.com/gosuit/boot
```

## Usage

```golang
package module

// main/module/module.go

import "github.com/gosuit/boot"

type ModuleConfig struct {
	Value string
}

type Module struct {
	Cfg *ModuleConfig
}

func Boot[G any]() any {
	return boot.Boot[G, ModuleConfig](NewModule)
}

func NewModule(cfg *ModuleConfig) *Module {
	return &Module{
		Cfg: cfg,
	}
}

```

```golang
package main

// main/main.go

import (
	"fmt"

	"github.com/gosuit/boot/main/module"
	"go.uber.org/fx"
)

type AppConfig struct {
	Module *module.ModuleConfig
}

func main() {
	fx.New(
		fx.Provide(NewConfig),
		fx.Provide(module.Boot[AppConfig]()),
		fx.Invoke(func(mod *module.Module) {
			fmt.Println(mod.Cfg) // &{value}
		}),
	).Run()
}

func NewConfig() *AppConfig {
	return &AppConfig{
		Module: &module.ModuleConfig{
			Value: "value",
		},
	}
}
```

### Module injection

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
