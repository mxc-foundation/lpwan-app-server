/// <reference types="jest" />
import topupStore from './TopupStore';

it('getTopUpDestination', () => {
    const moneyAbbr = 0;
    const orgId = 1;
    topupStore.getTopUpDestination(moneyAbbr, orgId, 
        result => {
            // finished
            console.log(result.activeAccount);
            expect(result.activeAccount).toBeDefined();
        }, 
        err => {
            // something went wrong
            fail('should not fail');
    });
});