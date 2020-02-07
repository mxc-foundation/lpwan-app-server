import os
import binascii

from RESTful_api.send_request import get_request, post_request
from script.global_admin_user import all_users, user_login

from RESTful_api.network_server import getNetworkServerID
from RESTful_api.service_profile import getServiceProfileID
from RESTful_api.device_profile import getDeviceProfileID
from RESTful_api.application import getApplicationID


def getDevice(user, applicationID):
    response = get_request(api_url='api/devices?limit=999&offset=0&applicationID={}'.format(applicationID), jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/devices?limit=999&offset=0&applicationID={} response code {}".format(applicationID, response.status_code))
        exit()

    tmpRes = response.text.replace('null', '123')
    result = eval(tmpRes)
    if eval(result['totalCount']) == 0:
        return "null"
    else:
        deviceList = result['result']

    return deviceList


def createDevice(applicationID, deviceProfileID):
    user = all_users.get_user_list()[0]
    devEUI = binascii.hexlify(os.urandom(8)).decode()
    device = '''
    {{
      "device": {{
        "applicationID": "{}",
        "description": "auto generated",
        "devEUI": "{}",
        "deviceProfileID": "{}",
        "name": "device_{}",
        "referenceAltitude": 0,
        "skipFCntCheck": true,
        "tags": {{}},
        "variables": {{}}
      }}
    }}
    '''.format(applicationID, devEUI, deviceProfileID, devEUI)

    response = post_request(api_url='api/devices', jwt=user.get_user_jwt(), data=device)
    if response.status_code != 200:
        print("api/devices response code {}".format(response.status_code))
        exit()

    print("device with mac {} created..".format(devEUI))


if __name__ == "__main__":
    # user login
    user_info = '{"password": "admin", "username": "admin"}'

    if not user_login(user_info):
        print("User {} login failed!".format(user_info))
        exit(1)

    user = all_users.get_user_list()[0]
    networkServerID = getNetworkServerID()
    serviceProfileID = getServiceProfileID(networkServerID)
    applicationID = getApplicationID(serviceProfileID)
    deviceProfileID = getDeviceProfileID(applicationID, networkServerID)

    createDevice(applicationID, deviceProfileID)