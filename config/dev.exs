use Mix.Config

config :flare, Flare.Infra.HTTP.Endpoint,
  http: [port: 8080],
  debug_errors: true,
  code_reloader: true,
  check_origin: false,
  watchers: []

config :flare, :repository,
  provider: MongoDB,
  options: [
    pool: DBConnection.Poolboy
  ],
  resource: [
    bucket_size: 1000
  ]

config :logger, :console, format: "[$level] $message\n"

config :phoenix, :stacktrace_depth, 20

config :ex_aws,
  access_key_id: "AKIAJYORFQYWJNIY7DHQ",
  secret_access_key: "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",
  region: "sa-east-1"

config :flare, :subscription,
  worker: [
    reduce: [
      concurrency: 10
    ],
    dispatcher: [
      concurrency: 20,
      pipeline_concurrency: 10
    ],
    unit: [
      concurrency: 1000
    ]
  ]

# pegar fora do runner oq nao estiver no runner, isso vale pro provider.
config :flare, :worker,
  provider: Flare.Plugin.AWS.SQS,
  producers: [
    [
      runner: Flare.Plugin.AWS.SQS,
      ingress: [
        concurrency: 1,
        rateLimit: [
          [
            period: 1000,
            quantity: 100
          ]
        ]
      ],
      egress: [
        concurrency: 1,
        options: [wait_time_seconds: 20]
      ]
    ]
  ],
  consumers: [
    [
      concurrency: 10,
      queue: "flare-subscription-dispatcher",
      runner: Flare.Domain.Subscription.Bucket,
      options: [wait_time_seconds: 20]
    ]
  ]
