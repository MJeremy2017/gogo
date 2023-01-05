from gogokit import ViagogoClient

# TODO get client ID and secret
# waiting for the affiliate team to give client ID and secret (pending)
# All methods require OAuth2 authentication. To get OAuth2 credentials for your
# application, see http://developer.viagogo.net/#authentication.
client = ViagogoClient("6779ef20e75817b79602", "")


# Get an access token. See http://developer.viagogo.net/#getting-access-tokens
token = client.oauth.get_client_access_token()
client.set_token(token)


# Get a list of events, categories, venues and metro areas that match the given
# search query
search = client.search.get_search_results({ "query": "Real Madrid" })