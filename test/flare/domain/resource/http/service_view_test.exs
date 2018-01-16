defmodule Flare.Domain.Resource.HTTP.ServiceViewTest do
  use FlareWeb.ConnCase, async: true
  import Phoenix.View

  alias Flare.Domain.Resource
  alias Flare.Domain.Resource.HTTP.ServiceView, as: View

  test "render index.json" do
    assert render(View, "index.json", resources: resources_request()) == %{resources: resources_response()}
  end

  test "render show.json" do
    assert render(
      View,
      "show.json",
      resource: resources_request() |> Enum.at(0)
    ) == resources_response() |> Enum.at(0)
  end

  test "render error.json" do
    assert render(View, "error.json", detail: "Page not found") == %{detail: "Page not found"}
  end

  def resources_request do
    [
      %Resource{
        id: "id1",
        addresses: ["address.1", "address.2"],
        path: "path",
        change: %{
          "format" => "format",
          "field" => "field"
        },
        created_at: "1970-01-01T00:00:00.000Z"
      }
    ]
  end

  def resources_response do
    [
      %{
        "id" => "id1",
        "addresses" => ["address.1", "address.2"],
        "path" => "path",
        "change" => %{
          "format" => "format",
          "field" => "field"
        },
        "createdAt" => "1970-01-01T00:00:00.000Z"
      }
    ]
  end
end
