addEventListener('fetch', event => {
  event.respondWith(new Response('This service has been rewritten in Go. Please deploy the Go binary instead.', {status: 200}));
});
