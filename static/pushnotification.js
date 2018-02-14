/* eslint-env browser */
(function () {
	'use strict';

	if (!('serviceWorker' in navigator)) {
		return;
	}

	function getReqBody(subscription) {
		const  key = subscription.getKey ? subscription.getKey('p256dh') : '';

		return {
			endpoint: subscription.endpoint,
			key: key ? btoa(String.fromCharCode.apply(null, new Uint8Array(key))) : '',
		};
	}

	navigator.serviceWorker.register('./service-worker.js')
		.then((registration) => {
			return registration.pushManager.getSubscription()
				.then((subscription) => {
					if (subscription) {
						return subscription;
					}

					return registration.pushManager.subscribe({
						userVisibleOnly: true
					});
				})
		})
		.then((subscription) => {
			console.log(JSON.stringify(getReqBody(subscription)));
		})
		.catch((err) => {
			console.error(err);
		});
}());
