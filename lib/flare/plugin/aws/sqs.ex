defmodule Flare.Plugin.AWS.SQS do
  use GenServer
  require Logger

  def start_link(opts \\ []) do
    GenServer.start_link(__MODULE__, opts)
  end

  # [concurrency: 10, runner: Flare.Domain.Subscription.Bucket]
  def init(state) do
    # tem que ver se a fila existe antes de rodar, se nao existir, criar.
    Logger.info("Starting a new Worker Runner SQS")
    schedule()
    {:ok, state}
  end

  def handle_info(:work, state) do
    IO.inspect(state)

    ExAws.SQS.receive_message("flare-subscription-dispatcher", wait_time_seconds: 20)
    |> ExAws.request()
    |> IO.inspect()

    schedule()
    {:noreply, state}
  end

  def schedule() do
    Process.send_after(self(), :work, 1000)
  end
end
