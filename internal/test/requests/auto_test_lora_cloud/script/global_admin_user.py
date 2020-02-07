# data structures
from model.users import User, UserList
from model.organizations import Organization, OrganizationList
from model.applications import Application, ApplicationList
from model.gateways import Gateway, GatewayList
# APIs
from RESTful_api import api_internal, api_organization, api_application, api_gateway

# variables for the whole cloud
all_users = UserList()
all_organizations = OrganizationList()
all_applications = ApplicationList()
all_gateways = GatewayList()


def script_get_all_users():
    return all_users.get_user_list()


def script_get_all_organizations():
    return all_organizations.get_organization_list()


def script_get_all_gateways():
    return all_gateways.get_gateway_list()


def script_get_all_applications():
    return all_applications.get_application_list()


def user_login(user_info):
    '''
    updated
    :param user_info:  '{"password": "", "username": ""}'
    :return: True or False
    '''
    user = api_internal.internal_login(user_info)
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
    res = api_internal.internal_profile(user)
    if "null" == res:
        return False
    return True


def get_organization_list(jwt, limit, offset, search=''):
    '''

    :param jwt:
    :param limit:
    :param offset:
    :param search: key word
    :return: True or False
    '''
    org_obj_list = api_organization.organization_list(jwt, limit, offset)
    if "null" == org_obj_list:
        return False
    all_organizations.append_organization_list(org_obj_list)
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
    user_obj_list = api_organization.organization_user_list(jwt, org_id, limit, offset)
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


def get_application_list_of_organization(jwt, organization, limit, offset):
    '''
    :param jwt:
    :param organization: organization object
    :param limit:
    :param offset:
    :return: True or False
    '''
    org_id = organization.get_organization_id()
    app_obj_list = api_application.application_list(jwt, org_id, limit, offset)
    if "null" == app_obj_list:
        return False

    all_applications.append_application_list(app_obj_list)
    # assign app id list to organization object
    app_id_list = []
    for item in app_obj_list:
        app_id_list.append(item.get_application_id())
    organization.add_app_id_list(app_id_list)
    return True


def get_gateway_list_of_organization(jwt, organization, limit, offset):
    '''

    :param jwt:
    :param organization:
    :param limit:
    :param offset:
    :return: True or False
    '''
    org_id = organization.get_organization_id()
    gateway_obj_list = api_gateway.gateway_list(jwt, org_id, limit, offset)
    if "null" == gateway_obj_list:
        return False

    all_gateways.append_gateway_list(gateway_obj_list)
    # assign gateway id list to organization object
    gateway_id_list = []
    for item in gateway_obj_list:
        gateway_id_list.append(item.get_gateway_id())
    organization.add_gateway_id_list(gateway_id_list)
    return True
