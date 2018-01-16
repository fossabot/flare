defmodule Flare.Domain.Resource.HTTP.ServiceTest do
  use FlareWeb.ConnCase

  setup context do
    case context[:provider] do
      :memory ->
        Application.put_env(
          :flare,
          :repository,
          provider: Memory
        )

      :mongodb ->
        Application.put_env(
          :flare,
          :repository,
          provider: MongoDB,
          options: [pool: DBConnection.Poolboy]
        )
    end
  end

  Enum.each([:memory, :mongodb], fn(provider) ->
    @tag provider: provider
    describe "given a #{provider} repository" do
      test("index/2 responds with all Users", %{conn: conn}, do: index(conn))
    end
  end)

  def index(conn) do
    conn
    |> get(resource_path(conn, :index))
    |> json_response(200)
  end
end

# describe "create/2" do
# 	test "Creates, and responds with a newly created user if attributes are valid"
# 	test "Returns an error and does not create a user if attributes are invalid"
# end

# describe "show/2" do
# 	test "Responds with user info if the user is found"
# 	test "Responds with a message indicating user not found"
# end

# describe "update/2" do
# 	test "Edits, and responds with the user if attributes are valid"
# 	test "Returns an error and does not edit the user if attributes are invalid"
# end

# test "delete/2 and responds with :ok if the user was deleted"
# end
