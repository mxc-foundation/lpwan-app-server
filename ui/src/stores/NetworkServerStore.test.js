/// <reference types="jest" />
import networkServerStore from './NetworkServerStore';
import SessionStore from './SessionStore';

jest.mock('history',  () =>  ({
  createHashHistory: jest.fn() 
}))

beforeAll(async (done) => {

    await SessionStore.login({
        username: 'admin',
        password: 'admin',
        isVerified: true
    });
    done();
});


describe('NetworkServerStore', () => {
    
    it('get existing', async (done) => {
        const id = 1;
        const result = await networkServerStore.get(id);
    
        // finished
        expect(result.networkServer).toBeDefined();
        done();
    });
    
    it('get not existing', async (done) => {
        const id = 123456789;
        const result = await networkServerStore.get(id);
        
        // finished
        expect(result).not.toBeDefined();
        done();
    });

    it('list', async (done) => {
        const id = 1;
        const organizationID = 1;
        const limit = 10;
        const offset = 0;
        
        const result = await networkServerStore.list(organizationID, limit, offset);
        
        expect(result).toBeDefined();
        done();
    }); 

    /* it('delete', async (done) => {
        const id = '1';
        
        const result = await networkServerStore.delete(id);
        
        expect(result).toBeDefined();
        done();
    }); */

    /* it('update', async (done) => {
        const networkServer = {};
        
        const result = await networkServerStore.update(networkServer);
        
        expect(result).toBeDefined();
        done();
    });

    it('create', async (done) => {
        const networkServer = {};
        
        const result = await networkServerStore.create(networkServer) ;
        
        expect(result).toBeDefined();
        done();
    }); */

});