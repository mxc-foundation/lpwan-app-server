let initializedUnhandledRejection = false;

// hackfix for potentially underreported promise rejects
// see: https://github.com/nodejs/node/issues/9523#issuecomment-453505625
exports.initializeUnhandledRejection = function initializeUnhandledRejection() {
    if (initializedUnhandledRejection) {
        return;
    }
    initializedUnhandledRejection = true;
    process.on('unhandledRejection', error => {
        console.error('[ERROR] Unhandled Rejection (make sure to always catch your `async` errors)\n  ', error);
        //process.exit(1);
    });
};