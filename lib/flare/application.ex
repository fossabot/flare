defmodule Flare.Application do
  use Application
  use Supervisor

  def start(_type, _args) do
    Application.put_env(:flare, :diego, :bernardes)

    children =
      [
        supervisor(Flare.Infra.HTTP.Endpoint, []),
        supervisor(Flare.Plugin.Memory.Resource, []),
        supervisor(Flare.Infra.Worker, []),
        worker(Mongo, [[name: :mongo, database: "flare", pool: DBConnection.Poolboy]])
      ] ++ worker_provider()

    opts = [strategy: :one_for_one, name: Flare.Supervisor]
    # {:ok, self()}
    # dps de iniciar, se iniciar, mandar um evento para comecar o procesamento.
    Supervisor.start_link(children, opts)
  end

  def worker_provider() do
    # case Application.get_env(:flare, :worker) do
    #   nil -> []
    #   [provider: provider, consumers: consumers] -> [supervisor(provider, consumers)]
    # end

    []
  end

  def config_change(changed, _new, removed) do
    Flare.Infra.HTTP.Endpoint.config_change(changed, removed)
    :ok
  end
end
