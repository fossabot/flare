defmodule Flare.Domain.Resource.HTTP.ServiceTest do
	use FlareWeb.ConnCase

	# aqui vamos carregar o reposiutory em memoria do

	setup do
		# fazer o insert aqui no mem e mongodb
		# rodar um de cada vez e ir trocando no repositorio o cara a ser testado
	end

	test "index/2 responds with all Users", %{conn: conn} do
		response = conn
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