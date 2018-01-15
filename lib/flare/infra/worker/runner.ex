defmodule Flare.Infra.Worker.Runner do
  use GenServer
  require Logger

  def start_link(opts \\ []) do
    GenServer.start_link(__MODULE__, opts)
  end

  def init(_) do
    Logger.info("Starting a new Worker Runner")
    {:ok, []}
  end
end
