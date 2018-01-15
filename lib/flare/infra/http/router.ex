defmodule Flare.Infra.HTTP.Router do
  use Flare.Infra.HTTP, :router

  pipeline :api do
    plug(:accepts, ["json"])
  end

  scope "/", Flare.Domain do
    pipe_through(:api)

    resources "/resources", Resource.HTTP.Service,
      only: [:index, :show, :create, :delete],
      name: "resource" do
      resources(
        "/subscriptions",
        Subscription.HTTP.Service,
        only: [:index, :show, :create, :delete],
        name: "subscription"
      )
    end

    resources(
      "/documents",
      Document.HTTP.Service,
      only: [:show, :update, :delete],
      name: "document"
    )
  end
end
