defmodule Flare.Domain.Subscription.HTTP.ServiceView do
  use Flare.Infra.HTTP, :view

  def render("index.json", %{subscriptions: subscriptions}) do
    %{
      subscriptions:
        render_many(subscriptions, __MODULE__, "subscription.json", as: :subscription)
    }
  end

  def render("subscription.json", %{subscription: subscription}) do
    %{
      id: Map.get(subscription, "id")
    }
  end

  def render("error.json", err) do
    %{
      detail: err.message
    }
  end
end
