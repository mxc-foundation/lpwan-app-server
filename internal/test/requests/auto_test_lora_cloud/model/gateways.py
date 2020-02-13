class Gateway:
    def __init__(self):
        self.id = ''
        self.name = ''
        self.description = ''

        self.latitude = ''
        self.longitude = ''
        self.altitude = ''
        self.source = ''
        self.accuracy = ''

        self.createdAt = ''
        self.updatedAt = ''
        self.firstSeenAt = ''
        self.lastSeenAt = ''

        self.organizationID = ''
        self.discoveryEnabled = ''
        self.gatewayProfileID = ''
        self.networkServerID = ''
        self.boards = []

    def init_gateway(self, data):
        if self._init_from_gateway_list(data):
            return
        if self._init_from_gateway_id(data):
            return

    def _init_from_gateway_list(self, data):
        try:
            data["id"]
            data["name"]
            data["description"]
            data["createdAt"]
            data["updatedAt"]
            data["organizationID"]
            data["networkServerID"]
        except KeyError:
            return False

        self.id = data["id"]
        self.name = data["name"]
        self.description = data["description"]
        self.createdAt = data["createdAt"]
        self.updatedAt = data["updatedAt"]
        self.organizationID = data["organizationID"]
        self.networkServerID = data["networkServerID"]
        return True

    def _init_from_gateway_id(self, data):
        try:
            data["gateway"]["id"]
            data["gateway"]["name"]
            data["gateway"]["description"]

            data["gateway"]["location"]["latitude"]
            data["gateway"]["location"]["longitude"]
            data["gateway"]["location"]["altitude"]
            data["gateway"]["location"]["source"]
            data["gateway"]["location"]["accuracy"]

            data["organizationID"]
            data["discoveryEnabled"]
            data["networkServerID"]
            data["gatewayProfileID"]

            if 0 != len(data["boards"]):
                data["boards"][0]["fpgaID"]
                data["boards"][0]["fineTimestampKey"]

            data["createdAt"]
            data["updatedAt"]
            data["firstSeenAt"]
            data["lastSeenAt"]
        except KeyError:
            return False

        self.id = data["gateway"]["id"]
        self.name = data["gateway"]["name"]
        self.description = data["gateway"]["description"]

        self.latitude = data["gateway"]["location"]["latitude"]
        self.longitude = data["gateway"]["location"]["longitude"]
        self.altitude = data["gateway"]["location"]["altitude"]
        self.source = data["gateway"]["location"]["source"]
        self.accuracy = data["gateway"]["location"]["accuracy"]

        self.createdAt = data["createdAt"]
        self.updatedAt = data["updatedAt"]
        self.firstSeenAt = data["firstSeenAt"]
        self.lastSeenAt = data["lastSeenAt"]

        self.organizationID = data["organizationID"]
        self.discoveryEnabled = data["discoveryEnabled"]
        self.gatewayProfileID = data["gatewayProfileID"]
        self.networkServerID = data["networkServerID"]

        for item in data["boards"]:
            board = dict()
            board["fpgaID"] = item["fpgaID"]
            board["fineTimestampKey"] = item["fineTimestampKey"]
            self.boards.append(board)

        return True

    def get_gateway_name(self):
        return self.name

    def get_gateway_id(self):
        return self.id

    def get_organization_id(self):
        return self.organizationID

    def is_id(self, mac):
        return mac == self.id


class GatewayList:
    def __init__(self):
        self.gateway_list = []

    def insert_gateway(self, gateway):
        self.gateway_list.append(gateway)

    def append_gateway_list(self, gateway_list):
        self.gateway_list += gateway_list

    def get_gateway_list(self):
        return self.gateway_list

    def get_gateway_by_mac(self, mac):
        for item in self.gateway_list:
            if item.is_id(mac):
                return item
        return "null"

