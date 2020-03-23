# This model calls all internal APIs
# Accept input:
#   1. Grpc-Metadata-Authorization
#   2. jwt
# Return http response status code and text
from RESTful_api.send_request import get_request, post_request
from model.users import User
false = False
true = True


# login with username and password
def internal_login(user_info):
    '''

    :param user_info: '{"password": "", "username": ""}'
    :return: user object
    '''
    # check input
    try:
        user = eval(user_info)
        user["password"]
        user["username"]
    except KeyError:
        return "null"
    except TypeError:
        return "null"

    response = post_request(api_url='api/internal/login', jwt='', data=user_info)
    if 200 != response.status_code:
        print("Internal login request error: {} \n {}".format(response.status_code, response.text))
        return "null"

    try:
        result = eval(response.text)
        jwt = result["jwt"]
    except KeyError:
        return "null"
    except TypeError:
        return "null"

    user = User()
    user.set_user_jwt(jwt)
    return user


# get data for a particular user
def internal_profile(user):
    '''

    :param user:
    :return: user object
    '''
    jwt = user.get_user_jwt()
    response = get_request(api_url='api/internal/profile', jwt=jwt, data='')
    if 200 != response.status_code:
        print("Internal profile request error: {} \n {}".format(response.status_code, response.text))
        return "null"

    try:
        profile = eval(response.text)
    except TypeError:
        return "null"

    user.init_user(profile)
    print("User {}({}) successfully "
          "created! \n".format(user.get_username(), user.get_user_id()))
    return user

