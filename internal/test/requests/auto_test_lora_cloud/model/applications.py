class Application:
    def __init__(self):
        self.id = ''
        self.name = ''
        self.description = ''
        self.organizationID = ''
        self.serviceProfileID = ''
        self.serviceProfileName = ''
        self.payloadCodec = ''
        self.payloadEncoderScript = ''
        self.payloadDecoderScript = ''
        self.node_id_list = []

    def init_application(self, data):
        if self._init_from_application_list(data):
            return
        if self._init_from_application_id(data):
            return

    def _init_from_application_list(self, data):
        try:
            data["id"]
            data["name"]
            data["description"]
            data["organizationID"]
            data["serviceProfileID"]
            data["serviceProfileName"]
        except KeyError:
            return False

        self.id = data["id"]
        self.name = data["name"]
        self.description = data["description"]
        self.organizationID = data["organizationID"]
        self.serviceProfileID = data["serviceProfileID"]
        self.serviceProfileName = data["serviceProfileName"]
        return True

    def _init_from_application_id(self, data):
        try:
            data["id"]
            data["name"]
            data["description"]
            data["organizationID"]
            data["serviceProfileID"]
            data["payloadCodec"]
            data["payloadEncoderScript"]
            data["payloadDecoderScript"]
        except KeyError:
            return False

        self.id = data["id"]
        self.name = data["name"]
        self.description = data["description"]
        self.organizationID = data["organizationID"]
        self.serviceProfileID = data["serviceProfileID"]
        self.payloadCodec = data["payloadCodec"]
        self.payloadEncoderScript = data["payloadEncoderScript"]
        self.payloadDecoderScript = data["payloadDecoderScript"]
        return True

    def get_application_id(self):
        return self.id

    def get_application_name(self):
        return self.name

    def add_node_to_application(self, node_id):
        self.node_id_list.append(node_id)

    def add_node_list(self, node_list):
        self.node_id_list = node_list

    def get_node_list(self):
        return self.node_id_list

    def is_id(self, app_id):
        return app_id == self.id


class ApplicationList:
    def __init__(self):
        self.application_list = []

    def insert_application(self, application):
        self.application_list.append(application)

    def append_application_list(self, applist):
        self.application_list += applist

    def get_application_list(self):
        return self.application_list

    def get_application_by_id(self, app_id):
        for item in self.application_list:
            if item.is_id(app_id):
                return item
        return "null"
