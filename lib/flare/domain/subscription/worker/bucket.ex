defmodule Flare.Domain.Subscription.Worker.Bucket do
  # @repository Flare.Provider.repository(:resource)

  def perform(%{resourceId: resourceId}) do
    IO.inspect(resourceId)
    # count = @repository.bucketsCount(resourceId) - 1

    # Task.async_stream(
    #   0..count,
    #   __MODULE__,
    #   :enqueueBucket,
    #   [resourceId: resourceId],
    #   max_concurrency: 10
    # )
    # |> Enum.any(fn x ->
    #   case x do
    #     {:error, _} -> true
    #     _ -> nil
    #   end
    # end)

    # Enum.each(0..count, fn i ->
    #   ExAws.SQS.send_message("flare-subscription-dispatcher", "#{i}") |> ExAws.request()
    # end)
  end

  def enqueueBucket(_i, _opts) do
    # ExAws.SQS.send_message("flare-subscription-dispatcher", "#{i}") |> ExAws.request()
    {:exit, "deu merda"}
    # IO.puts("antes de dar um quit no processo...")
    # Process.exit(self(), :exit)
  end
end

# Flare.Domain.Subscription.Worker.Bucket.perform(%{resourceId: "9799aabe-81e9-49c2-a49b-99817751d036"})
