from RESTful_api.send_request import get_request, post_request
from script.global_admin_user import all_users, user_login
from RESTful_api.network_server import getNetworkServerID


def getServiceProfile(user):
    response = get_request(api_url='api/service-profiles?limit=999&offset=0&organizationID=1', jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/service-profiles?limit=999&offset=0&organizationID=1 response code {}".format(response.status_code))
        exit()

    tmpRes = response.text.replace('null', '123')
    result = eval(tmpRes)
    if eval(result['totalCount']) == 0:
        return "null"
    else:
        serviceProfileList = result['result']

    return serviceProfileList


def createServiceProfile(user, networkServerID):
    serviceProfile = '''
    {{
      "serviceProfile": {{
        "addGWMetaData": true,
        "channelMask": "",
        "devStatusReqFreq": 0,
        "dlBucketSize": 0,
        "dlRate": 0,
        "dlRatePolicy": "DROP",
        "drMax": 5,
        "drMin": 1,
        "hrAllowed": true,
        "id": "0",
        "minGWDiversity": 0,
        "name": "serviceProfile",
        "networkServerID": "{}",
        "nwkGeoLoc": true,
        "organizationID": "1",
        "prAllowed": true,
        "raAllowed": true,
        "reportDevStatusBattery": true,
        "reportDevStatusMargin": true,
        "targetPER": 0,
        "ulBucketSize": 0,
        "ulRate": 0,
        "ulRatePolicy": "DROP"
      }}
    }}
    '''.format(networkServerID)

    response = post_request(api_url='api/service-profiles', jwt=user.get_user_jwt(), data=serviceProfile)
    if response.status_code != 200:
        print("api/service-profiles response code {}".format(response.status_code))
        exit()

    result = eval(response.text)
    return result['id']


def getServiceProfileID(networkServerID):
    user = all_users.get_user_list()[0]
    serviceProfileList = getServiceProfile(user)
    if "null" == serviceProfileList:
        serviceProfileID = createServiceProfile(user, networkServerID)
        return serviceProfileID
    else:
        serviceProfile = serviceProfileList[0]
        return serviceProfile['id']


if __name__ == "__main__":
    # user login
    user_info = '{"password": "admin", "username": "admin"}'

    if not user_login(user_info):
        print("User {} login failed!".format(user_info))
        exit(1)

    user = all_users.get_user_list()[0]
    networkServerID = getNetworkServerID()

    serviceProfileList = getServiceProfile(user)
    if "null" == serviceProfileList:
        serviceProfileID = createServiceProfile(user, networkServerID)
        print(serviceProfileID)
    else:
        serviceProfile = serviceProfileList[0]
        print(serviceProfile['id'])