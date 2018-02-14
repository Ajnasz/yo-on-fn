/* eslint-env browser */

function getEndpoint() {
	return self.registration.pushManager.getSubscription()
	.then((subscription) => {
		if (subscription) {
			return subscription.endpoint;
		}

		throw new Error('User not subscribed');
	});
}

self.addEventListener('activate', (event) => {
	event.waitUntil(self.clients.claim());
});

self.addEventListener('push', (event) => {
	'use strict';

	const data = event.data ? event.data.text() : '';

	console.log('PUSH NOTIFICATION', event);

	event.waitUntil(getEndpoint()
		.then(() => self.registration.showNotification('YO!', {
			body: data,
		}))
		.catch( (err) => {
			console.error('PUSH NOTIFICATION', err);
		}));
});
