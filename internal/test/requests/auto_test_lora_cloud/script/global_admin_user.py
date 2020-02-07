# data structures
from model.users import User, UserList
# APIs
from RESTful_api import internal

# variables for the whole cloud
all_users = UserList()


def user_login(user_info):
    '''
    updated
    :param user_info:  '{"password": "", "username": ""}'
    :return: True or False
    '''
    user = internal.internal_login(user_info)
    if "null" == user:
        return False
    all_users.insert_user(user)
    return True


def user_profile(user):
    '''
    updated
    :param user: user object
    :return: True or False
    '''
    res = internal.internal_profile(user)
    if "null" == res:
        return False
    return True


def get_user_list_of_organizaion(jwt, organization, limit, offset):
    '''
    Global admin does not belong to any organization,
    therefore here we can only get normal organization users.
    Assign values of organization to user object
    Assign user id list to organization object
    :param jwt:
    :param organization: organization object
    :param limit:
    :param offset:
    :return: True or False
    '''
    org_id = organization.get_organization_id()
    user_obj_list = organization.organization_user_list(jwt, org_id, limit, offset)
    if "null" == user_obj_list:
        return False

    user_id_list = []
    # assign values of organization to user object
    for item in user_obj_list:
        item.mark_user_in_organization(org_id, item.is_user_oadmin())
        user_id_list.append(item.get_user_id())
    all_users.append_user_list(user_obj_list)
    # assign user id list to organization object
    organization.add_user_id_list(user_id_list)
    return True
