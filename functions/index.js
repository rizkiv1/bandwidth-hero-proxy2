export default {
  async fetch(request, env, ctx) {
    return new Response("This service has been rewritten in Go. Please deploy the Go binary instead.", { status: 200 });
  }
};
