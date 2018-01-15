defmodule Flare.Domain.Resource.HTTP.ServiceView do
  use Flare.Infra.HTTP, :view

  def render("index.json", %{resources: resources}) do
    %{resources: render_many(resources, __MODULE__, "show.json", as: :resource)}
  end

  def render("show.json", %{resource: resource}) do
    %{
      "id" => resource.id,
      "addresses" => resource.addresses,
      "path" => resource.path,
      "change" => resource.change,
      "createdAt" => resource.created_at
    }
  end

  def render("error.json", %{detail: detail}) do
    %{detail: detail}
  end
end
