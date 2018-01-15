defmodule Flare.Domain.Subscription.HTTP.Service do
  use Flare.Infra.HTTP, :controller
  @repository Flare.Plugin.repository(:subscription)

  def index(conn, _params) do
    render(conn, "index.json", subscriptions: @repository.all())
  end

  def create(conn, params) do
    case @repository.create(params) do
      {:ok, _} ->
        render(conn, "index.json", subscriptions: @repository.all())

      {:error, err} ->
        render(conn, "error.json", err)
    end
  end
end
