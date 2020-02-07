from RESTful_api.send_request import get_request, post_request
from script.global_admin_user import all_users, user_login
from RESTful_api.service_profile import getServiceProfileID
from RESTful_api.network_server import getNetworkServerID


def getApplications(user):
    response = get_request(api_url='api/applications?limit=999&offset=0&organizationID=1', jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/applications?limit=999&offset=0&organizationID=1 response code {}".format(response.status_code))
        exit()

    tmpRes = response.text.replace('null', '123')
    result = eval(tmpRes)
    if eval(result['totalCount']) == 0:
        return "null"
    else:
        applicationList = result['result']

    return applicationList


def createApplication(user, serviceProfileID):
    application = '''
    {{
      "application": {{
        "description": "auto generated",
        "id": "0",
        "name": "application",
        "organizationID": "1",
        "payloadCodec": "",
        "payloadDecoderScript": "",
        "payloadEncoderScript": "",
        "serviceProfileID": "{}"
      }}
    }}
    '''.format(serviceProfileID)

    response = post_request(api_url='api/applications', jwt=user.get_user_jwt(), data=application)
    if response.status_code != 200:
        print("api/applications response code {}".format(response.status_code))
        exit()

    result = eval(response.text)
    return eval(result['id'])


def getApplicationID(serviceProfileID):
    user = all_users.get_user_list()[0]
    applicationList = getApplications(user)
    if "null" == applicationList:
        applicationID = createApplication(user, serviceProfileID)
        return applicationID
    else:
        application = applicationList[0]
        return eval(application['id'])


if __name__ == "__main__":
    # user login
    user_info = '{"password": "admin", "username": "admin"}'

    if not user_login(user_info):
        print("User {} login failed!".format(user_info))
        exit(1)

    user = all_users.get_user_list()[0]
    networkServerID = getNetworkServerID()
    serviceProfileID = getServiceProfileID(networkServerID)

    applicationList = getApplications(user)
    if "null" == applicationList:
        applicationID = createApplication(user, serviceProfileID)
        print(applicationID)
    else:
        application = applicationList[0]
        print(eval(application['id']))