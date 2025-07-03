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
				ctn.Get(util.TypeNameFor[commands.HandlerRegistry]()).(*commands.HandlerRegistry),
			), nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[commands.MappingRegistry](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			registry := commands.NewMappingRegistry()
			commands.RegisterMapping[domain.AddCommandReq](registry, "add")
			commands.RegisterMapping[domain.SubCommandReq](registry, "sub")
			commands.RegisterMapping[domain.WaitCommandReq](registry, "wait")
			commands.RegisterMapping[domain.NestedCommandReq](registry, "nested")
			return registry, nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[commands.DecoderRegistry](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			registry := commands.NewDecoderRegistry()
			commands.RegisterDecoder[domain.AddCommandReq](registry, commands.DefaultCommandReqDecoder[domain.AddCommandReq]())
			commands.RegisterDecoder[domain.SubCommandReq](registry, commands.DefaultCommandReqDecoder[domain.SubCommandReq]())
			commands.RegisterDecoder[domain.WaitCommandReq](registry, commands.DefaultCommandReqDecoder[domain.WaitCommandReq]())
			commands.RegisterDecoder[domain.NestedCommandReq](registry, commands.DefaultCommandReqDecoder[domain.NestedCommandReq]())
			return registry, nil
		},
	})

	_ = builder.Add(di.Def{
		Name:  util.TypeNameFor[commands.HandlerRegistry](),
		Scope: di.App,
		Build: func(ctn di.Container) (any, error) {
			registry := commands.NewHandlerRegistry()
			commands.RegisterHandler[domain.AddCommandReq, domain.AddCommandRes](registry, func() commands.Handler[domain.AddCommandReq, domain.AddCommandRes] {
				return util.GetFromContainer[*domain.AddHandler](ctn)
			})
			commands.RegisterHandler[domain.SubCommandReq, domain.SubCommandRes](registry, func() commands.Handler[domain.SubCommandReq, domain.SubCommandRes] {
				return util.GetFromContainer[*domain.SubHandler](ctn)
			})
			commands.RegisterHandler[domain.WaitCommandReq, domain.WaitCommandRes](registry, func() commands.Handler[domain.WaitCommandReq, domain.WaitCommandRes] {
				return util.GetFromContainer[*domain.WaitHandler](ctn)
			})
			commands.RegisterHandler[domain.NestedCommandReq, domain.NestedCommandRes](registry, func() commands.Handler[domain.NestedCommandReq, domain.NestedCommandRes] {
				return util.GetFromContainer[*domain.NestedHandler](ctn)
			})
			return registry, nil
		},
	})

	return nil, builder.Build()
}
