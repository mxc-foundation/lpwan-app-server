import profilePic from './data/mockProfilePic';

const getUser = {
  "id": "100",
  "profilePic": profilePic,
  "username": "My Username",
  "sessionTTL": "1000",
  "isAdmin": true,
  "isActive": true,
  "email": "admin@mock.com",
  "note": "My Note"
};

class MockSessionStoreApi {
  static getUser() {
    return getUser;
  }
}

export default MockSessionStoreApi;
