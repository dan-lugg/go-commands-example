package app

import (
	"fmt"
	"github.com/dan-lugg/go-commands-example/app/domain"
	"github.com/dan-lugg/go-commands-example/app/util"
	"github.com/dan-lugg/go-commands/commands"
	"github.com/sarulabs/di/v2"
)

func BuildContainer() (err error, container di.Container) {
	builder, err := di.NewBuilder()
	if err != nil {
		return fmt.Errorf("failed to create DI builder: %w", err), di.Container{}
	}

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[domain.AddHandler](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			return &domain.AddHandler{}, nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[domain.SubHandler](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			return &domain.SubHandler{}, nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[domain.WaitHandler](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			return &domain.WaitHandler{}, nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[domain.NestedHandler](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			return domain.NewNestedHandler(
				ctn.Get(util.TypeNameFor[commands.HandlerCatalog]()).(*commands.HandlerCatalog),
			), nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[commands.MappingCatalog](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			return commands.NewMappingCatalog(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[commands.DecoderCatalog](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			return commands.NewDecoderCatalog(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[commands.HandlerCatalog](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			return commands.NewHandlerCatalog(), nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[commands.Manager](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			manager := commands.NewManager(
				ctn.Get(util.TypeNameFor[commands.MappingCatalog]()).(*commands.MappingCatalog),
				ctn.Get(util.TypeNameFor[commands.DecoderCatalog]()).(*commands.DecoderCatalog),
				ctn.Get(util.TypeNameFor[commands.HandlerCatalog]()).(*commands.HandlerCatalog),
			)
			commands.Insert[domain.AddCommandReq, domain.AddCommandRes](
				manager,
				"add",
				commands.DefaultDecoder[domain.AddCommandReq](),
				func() commands.Handler[domain.AddCommandReq, domain.AddCommandRes] {
					return util.GetFromContainer[*domain.AddHandler](ctn)
				})
			commands.Insert[domain.SubCommandReq, domain.SubCommandRes](
				manager,
				"sub",
				commands.DefaultDecoder[domain.SubCommandReq](),
				func() commands.Handler[domain.SubCommandReq, domain.SubCommandRes] {
					return util.GetFromContainer[*domain.SubHandler](ctn)
				})
			commands.Insert[domain.WaitCommandReq, domain.WaitCommandRes](
				manager,
				"wait",
				commands.DefaultDecoder[domain.WaitCommandReq](),
				func() commands.Handler[domain.WaitCommandReq, domain.WaitCommandRes] {
					return util.GetFromContainer[*domain.WaitHandler](ctn)
				})
			commands.Insert[domain.NestedCommandReq, domain.NestedCommandRes](
				manager,
				"nested",
				commands.DefaultDecoder[domain.NestedCommandReq](),
				func() commands.Handler[domain.NestedCommandReq, domain.NestedCommandRes] {
					return util.GetFromContainer[*domain.NestedHandler](ctn)
				})
			return manager, nil
		},
	})
	
	return nil, builder.Build()
}
