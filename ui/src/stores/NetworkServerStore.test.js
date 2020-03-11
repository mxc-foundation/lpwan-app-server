/// <reference types="jest" />
import networkServerStore from './NetworkServerStore';


describe('NetworkServerStore', () => {
    
it('get', (done) => {
    const id = 1;
    networkServerStore.get(id, 
        result => {
            // finished
            console.log(result.obj);
            expect(result.obj).toBeDefined();
            done();
        },
        err => {
            fail('could not get id');
        });
});

/* it('list', () => {
    const id = 1;
    const organizationID = 1;
    const limit = 10;
    const offset = 0;
    networkServerStore.list(organizationID, limit, offset, callbackFunc, errorCallbackFunc
        result => {
            console.log(result.activeAccount);
            expect(result.activeAccount).toBeDefined();
        });
}); */

});