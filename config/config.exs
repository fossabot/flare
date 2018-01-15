use Mix.Config

config :flare, Flare.Infra.HTTP.Endpoint,
  url: [host: "localhost"],
  secret_key_base: "l39nqH24in7GJTm8k6dBmw2J8dtckm1QuH9J1Hz0vjyE1QIkaGftgFi7d8YnjxJ3",
  render_errors: [view: Flare.Infra.HTTP.ErrorView, accepts: ~w(json)]

config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

import_config "#{Mix.env()}.exs"
