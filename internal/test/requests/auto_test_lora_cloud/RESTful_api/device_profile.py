from RESTful_api.send_request import get_request, post_request
from script.global_admin_user import all_users, user_login
from RESTful_api.network_server import getNetworkServerID
from RESTful_api.application import getApplicationID
from RESTful_api.service_profile import getServiceProfileID


def getDeviceProfile(user, applicationID):
    response = get_request(api_url='api/device-profiles?limit=999&offset=0&organizationID=1&applicationID={}'.format(applicationID), jwt=user.get_user_jwt(), data='')
    if response.status_code != 200:
        print("api/device-profiles?limit=999&offset=0&organizationID=1&applicationID={} response code {}".format(applicationID, response.status_code))
        exit()

    tmpRes = response.text.replace('null', '123')
    result = eval(tmpRes)
    if eval(result['totalCount']) == 0:
        return "null"
    else:
        deviceProfileList = result['result']

    return deviceProfileList


def createDeviceProfile(user, networkServerID):
    deviceProfile = '''
    {{
      "deviceProfile": {{
        "classBTimeout": 0,
        "classCTimeout": 0,
        "factoryPresetFreqs": [
          0
        ],
        "geolocBufferTTL": 0,
        "geolocMinBufferSize": 0,
        "id": "0",
        "macVersion": "1.0.0",
        "maxDutyCycle": 0,
        "maxEIRP": 0,
        "name": "deviceProfile",
        "networkServerID": "{}",
        "organizationID": "1",
        "payloadCodec": "None",
        "payloadDecoderScript": "",
        "payloadEncoderScript": "",
        "pingSlotDR": 0,
        "pingSlotFreq": 0,
        "pingSlotPeriod": 0,
        "regParamsRevision": "A",
        "rfRegion": "",
        "rxDROffset1": 0,
        "rxDataRate2": 0,
        "rxDelay1": 0,
        "rxFreq2": 0,
        "supports32BitFCnt": true,
        "supportsClassB": true,
        "supportsClassC": true,
        "supportsJoin": true
      }}
    }}
    '''.format(networkServerID)

    response = post_request(api_url='api/device-profiles', jwt=user.get_user_jwt(), data=deviceProfile)
    if response.status_code != 200:
        print("api/device-profiles response code {}".format(response.status_code))
        exit()

    result = eval(response.text)
    return result['id']


def getDeviceProfileID(applicationID, networkServerID):
    user = all_users.get_user_list()[0]
    deviceProfileList = getDeviceProfile(user, applicationID)
    if "null" == deviceProfileList:
        deviceProfileID = createDeviceProfile(user, networkServerID)
        return deviceProfileID
    else:
        deviceProfile = deviceProfileList[0]
        return deviceProfile['id']


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

    deviceProfileList = getDeviceProfile(user, applicationID)
    if "null" == deviceProfileList:
        deviceProfileID = createDeviceProfile(user, networkServerID)
        print(deviceProfileID)
    else:
        deviceProfile = deviceProfileList[0]
        print(deviceProfile['id'])