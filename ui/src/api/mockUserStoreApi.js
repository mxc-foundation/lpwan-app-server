import delay from './delay';
import profilePic from './data/mockProfilePic';

const get = {
  "user": {
    "id": "100",
    "profilePic": profilePic,
    "username": "My Username",
    "sessionTTL": "1000",
    "isAdmin": true,
    "isActive": true,
    "email": "admin@mock.com",
    "note": "My Note"
  },
  "createdAt": "2019-12-06 15:46:12.624982",
  "updatedAt": "2019-12-07 15:46:12.624982"
};

class MockUserStoreApi {
  static get() {
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        resolve(Object.assign({}, get));
      }, delay);
    });
  }
}

export default MockUserStoreApi;
