defmodule Flare.Domain.Subscription.Repository do
  @callback all() :: {:ok, [%{}]} | {:error, String.t()}
  @callback create(%{}) :: {:ok}
end
