# This model calls all internal APIs
# Accept input:
#   1. Grpc-Metadata-Authorization
#   2. jwt
# Return http response status code and text
from RESTful_api.send_request import get_request, post_request
from model.users import User
false = False
true = True

