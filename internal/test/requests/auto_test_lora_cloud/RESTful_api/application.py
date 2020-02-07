from RESTful_api.send_request import get_request, post_request
from model.applications import Application
false = False
true = True


# get the available application list from a particular organization
def application_list(jwt, oid, limit, offset):
    '''
    updated
    :param jwt:
    :param oid:
    :param limit:
    :param offset:
    :return: list of application object
    '''
    response = get_request('api/applications?limit={}&offset={}&organizationID={}'.format(limit, offset, oid),
                           jwt=jwt, data='')
    if 200 != response.status_code:
        print("Application list request error: {} \n {}".format(response.status_code, response.text))
        return "null"

    app_obj_list = []
    try:
        result = eval(response.text)
        applist = result["result"]
    except TypeError:
        return "null"
    except KeyError:
        return "null"

    for item in applist:
        app = Application()
        app.init_application(item)
        app_obj_list.append(app)

    return app_obj_list


# get a requested application
def application_data(jwt, app_id):
    '''

    :param jwt: jwt of a particular user
    :param app_id: application id
    :return: application object
    '''
    response = get_request('api/applications/{}'.format(app_id), jwt=jwt, data='')
    if 200 != response.status_code:
        print("Application data request error: {} \n {}".format(response.status_code, response.text))
        return "null"

    try:
        profile = eval(response.text)
    except TypeError:
        return "null"

    app = Application()
    app.init_application(profile)
    print("Application {}({}) created successfully! \n".format(app.get_application_name(),
                                                               app.get_application_id()))
    return app

