defmodule Flare.Domain.Resource.HTTP.Service do
  use Flare.Infra.HTTP, :controller
  alias Flare.Infra.HTTP.ErrorView
  @repository Flare.Plugin.repository(:resource)

  # a ideia eh buscar a merda toda do application.get_env em todo request mesmo.

  def index(conn, _params) do
    IO.inspect(conn)

    case @repository.all() do
      {:ok, result} -> render(conn, "index.json", resources: Map.get(result, :resources))
    end
  end

  def handle_info({:EXIT, from, reason}, state) do
    IO.puts("Exit pid: #{inspect(from)} reason: #{inspect(reason)}")
    child = state[from]
    {:ok, pid} = child.start_link()
    {:noreply, Map.put(state, pid, child)}
  end

  def show(conn, %{"id" => id}) do
    case @repository.one(id) do
      {:ok, doc} -> render(conn, "show.json", resource: doc)
      {:error, :not_found} -> render(conn, ErrorView, "404.json")
    end
  end

  def create(conn, params) do
    case @repository.create(params) do
      {:ok, id} ->
        case @repository.one(id) do
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
    case @repository.delete(id) do
      :ok -> put_status(conn, :not_found)
      {:error, :not_found} -> render(conn, ErrorView, "404.json")
    end
  end
end
