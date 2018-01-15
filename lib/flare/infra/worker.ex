defmodule Flare.Infra.Worker do
  # :discard
  # :ok
  # :error
  # 
  # preciso usar o poolboy para poder fazer o worker pool.

  use Supervisor
  require Logger

  def start_link do
    Supervisor.start_link(__MODULE__, [], name: Flare.Infra.Worker)
  end

  def init(_) do
    Logger.info("Starting Workers")

    children = [
      worker(Flare.Infra.Worker.Runner, [])
    ]

    opts = [strategy: :one_for_one, name: __MODULE__]
    supervise(children, opts)
  end
end
