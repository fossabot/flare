use Mix.Config

config :flare, :log,
  level: "debug",
  output: "stdout",
  format: "human"

config :flare, :http,
  server: [
    addr: ":8080",
    default_limit: 30,
    timeout: "1s"
  ],
  client: [
    max_idle_connections: 100
  ]

config :flare, :repository, provider: "mongodb"

config :flare, :worker,
  provider: "aws.sqs",
  subscription: [
    partition: [
      concurrency: 10,
      concurrency_output: 100
    ],
    spread: [
      concurrency: 100,
      concurrency_output: 100,
      concurrency_repository: 100
    ],
    delivery: [
      concurrency: 1000
    ]
  ]

config :flare, :provider,
  aws: [
    session: [
      key: "key",
      secret: "secret",
      region: "region"
    ],
    sqs: [
      queue: [
        [
          worker: "subscription.partition",
          visibility_timeout: "30s",
          retention_period: "10s",
          max_message_size: "100",
          delivery_delay: "10s",
          receive_wait_time: "20s"
        ],
        [
          worker: "subscription.spread"
        ],
        [
          worker: "subscription.delivery"
        ]
      ]
    ]
  ],
  mongodb: [
    addrs: ["localhost:27017"],
    database: "flare",
    username: "flare",
    password: "flare",
    replica_set: "",
    pool_limit: 4096,
    timeout: "1s"
  ]
