from RESTful_api.send_request import get_request, post_request
from model.organizations import Organization
from model.users import User
false = False
true = True


# get data for a particular organization
def organization_data(jwt, oid):
    '''

    :param jwt: jwt for a specific user
    :param id: organization id
    :return: organization object
    '''
    response = get_request('api/organizations/{}'.format(oid), jwt=jwt, data='')
    if 200 != response.status_code:
        print("Organization data request error: {} \n {}".format(response["status_code"], response["text"]))
        return "null"

    try:
        profile = eval(response.text)
    except TypeError:
        return "null"

    organization = Organization()
    organization.init_organization(profile)
    print("Organization {}({}) successfully "
          "created! \n".format(organization.get_organization_dname(), organization.get_organization_id()))
    return organization


# get organization's user list
def organization_user_list(jwt, oid, limit, offset):
    '''

    :param jwt:
    :param oid:
    :param limit:
    :param offset:
    :return: list of user objects
    '''
    response = get_request('api/organizations/{}/users?limit={}&offset={}'.format(oid, limit, offset), jwt=jwt)
    if 200 != response.status_code:
        print("Organization user list request error: {} \n {}".format(response.status_code, response.text))
        return "null"

    try:
        result = eval(response.text)
        ulist = result["result"]
    except TypeError:
        return "null"
    except KeyError:
        return "null"

    user_obj_list = []
    for item in ulist:
        user = User()
        user.init_user(item)
        user_obj_list.append(user)
    return user_obj_list


# get data for a particular organization user
def organization_user_data(jwt, oid, uid):
    '''

    :param jwt: jwt of the loged in user
    :param oid: organization id
    :return: user object
    '''
    response = get_request(api_url='api/organizations/{}/users/{}'.format(oid, uid), jwt=jwt, data='')
    if 200 != response.status_code:
        print("Organization user data request error: {} \n {}".format(response.status_code, response.text))
        return "null"

    try:
        profile = eval(response.text)
    except TypeError:
        return "null"

    user = User()
    user.init_user(profile)
    print("User {}({}) successfully "
          "created! \n".format(user.get_username(), user.get_user_id()))
    return user


# get organization list
def organization_list(jwt, limit, offset, search=''):
    '''
    updated
    :param jwt: jwt of the global admin user
    :param limit: max length of response list
    :param offset: start with
    :param search: key words
    :return: return a list of organization objects
    '''
    response = get_request('api/organizations?limit={}&offset={}'.format(limit, offset),
                           jwt=jwt, data='')
    if 200 != response.status_code:
        print("Organization list request error: {} \n {}".format(
            response.status_code, response.text))
        return "null"

    try:
        result = eval(response.text)
        rlist = result["result"]
    except TypeError:
        return "null"
    except KeyError:
        return "null"

    org_list = []
    for item in rlist:
        organization = Organization()
        organization.init_organization(item)
        org_list.append(organization)
    return org_list
