defmodule Flare.Domain.Resource.Repository do
  alias Flare.Domain.Resource
  alias Flare.Domain.Resource.RepositoryError
  alias Flare.Infra.HTTP.Pagination

  @callback all() :: {:ok, %{resources: [Resource], pagination: Pagination}} | {:error, RepositoryError}
  @callback one(String.t()) :: {:ok, Resource} | {:error, RepositoryError}
  @callback create(Resource) :: :ok | {:error, RepositoryError}
  @callback destroy(String.t()) :: :ok | {:error, RepositoryError}
end

