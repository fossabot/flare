defmodule Flare.Domain.Resource.HTTP do
  alias Flare.Domain.Resource.HTTP.Service

  def init do
    case Service.repository() do
      nil -> {:error, "missing repository"}
      _ -> :ok
    end
  end
end
