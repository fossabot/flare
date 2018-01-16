use Mix.Config

config :flare, Flare.Infra.HTTP.Endpoint,
  http: [port: 4001],
  server: false

config :logger, level: :warn

config :flare, :repository,
  provider: Memory