defmodule Flare.Domain.Resource.HTTP do
  alias Flare.Domain.Resource.HTTP.Service

  # precisa achar um nome melhor pra esse cara pq vamos verificar do worker tb que nao eh http
  # acredito que o melhor seja "Flare.Domain.Resource.Init" e a funcao check pra verificar tudo

  def init do
    case Service.repository() do
      nil -> {:error, "missing repository"}
      _ -> :ok
    end
  end
end
