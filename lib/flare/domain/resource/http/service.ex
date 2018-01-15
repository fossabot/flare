defmodule Flare.Domain.Resource.HTTP.Service do
  use Flare.Infra.HTTP, :controller
  alias Flare.Infra.HTTP.ErrorView

  def index(conn, _params) do
    repo = repository()
    case repo.all() do
      {:ok, result} -> render(conn, "index.json", resources: Map.get(result, :resources))
    end
  end

  def show(conn, %{"id" => id}) do
    repo = repository()
    case repo.one(id) do
      {:ok, doc} -> render(conn, "show.json", resource: doc)
      {:error, :not_found} -> render(conn, ErrorView, "404.json")
    end
  end

  def create(conn, params) do
    repo = repository()
    case repo.create(params) do
      {:ok, id} ->
        case repo.one(id) do
          {:ok, doc} -> render(conn, "show.json", resource: doc)
          {:error, :not_found} -> render(conn, ErrorView, "404.json")
        end

      {:error, detail} ->
        conn
        |> put_status(:internal_server_error)
        |> render("error.json", %{detail: detail})
    end
  end

  def delete(conn, %{"id" => id}) do
    repo = repository()
    case repo.delete(id) do
      :ok -> put_status(conn, :not_found)
      {:error, :not_found} -> render(conn, ErrorView, "404.json")
    end
  end

  def repository, do: Flare.Plugin.repository(:resource)
end
